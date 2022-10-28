package ethereum

import (
	"context"

	"github.com/savour-labs/key-locker/blockchain"
	"github.com/savour-labs/key-locker/blockchain/fallback"
	"github.com/savour-labs/key-locker/config"
	"github.com/savour-labs/key-locker/proto/keylocker"
)

const ChainName = "Ethereum"

type KeyAdaptor struct {
	fallback.KeyAdaptor
	clients *KeyLockerClient
}

func NewChainAdaptor(conf *config.Config) (blockchain.KeyAdaptor, error) {
	client, err := NewKeyLockerClient(conf)
	if err != nil {
		return nil, err
	}
	return &KeyAdaptor{
		clients: client,
	}, nil
}

func (a *KeyAdaptor) GetSupportChain(req *keylocker.SupportChainReq) (*keylocker.SupportChainRep, error) {
	return &keylocker.SupportChainRep{
		Code: keylocker.ReturnCode_SUCCESS,
		Msg:  "get support chain success",
	}, nil
}

func (a *KeyAdaptor) GetSocialKey(ctx context.Context, req *keylocker.GetSocialKeyReq) (*keylocker.GetSocialKeyRep, error) {
	// todo: 实现通过调用合约 get key
	return &keylocker.GetSocialKeyRep{
		Code: keylocker.ReturnCode_SUCCESS,
		Msg:  "get social key success",
	}, nil
}

func (a *KeyAdaptor) SetSocialKey(ctx context.Context, req *keylocker.SetSocialKeyReq) (*keylocker.SetSocialKeyRep, error) {
	// todo: 实现通过调用合约 set key
	return &keylocker.SetSocialKeyRep{
		Code: keylocker.ReturnCode_SUCCESS,
		Msg:  "set social key success",
	}, nil
}
