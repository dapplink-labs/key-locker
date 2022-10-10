package fallback

import (
	"github.com/savour-labs/key-locker/proto/keylocker"
)

type KeyAdaptor struct{}

func (a *KeyAdaptor) GetSupportChain(req *keylocker.SupportChainReq) (*keylocker.SupportChainRep, error) {
	panic("implement me")
}

func (w *KeyAdaptor) SetSocialKey(request *keylocker.SetSocialKeyReq) (*keylocker.SetSocialKeyRep, error) {
	panic("implement me")
}

func (a *KeyAdaptor) GetGasPrice(req *keylocker.GetSocialKeyReq) (*keylocker.SetSocialKeyRep, error) {
	panic("implement me")
}
