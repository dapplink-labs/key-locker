package filecoin

import (
	"context"

	"github.com/savour-labs/key-locker/blockchain"
	"github.com/savour-labs/key-locker/blockchain/fallback"
	"github.com/savour-labs/key-locker/blockchain/multiclient"
	"github.com/savour-labs/key-locker/config"
	"github.com/savour-labs/key-locker/proto/common"
	"github.com/savour-labs/key-locker/proto/keylocker"
)

const ChainName = "Filcoin"

type KeyAdaptor struct {
	fallback.KeyAdaptor
	clients *multiclient.MultiClient
}

func NewChainAdaptor(conf *config.Config) (blockchain.KeyAdaptor, error) {
	return &KeyAdaptor{
		clients: nil,
	}, nil
}

func NewLocalKeyAdaptor(network config.NetWorkType) blockchain.KeyAdaptor {
	return newKeyAdaptor()
}

func newKeyAdaptor() blockchain.KeyAdaptor {
	return &KeyAdaptor{
		clients: multiclient.New([]multiclient.Client{}),
	}
}

func (a *KeyAdaptor) getClient() {

}

func (a *KeyAdaptor) GetSupportChain(req *keylocker.SupportChainReq) (*keylocker.SupportChainRep, error) {
	return &keylocker.SupportChainRep{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get support chain success",
	}, nil
}

func (a *KeyAdaptor) GetSocialKey(ctx context.Context, req *keylocker.GetSocialKeyReq) (*keylocker.GetSocialKeyRep, error) {
	// todo: 调用 filcoin 接口获取 Social Key
	return &keylocker.GetSocialKeyRep{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get filcoin social key success",
	}, nil
}

func (a *KeyAdaptor) SetSocialKey(ctx context.Context, req *keylocker.SetSocialKeyReq) (*keylocker.SetSocialKeyRep, error) {
	// todo: 调用 filcoin 接口上传 Social Key
	return &keylocker.SetSocialKeyRep{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "set filcoin social key success",
	}, nil
}
