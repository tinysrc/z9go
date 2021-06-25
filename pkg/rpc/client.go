package rpc

import (
	"github.com/tinysrc/z9go/pkg/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Client struct {
	*grpc.ClientConn
}

// NewClient impl
func NewClient() *Client {
	return &Client{}
}

// Dial impl
func (c *Client) Dial(target string, opts ...grpc.DialOption) error {
	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		log.Fatal("grpc dial failed", zap.String("target", target))
		return err
	}
	c.ClientConn = conn
	log.Info("grpc dial success", zap.String("target", target))
	return nil
}
