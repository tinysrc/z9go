package rpc

import (
	"context"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/tinysrc/z9go/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Server struct
type Server struct {
	Server *grpc.Server
}

func dummyAuth(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

// NewServer impl
func NewServer(creds *credentials.TransportCredentials, auth grpc_auth.AuthFunc) (s *Server) {
	if auth == nil {
		auth = dummyAuth
	}
	s = &Server{}
	if creds != nil {
		s.Server = grpc.NewServer(
			grpc.Creds(*creds),
			unaryInterceptor(auth),
			streamInterceptor(auth),
		)
	} else {
		s.Server = grpc.NewServer(
			unaryInterceptor(auth),
			streamInterceptor(auth),
		)
	}
	return
}

func streamInterceptor(auth grpc_auth.AuthFunc) grpc.ServerOption {
	return grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		grpc_ctxtags.StreamServerInterceptor(),
		grpc_opentracing.StreamServerInterceptor(),
		grpc_prometheus.StreamServerInterceptor,
		grpc_zap.StreamServerInterceptor(log.Logger),
		grpc_auth.StreamServerInterceptor(auth),
		grpc_recovery.StreamServerInterceptor(),
	))
}

func unaryInterceptor(auth grpc_auth.AuthFunc) grpc.ServerOption {
	return grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_ctxtags.UnaryServerInterceptor(),
		grpc_opentracing.UnaryServerInterceptor(),
		grpc_prometheus.UnaryServerInterceptor,
		grpc_zap.UnaryServerInterceptor(log.Logger),
		grpc_auth.UnaryServerInterceptor(auth),
		grpc_recovery.UnaryServerInterceptor(),
	))
}

func (s *Server) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.Server.Serve(lis)
}

func (s *Server) Stop() {
	s.Server.GracefulStop()
}
