package blockchain

import (
	"github.com/savour-labs/key-locker/proto/keylocker"
)

type KeyAdaptor interface {
	GetSupportChain(req *keylocker.SupportChainReq) (*keylocker.SupportChainRep, error)
	SetSocialKey(req *keylocker.SetSocialKeyReq) (*keylocker.SetSocialKeyRep, error)
	GetSocialKey(req *keylocker.GetSocialKeyReq) (*keylocker.GetSocialKeyRep, error)
}
