package rpc

import (
	"errors"
	"net"

	"google.golang.org/grpc"
)

// Server struct
type Server struct {
	*grpc.Server
}

// NewServer impl
func NewServer(opts ...grpc.ServerOption) *Server {
	s := &Server{}
	srv := grpc.NewServer(opts...)
	if srv == nil {
		panic(errors.New("new server failed"))
	}
	s.Server = srv
	return s
}

// Run impl
func (s *Server) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.Server.Serve(lis)
}

// Stop impl
func (s *Server) Stop() {
	s.Server.GracefulStop()
}
