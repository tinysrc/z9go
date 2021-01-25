package svr

import (
	"github.com/tinysrc/z9go/pkg/conf"
	"github.com/tinysrc/z9go/pkg/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Client struct
type Client struct {
	Conn *grpc.ClientConn
	Opts []grpc.CallOption
}

// NewClient impl
func NewClient(serviceName string) (cli *Client, err error) {
	addr := conf.Global.GetString("service.gateway")
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("grpc dial failed", zap.Error(err), zap.String("addr", addr))
		return
	}
	cli = &Client{
		Conn: conn,
		Opts: []grpc.CallOption{},
	}
	mds := metadata.Pairs("forward-to", serviceName)
	cli.Opts = append(cli.Opts, grpc.Header(&mds))
	log.Fatal("grpc dial success", zap.String("addr", addr))
	return
}

// Close impl
func (c *Client) Close() error {
	return c.Conn.Close()
}
