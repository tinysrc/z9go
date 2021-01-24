package svr

import (
	"net"

	"github.com/tinysrc/z9go/pkg/conf"
	"github.com/tinysrc/z9go/pkg/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Server interface
type Server interface {
	Core() *grpc.Server
	Run() error
}

type server struct {
	addr   string
	listen net.Listener
	server *grpc.Server
}

// NewServer new instance
func NewServer() Server {
	addr := conf.Global.GetString("service.addr")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("create listen failed", zap.String("addr", addr))
		return nil
	}
	svr := grpc.NewServer()
	return &server{
		addr:   addr,
		listen: lis,
		server: svr,
	}
}

func (s *server) Core() *grpc.Server {
	return s.server
}

func (s *server) Run() (err error) {
	if err := s.server.Serve(s.listen); err != nil {
		log.Error("server serve failed", zap.Error(err))
	}
	return
}
