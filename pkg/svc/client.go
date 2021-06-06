package svc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"time"

	"github.com/tinysrc/z9go/pkg/conf"
	"github.com/tinysrc/z9go/pkg/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

// Client struct
type Client struct {
	handlers []grpc.UnaryClientInterceptor
	creds    credentials.TransportCredentials
	conn     *grpc.ClientConn
	md       metadata.MD
}

// NewClient impl
func NewClient() *Client {
	// 加载客户端私钥和证书
	serverName := conf.Global.GetString("service.tls.client.serverName")
	certFile := conf.Global.GetString("service.tls.client.certFile")
	keyFile := conf.Global.GetString("service.tls.client.keyFile")
	caFile := conf.Global.GetString("service.tls.caFile")
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}
	// 根证书
	rootCAs := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caFile)
	if err != nil {
		panic(err)
	}
	if !rootCAs.AppendCertsFromPEM(ca) {
		panic("rootCAs append failed")
	}
	// 创建凭证
	creds := credentials.NewTLS(&tls.Config{
		ServerName:   serverName,
		Certificates: []tls.Certificate{cert},
		RootCAs:      rootCAs,
	})
	return &Client{
		creds: creds,
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
		// timeout
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		// invoker
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// Dial impl
func (c *Client) Dial(target string) (conn *grpc.ClientConn, err error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(c.creds)}
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
