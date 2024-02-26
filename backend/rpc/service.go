package rpc

import (
	"net"

	"github.com/ethereum/go-ethereum/log"
	"github.com/savour-labs/key-locker/config"
	"github.com/savour-labs/key-locker/proto/keylocker"
	"github.com/savour-labs/key-locker/serverimpl"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartService(conf *config.Config) {
	//dispatcher, err := keydispatcher.New(conf)
	serverImpl, err := serverimpl.NewServerImpl(conf)
	if err != nil {
		log.Error("Setup dispatcher failed", "err", err)
		panic(err)
	}
	//grpcServer := grpc.NewServer(grpc.UnaryInterceptor(dispatcher.Interceptor))
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(serverImpl.Interceptor))
	defer grpcServer.GracefulStop()
	//keylocker.RegisterLeyLockerServiceServer(grpcServer, dispatcher)
	keylocker.RegisterLeyLockerServiceServer(grpcServer, serverImpl)
	listen, err := net.Listen("tcp", ":"+conf.RpcServer.Port)
	if err != nil {
		log.Error("net listen failed", "err", err)
		panic(err)
	}
	reflection.Register(grpcServer)
	log.Info("savour dao start success", "port", conf.RpcServer.Port)
	if err := grpcServer.Serve(listen); err != nil {
		log.Error("grpc server serve failed", "err", err)
		panic(err)
	}
}
