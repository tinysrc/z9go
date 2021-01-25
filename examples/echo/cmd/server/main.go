package main

import (
	"context"

	"github.com/tinysrc/z9go/examples/echo/pb"
	"github.com/tinysrc/z9go/pkg/log"
	"github.com/tinysrc/z9go/pkg/srv"
	"go.uber.org/zap"
)

type service struct {
	pb.UnimplementedEchoServer
}

func (s *service) Echo(ctx context.Context, in *pb.StringMessage) (*pb.StringMessage, error) {
	log.Debug("Echo", zap.Any("in", in))
	return &pb.StringMessage{Value: in.Value}, nil
}

func main() {
	defer log.Close()
	s := srv.NewServer()
	pb.RegisterEchoServer(s.Server, &service{})
	s.Run()
}
