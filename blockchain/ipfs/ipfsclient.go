package ipfs

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	files "github.com/ipfs/go-ipfs-files"
	icore "github.com/ipfs/interface-go-ipfs-core"
	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/ipfs/kubo/config"
	"github.com/ipfs/kubo/core"
	"github.com/ipfs/kubo/core/coreapi"
	"github.com/ipfs/kubo/core/node/libp2p"
	"github.com/ipfs/kubo/plugin/loader"
	"github.com/ipfs/kubo/repo/fsrepo"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
)

type Client struct {
	ipfs        icore.CoreAPI
	node        *core.IpfsNode
	networkNode []string
	repoPath    string
}

func New(ctx context.Context, networkNode []string, repoPath string) (*Client, error) {
	ipfs, node, err := spawnEphemeral(ctx, repoPath)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := connectToPeers(ctx, ipfs, networkNode); err != nil {
			log.Printf("failed connect to peers: %s", err)
		}
	}()

	return &Client{
		ipfs:        ipfs,
		node:        node,
		networkNode: networkNode,
		repoPath:    repoPath,
	}, nil
}

// AddFile 添加文件，返回cid
func (c *Client) AddFile(ctx context.Context, file []byte) (string, error) {
	peerCidFile, err := c.ipfs.Unixfs().Add(ctx, files.NewBytesFile(file))

	if err != nil {
		return "", errors.WithMessage(err, "add file fail")
	}

	return peerCidFile.String(), nil
}

// GetFile get file from local and network
func (c *Client) GetFile(ctx context.Context, cidStr string) ([]byte, error) {
	cid := icorepath.New(cidStr)
	node, err := c.ipfs.Unixfs().Get(ctx, cid)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(files.ToFile(node))
}

var loadPluginsOnce sync.Once

// Spawns a node to be used just for this run (i.e. creates a tmp repo)
func spawnEphemeral(ctx context.Context, repoPath string) (icore.CoreAPI, *core.IpfsNode, error) {
	var onceErr error
	loadPluginsOnce.Do(func() {
		onceErr = setupPlugins("")
	})
	if onceErr != nil {
		return nil, nil, onceErr
	}

	// private network
	// cfg := &config.Config{
	// 	Identity: config.Identity{
	// 		PeerID:  "12D3KooWLXzGF1pXMYsNgv7yCGKdRwPqF2WAuqJwNFWvP2h3fNSp",
	// 		PrivKey: "CAESQOOJa3+r1ahGwYkEXuQkzgN8CvCQfccYkiNuKJ4XR+tRnz5+tOrTuExk12SMjbrejfrgtc6zR+HRoP624YCdKx0=",
	// 	},
	// 	Addresses: config.Addresses{
	// 		Swarm: []string{"/ip4/0.0.0.0/tcp/4001", "/ip4/0.0.0.0/udp/4001/quic"},
	// 		API:   []string{"/ip4/127.0.0.1/tcp/5001"},
	// 	},
	// }
	//
	// r := &repo.Mock{
	// 	C: *cfg,
	// 	D: syncds.MutexWrap(datastore.NewMapDatastore()),
	// }
	//
	// nodeOptions := &core.BuildCfg{
	// 	Online:  true,
	// 	Routing: libp2p.DHTOption, // This option sets the node to be a full DHT node (both fetching and storing DHT Records)
	// 	// Routing: libp2p.DHTClientOption, // This option sets the node to be a client DHT node (only fetching records)
	// 	Repo: r,
	// }
	//
	// node, err := core.NewNode(ctx, nodeOptions)
	// if err != nil {
	// 	return nil, nil, err
	// }

	node, err := createNode(ctx, repoPath)
	if err != nil {
		return nil, nil, err
	}

	api, err := coreapi.NewCoreAPI(node)

	return api, node, err
}

func setupPlugins(externalPluginsPath string) error {
	// Load any external plugins if available on externalPluginsPath
	plugins, err := loader.NewPluginLoader(filepath.Join(externalPluginsPath, "plugins"))
	if err != nil {
		return fmt.Errorf("error loading plugins: %s", err)
	}

	// Load preloaded and external plugins
	if err := plugins.Initialize(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	if err := plugins.Inject(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	return nil
}

// CreateRepo Create a Temporary Repo
func CreateRepo() (string, error) {
	var onceErr error
	loadPluginsOnce.Do(func() {
		onceErr = setupPlugins("")
	})
	if onceErr != nil {
		return "", onceErr
	}

	repoPath, err := os.MkdirTemp("", "ipfs-shell")
	if err != nil {
		return "", fmt.Errorf("failed to get temp dir: %s", err)
	}

	// Create a config with default options and a 2048 bit key
	cfg, err := config.Init(io.Discard, 2048)
	if err != nil {
		return "", err
	}

	// Create the repo with the config
	err = fsrepo.Init(repoPath, cfg)
	if err != nil {
		return "", fmt.Errorf("failed to init ephemeral node: %s", err)
	}

	return repoPath, nil
}

// createNode Creates an IPFS node and returns its coreAPI
func createNode(ctx context.Context, repoPath string) (*core.IpfsNode, error) {
	// Open the repo
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}

	// Construct the node
	nodeOptions := &core.BuildCfg{
		Online:  true,
		Routing: libp2p.DHTOption, // This option sets the node to be a full DHT node (both fetching and storing DHT Records)
		// Routing: libp2p.DHTClientOption, // This option sets the node to be a client DHT node (only fetching records)
		Repo: repo,
	}

	return core.NewNode(ctx, nodeOptions)
}

// connectToPeers connect network nodes
func connectToPeers(ctx context.Context, ipfs icore.CoreAPI, peers []string) error {
	var wg sync.WaitGroup
	peerInfos := make(map[peer.ID]*peer.AddrInfo, len(peers))
	for _, addrStr := range peers {
		addr, err := ma.NewMultiaddr(addrStr)
		if err != nil {
			return err
		}
		pii, err := peer.AddrInfoFromP2pAddr(addr)
		if err != nil {
			return err
		}
		pi, ok := peerInfos[pii.ID]
		if !ok {
			pi = &peer.AddrInfo{ID: pii.ID}
			peerInfos[pi.ID] = pi
		}
		pi.Addrs = append(pi.Addrs, pii.Addrs...)
	}

	wg.Add(len(peerInfos))
	for _, peerInfo := range peerInfos {
		go func(peerInfo *peer.AddrInfo) {
			defer wg.Done()
			err := ipfs.Swarm().Connect(ctx, *peerInfo)
			if err != nil {
				log.Printf("failed to connect to %s: %s", peerInfo.ID, err)
			}
		}(peerInfo)
	}
	wg.Wait()
	return nil
}
