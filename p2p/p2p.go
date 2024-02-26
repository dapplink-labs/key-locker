package p2p

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	"github.com/multiformats/go-multiaddr"
	"github.com/savour-labs/key-locker/config"
	"github.com/savour-labs/key-locker/keydispatcher"
)

type P2PHost struct {
	host.Host
	*keydispatcher.Dispatcher
	//XXProtocol          *XXProtocol // any protocol impl
	FaultToleranceTimes int
	GetKeyProtocol      *GetKeyProtocol
}

type discoveryNotifee struct {
	PeerChan chan peer.AddrInfo
}

// interface to be called when new  peer is found
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	n.PeerChan <- pi
}

func NewP2PHost(conf *config.Config, dispatcher *keydispatcher.Dispatcher) (*P2PHost, error) {
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", conf.P2P.ListenHost, conf.P2P.ListenPort))
	privKeyStr := os.Getenv(conf.P2P.HostPrivateKeyEnvFlag)
	privKey, err := crypto.UnmarshalPrivateKey([]byte(privKeyStr))
	if err != nil {
		return nil, err
	}

	connmgr, err := connmgr.NewConnManager( // 连的peer如果太多到达hi，断开一些到low
		100, // Lowwater
		400, // HighWater,
		connmgr.WithGracePeriod(time.Minute),
	)

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	host, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(privKey), // 为了security部分更加安全，好像不设置更好？没想好！!
		libp2p.ConnectionManager(connmgr),
		libp2p.Security(noise.ID, noise.New), // 并不依赖Identity，有Identity就用这个，没有就自己创建和协商；；
	)
	if err != nil {
		return nil, err
	}
	p2pHost := &P2PHost{Host: host, FaultToleranceTimes: 3}
	p2pHost.GetKeyProtocol = NewGetKeyProtocol(p2pHost)
	p2pHost.Dispatcher = dispatcher
	return p2pHost, nil
}

func (h *P2PHost) Recovery(rendezvous string) {
	ctx := context.Background()
	peerChan := initMDNS(h, rendezvous)
	for { // allows multiple peers to join
		peer := <-peerChan // will block until we discover a peer
		if peer.ID > h.ID() {
			// if other end peer id greater than us, don't connect to it, just wait for it to connect us
			log.Debug("Found peer:", peer, " id is greater than us, wait for it to connect to us")
			continue
		}
		log.Debug("Found peer:", peer, ", connecting")

		if err := h.Connect(ctx, peer); err != nil {
			log.Debug("Connection failed:", err)
			continue
		}
		h.Peerstore().AddAddrs(peer.ID, peer.Addrs, time.Hour*12)
	}
}

// Initialize the MDNS service
func initMDNS(peerhost host.Host, rendezvous string) chan peer.AddrInfo {
	// register with service so that we get notified about peer discovery
	n := &discoveryNotifee{}
	n.PeerChan = make(chan peer.AddrInfo)

	// An hour might be a long long period in practical applications. But this is fine for us
	ser := mdns.NewMdnsService(peerhost, rendezvous, n)
	if err := ser.Start(); err != nil {
		panic(err)
	}
	return n.PeerChan
}
