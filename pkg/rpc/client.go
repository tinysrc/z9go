package rpc

import (
	"context"
	"fmt"
	"time"

	"github.com/tinysrc/z9go/pkg/log"
	"github.com/tinysrc/z9go/pkg/mw/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	handlers []grpc.UnaryClientInterceptor
	conn     *grpc.ClientConn
	sign     string
}

// NewClient impl
func NewClient(sign string) *Client {
	return &Client{
		sign: sign,
	}
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
		// auth
		claims := auth.CustomClaims{}
		token, _ := auth.NewJWT(c.sign).MakeToken(claims)
		ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", fmt.Sprintf("Basic %s", token))
		// timeout
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		// invoker
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// Dial impl
func (c *Client) Dial(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithChainUnaryInterceptor(c.allHandlers()...))
	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		log.Fatal("grpc dial failed", zap.String("target", target))
		return nil, err
	}
	c.conn = conn
	log.Info("grpc dial success", zap.String("target", target))
	return conn, nil
}
