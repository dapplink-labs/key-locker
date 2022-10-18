package blockchain

import (
	"context"

	"github.com/savour-labs/key-locker/proto/keylocker"
)

type KeyAdaptor interface {
	GetSupportChain(req *keylocker.SupportChainReq) (*keylocker.SupportChainRep, error)
	SetSocialKey(ctx context.Context, req *keylocker.SetSocialKeyReq) (*keylocker.SetSocialKeyRep, error)
	GetSocialKey(ctx context.Context, req *keylocker.GetSocialKeyReq) (*keylocker.GetSocialKeyRep, error)
}
