package keydispatcher

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/savour-labs/key-locker/blockchain"
	"github.com/savour-labs/key-locker/blockchain/ethereum"
	"github.com/savour-labs/key-locker/blockchain/filecoin"
	"github.com/savour-labs/key-locker/blockchain/ipfs"
	"github.com/savour-labs/key-locker/blockchain/moonbeam"
	"github.com/savour-labs/key-locker/config"
)

type CommonRequest interface {
	GetChain() string
}

type ChainType = string

type Dispatcher struct {
	Registry map[ChainType]blockchain.KeyAdaptor
}

func New(conf *config.Config) (*Dispatcher, error) {
	dispatcher := Dispatcher{
		Registry: make(map[ChainType]blockchain.KeyAdaptor),
	}
	keyAdaptorFactoryMap := map[string]func(conf *config.Config) (blockchain.KeyAdaptor, error){
		ethereum.ChainName: ethereum.NewChainAdaptor,
		moonbeam.ChainName: moonbeam.NewChainAdaptor,
		ipfs.ChainName:     ipfs.NewChainAdaptor,
		filecoin.ChainName: filecoin.NewChainAdaptor,
	}
	supportedChains := []string{ethereum.ChainName, moonbeam.ChainName, ipfs.ChainName, filecoin.ChainName}
	for _, c := range conf.Chains {
		if factory, ok := keyAdaptorFactoryMap[c]; ok {
			adaptor, err := factory(conf)
			if err != nil {
				log.Crit("failed to setup chain", "chain", c, "error", err)
			}
			dispatcher.Registry[c] = adaptor
		} else {
			log.Error("unsupported chain", "chain", c, "supportedChains", supportedChains)
		}
	}
	return &dispatcher, nil
}

//
//func (d *Dispatcher) Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
//	defer func() {
//		if e := recover(); e != nil {
//			log.Error("panic error", "msg", e)
//			log.Debug(string(debug.Stack()))
//			err = status.Errorf(codes.Internal, "Panic err: %v", e)
//		}
//	}()
//	pos := strings.LastIndex(info.FullMethod, "/")
//	method := info.FullMethod[pos+1:]
//	chain := req.(CommonRequest).GetChain()
//	log.Info(method, "chain", chain, "req", req)
//	resp, err = handler(ctx, req)
//	log.Debug("Finish handling", "resp", resp, "err", err)
//	return
//}
//
//func (d *Dispatcher) preHandler(req interface{}) (resp *keylocker.SupportChainRep) {
//	chain := req.(CommonRequest).GetChain() // req调用时的如参的结构内有chainID，并且实现了getChain方法拿出相应字段
//	if _, ok := d.registry[chain]; !ok {
//		return &keylocker.SupportChainRep{
//			Code:    keylocker.ReturnCode_ERROR,
//			Msg:     config.UnsupportedOperation,
//			Support: false,
//		}
//	}
//	return nil
//}
//
//func (d *Dispatcher) GetSupportChain(ctx context.Context, req *keylocker.SupportChainReq) (*keylocker.SupportChainRep, error) {
//	resp := d.preHandler(req) // 其实就是检查d.registry[chain]这个东西有没有，如果没有就返回响应错误
//	if resp != nil {
//		return &keylocker.SupportChainRep{
//			Code: keylocker.ReturnCode_ERROR,
//			Msg:  config.UnsupportedOperation,
//		}, nil
//	}
//	return d.registry[req.Chain].GetSupportChain(req)
//}
//
//func (d *Dispatcher) SetSocialKey(ctx context.Context, req *keylocker.SetSocialKeyReq) (*keylocker.SetSocialKeyRep, error) {
//	resp := d.preHandler(req)
//	if resp != nil {
//		return &keylocker.SetSocialKeyRep{
//			Code: keylocker.ReturnCode_ERROR,
//			Msg:  config.UnsupportedOperation,
//		}, nil
//	}
//	return d.registry[req.Chain].SetSocialKey(ctx, req)
//}
//
//func (d *Dispatcher) GetSocialKey(ctx context.Context, req *keylocker.GetSocialKeyReq) (*keylocker.GetSocialKeyRep, error) {
//	resp := d.preHandler(req)
//	if resp != nil {
//		return &keylocker.GetSocialKeyRep{
//			Code: keylocker.ReturnCode_ERROR,
//			Msg:  config.UnsupportedOperation,
//		}, nil
//	}
//	//return d.registry[req.Chain].GetSocialKey(ctx, req)
//	res, err := d.registry[req.Chain].GetSocialKey(ctx, req)
//	if err != nil {
//		if errors.Is(err, errors.New("Fault Tolerance")) {
//			d.p2pHost.FaultTolerance(p2p.GetKeyRequest)
//		}
//	}
//	return res, nil
//}
//
//func (d *Dispatcher) StreamHandler(s network.Stream) {
//
//}
