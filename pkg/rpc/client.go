package rpc

import (
	"context"
	"time"

	"github.com/tinysrc/z9go/pkg/log"
	"github.com/tinysrc/z9go/pkg/mw/utils"
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

func (c *Client) lastUnaryHandler() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// timeout
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		// invoker
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (c *Client) lastStreamHandler() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		// timeout
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		// client stream
		s, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			return nil, err
		}
		w := utils.WrapClientStream(s)
		w.WrappedCtx = ctx
		return w, nil
	}
}

// Dial impl
func (c *Client) Dial(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithChainUnaryInterceptor(c.lastUnaryHandler()))
	opts = append(opts, grpc.WithChainStreamInterceptor(c.lastStreamHandler()))
	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		log.Fatal("grpc dial failed", zap.String("target", target))
		return nil, err
	}
	c.ClientConn = conn
	log.Info("grpc dial success", zap.String("target", target))
	return conn, nil
}
