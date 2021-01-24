package main

import (
	"context"
	"net"

	"github.com/tinysrc/z9go/examples/echo/pb"
	"github.com/tinysrc/z9go/pkg/conf"
	"github.com/tinysrc/z9go/pkg/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type service struct {
	pb.UnimplementedEchoServiceServer
}

func (s *service) Echo(ctx context.Context, in *pb.StringMessage) (*pb.StringMessage, error) {
	return &pb.StringMessage{Value: in.Value}, nil
}

func main() {
	defer log.Close()
	addr := conf.Global.GetString("service.addr")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("create listen failed", zap.String("addr", addr))
		return
	}
	svr := grpc.NewServer()
	pb.RegisterEchoServiceServer(svr, &service{})
	if err := svr.Serve(lis); err != nil {
		log.Fatal("grpc serve failed", zap.Error(err))
	}
}
