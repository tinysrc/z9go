package svc

import (
	"context"
	"time"

	"github.com/tinysrc/z9go/pkg/conf"
	"github.com/tinysrc/z9go/pkg/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Client struct
type Client struct {
	handlers []grpc.UnaryClientInterceptor
	conn     *grpc.ClientConn
	md       metadata.MD
}

// NewClient impl
func NewClient() *Client {
	return &Client{}
}

// Use impl
func (c *Client) Use(handlers ...grpc.UnaryClientInterceptor) *Client {
	l1 := len(c.handlers)
	l2 := len(handlers)
	hs := make([]grpc.UnaryClientInterceptor, l1+l2)
	copy(hs, c.handlers)
	copy(hs[l1:], handlers)
	c.handlers = hs
	return c
}

func (c *Client) allHandlers() (hs []grpc.UnaryClientInterceptor) {
	l1 := len(c.handlers)
	hs = make([]grpc.UnaryClientInterceptor, l1+1)
	copy(hs, c.handlers)
	hs[l1] = c.getlastHandler()
	return
}

func (c *Client) getlastHandler() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// timeout
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		// invoker
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// Dial impl
func (c *Client) Dial(target string) (conn *grpc.ClientConn, err error) {
	opts := []grpc.DialOption{grpc.WithInsecure()}
	opts = append(opts, grpc.WithChainUnaryInterceptor(c.allHandlers()...))
	addr := conf.Global.GetString("service.gateway")
	conn, err = grpc.Dial(addr, opts...)
	if err != nil {
		log.Fatal("grpc dial failed", zap.Error(err), zap.String("addr", addr))
		return
	}
	c.conn = conn
	c.md = metadata.Pairs("Z9-Svc", target)
	log.Info("grpc dial success", zap.String("addr", addr))
	return
}

// NewCallCtx impl
func (c *Client) NewCallCtx() context.Context {
	return metadata.NewOutgoingContext(context.Background(), c.md)
}
