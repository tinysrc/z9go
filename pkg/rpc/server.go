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
)

// Server struct
type Server struct {
	Server *grpc.Server
}

func dummyAuth(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

// NewServer impl
func NewServer(authFunc grpc_auth.AuthFunc) *Server {
	if authFunc == nil {
		authFunc = dummyAuth
	}
	s := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_opentracing.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
			grpc_zap.StreamServerInterceptor(log.Logger),
			grpc_auth.StreamServerInterceptor(authFunc),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_zap.UnaryServerInterceptor(log.Logger),
			grpc_auth.UnaryServerInterceptor(authFunc),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)
	return &Server{
		Server: s,
	}
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
