package svc

import (
	"context"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/tinysrc/z9go/pkg/conf"
	"github.com/tinysrc/z9go/pkg/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Server struct
type Server struct {
	addr   string
	listen net.Listener
	Server *grpc.Server
}

func auth(ctx context.Context) (context.Context, error) {
	return nil, nil
}

// NewServer impl
func NewServer() *Server {
	addr := conf.Global.GetString("service.addr")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("create listen failed", zap.String("addr", addr))
		return nil
	}
	svr := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_opentracing.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
			grpc_zap.StreamServerInterceptor(log.Logger),
			grpc_auth.StreamServerInterceptor(auth),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_zap.UnaryServerInterceptor(log.Logger),
			grpc_auth.UnaryServerInterceptor(auth),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)
	return &Server{
		addr:   addr,
		listen: lis,
		Server: svr,
	}
}

// Run impl
func (s *Server) Run() (err error) {
	name := conf.Global.GetString("service.name")
	addr := conf.Global.GetString("service.addr")
	log.Info("service start", zap.String("serviceName", name), zap.String("serviceAddr", addr))
	if err := s.Server.Serve(s.listen); err != nil {
		log.Error("server serve failed", zap.Error(err))
	}
	return
}
