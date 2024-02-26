package serverimpl

import (
	"context"
	"errors"
	"runtime/debug"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/gogo/protobuf/proto"
	"github.com/savour-labs/key-locker/config"
	"github.com/savour-labs/key-locker/keydispatcher"
	"github.com/savour-labs/key-locker/p2p"
	"github.com/savour-labs/key-locker/proto/keylocker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerImpl struct {
	*keydispatcher.Dispatcher
	p2pHost *p2p.P2PHost
}

func NewServerImpl(conf *config.Config) (*ServerImpl, error) {
	dispatcher, err := keydispatcher.New(conf)
	if err != nil {
		return nil, err
	}
	p2pHost, err := p2p.NewP2PHost(conf, dispatcher)
	if err != nil {
		log.Error("new libp2p host failed:", err)
	}
	sImpl := &ServerImpl{dispatcher, p2pHost}

	if p2pHost != nil {
		go p2pHost.Recovery(conf.P2P.RendezousString)
	}
	return sImpl, nil
}

func (d *ServerImpl) Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error("panic error", "msg", e)
			log.Debug(string(debug.Stack()))
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()
	pos := strings.LastIndex(info.FullMethod, "/")
	method := info.FullMethod[pos+1:]
	chain := req.(keydispatcher.CommonRequest).GetChain()
	log.Info(method, "chain", chain, "req", req)
	resp, err = handler(ctx, req)
	log.Debug("Finish handling", "resp", resp, "err", err)
	return
}

func (d *ServerImpl) preHandler(req interface{}) (resp *keylocker.SupportChainRep) {
	chain := req.(keydispatcher.CommonRequest).GetChain() // req调用时的如参的结构内有chainID，并且实现了getChain方法拿出相应字段
	if _, ok := d.Registry[chain]; !ok {
		return &keylocker.SupportChainRep{
			Code:    keylocker.ReturnCode_ERROR,
			Msg:     config.UnsupportedOperation,
			Support: false,
		}
	}
	return nil
}

func (d *ServerImpl) GetSupportChain(ctx context.Context, req *keylocker.SupportChainReq) (*keylocker.SupportChainRep, error) {
	resp := d.preHandler(req) // 其实就是检查d.registry[chain]这个东西有没有，如果没有就返回响应错误
	if resp != nil {
		return &keylocker.SupportChainRep{
			Code: keylocker.ReturnCode_ERROR,
			Msg:  config.UnsupportedOperation,
		}, nil
	}
	return d.Registry[req.Chain].GetSupportChain(req)
}

func (d *ServerImpl) SetSocialKey(ctx context.Context, req *keylocker.SetSocialKeyReq) (*keylocker.SetSocialKeyRep, error) {
	resp := d.preHandler(req)
	if resp != nil {
		return &keylocker.SetSocialKeyRep{
			Code: keylocker.ReturnCode_ERROR,
			Msg:  config.UnsupportedOperation,
		}, nil
	}
	return d.Registry[req.Chain].SetSocialKey(ctx, req)
}

func (d *ServerImpl) GetSocialKey(ctx context.Context, req *keylocker.GetSocialKeyReq) (*keylocker.GetSocialKeyRep, error) {
	resp := d.preHandler(req)
	if resp != nil {
		return &keylocker.GetSocialKeyRep{
			Code: keylocker.ReturnCode_ERROR,
			Msg:  config.UnsupportedOperation,
		}, nil
	}
	//return d.registry[req.Chain].GetSocialKey(ctx, req)
	res, err := d.Registry[req.Chain].GetSocialKey(ctx, req)
	if err != nil {
		if errors.Is(err, errors.New("Fault Tolerance")) {
			reqBytes, err := proto.Marshal(req)
			if err != nil {
				log.Error("proto marshal origin request failed", err)
				return nil, errors.New("proto marshal origin request failed")
			}

			ticker := time.NewTicker(3 * time.Second)

			relayMsgId, err := d.p2pHost.FaultTolerance(p2p.GetKeyRequest, reqBytes)
			for {
				select {
				case <-ticker.C:
				case e := <-d.p2pHost.GetKeyProtocol.Requests[relayMsgId]:
					if e == nil {
						return nil, err
					}
					return e, nil
				}
			}
		}
		return nil, err
	}
	return res, nil
}
