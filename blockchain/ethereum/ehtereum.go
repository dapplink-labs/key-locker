package ethereum

import (
	"github.com/savour-labs/key-locker/blockchain"
	"github.com/savour-labs/key-locker/blockchain/fallback"
	"github.com/savour-labs/key-locker/blockchain/multiclient"
	"github.com/savour-labs/key-locker/config"
	"github.com/savour-labs/key-locker/proto/common"
	"github.com/savour-labs/key-locker/proto/keylocker"
)

const ChainName = "Ethereum"

type KeyAdaptor struct {
	fallback.KeyAdaptor
	clients *multiclient.MultiClient
}

func NewChainAdaptor(conf *config.Config) (blockchain.KeyAdaptor, error) {
	clients, err := newEthClients(conf)
	if err != nil {
		return nil, err
	}
	clis := make([]multiclient.Client, len(clients))
	for i, client := range clients {
		clis[i] = client
	}
	return &KeyAdaptor{
		clients: multiclient.New(clis),
	}, nil
}

func NewLocalKeyAdaptor(network config.NetWorkType) blockchain.KeyAdaptor {
	return newKeyAdaptor(newLocalEthClient(network))
}

func newKeyAdaptor(client *ethClient) blockchain.KeyAdaptor {
	return &KeyAdaptor{
		clients: multiclient.New([]multiclient.Client{client}),
	}
}

func (a *KeyAdaptor) getClient() *ethClient {
	return a.clients.BestClient().(*ethClient)
}

func (a *KeyAdaptor) GetSocialKey(req *keylocker.GetSocialKeyReq) (*keylocker.GetSocialKeyRep, error) {
	return &keylocker.GetSocialKeyRep{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get social key success",
	}, nil
}
