package svr

import (
	"github.com/tinysrc/z9go/pkg/conf"
	"github.com/tinysrc/z9go/pkg/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Client struct
type Client struct {
	Conn *grpc.ClientConn
	Kvs  []string
}

// NewClient impl
func NewClient(serviceName string) (*Client, error) {
	addr := conf.Global.GetString("service.gateway")
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("grpc dial failed", zap.Error(err), zap.String("addr", addr))
		return nil, err
	}
	cli := &Client{
		Conn: conn,
		Kvs:  []string{"GRPC-Service", serviceName},
	}
	log.Info("grpc dial success", zap.String("addr", addr))
	return cli, nil
}

// Close impl
func (c *Client) Close() error {
	return c.Conn.Close()
}
