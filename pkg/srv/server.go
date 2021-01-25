package srv

import (
	"net"

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

// NewServer impl
func NewServer() *Server {
	addr := conf.Global.GetString("service.addr")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("create listen failed", zap.String("addr", addr))
		return nil
	}
	svr := grpc.NewServer()
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
