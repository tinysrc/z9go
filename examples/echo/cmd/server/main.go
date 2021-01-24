package main

import (
	"context"

	"github.com/tinysrc/z9go/examples/echo/pb"
	"github.com/tinysrc/z9go/pkg/log"
	"github.com/tinysrc/z9go/pkg/svr"
)

type service struct {
	pb.UnimplementedEchoServer
}

func (s *service) Echo(ctx context.Context, in *pb.StringMessage) (*pb.StringMessage, error) {
	return &pb.StringMessage{Value: in.Value}, nil
}

func main() {
	defer log.Close()
	s := svr.NewServer()
	pb.RegisterEchoServer(s.Core(), &service{})
	s.Run()
}
