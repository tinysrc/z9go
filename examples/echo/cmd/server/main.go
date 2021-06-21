package main

import (
	"context"

	"github.com/tinysrc/z9go/examples/echo/pb"
	"github.com/tinysrc/z9go/pkg/conf"
	"github.com/tinysrc/z9go/pkg/log"
	"github.com/tinysrc/z9go/pkg/rpc"
	"go.uber.org/zap"
)

type service struct {
	pb.UnimplementedEchoServiceServer
}

func (s *service) Echo(ctx context.Context, in *pb.StringMessage) (*pb.StringMessage, error) {
	log.Debug("Echo", zap.Any("in", in))
	return &pb.StringMessage{Value: in.Value}, nil
}

func main() {
	defer log.Close()
	s := rpc.NewServer(nil)
	pb.RegisterEchoServiceServer(s.Server, &service{})
	addr := conf.Global.GetString("service.addr")
	s.Run(addr)
}
