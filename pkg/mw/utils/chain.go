package utils

import (
	"context"

	"google.golang.org/grpc"
)

func ChainUnaryServer(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	sz := len(interceptors)
	if sz == 0 {
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			return handler(ctx, req)
		}
	}
	if sz == 1 {
		return interceptors[0]
	}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		curr := handler
		for i := sz - 1; i > 0; i-- {
			temp, i := curr, i
			curr = func(myctx context.Context, myreq interface{}) (interface{}, error) {
				return interceptors[i](myctx, myreq, info, temp)
			}
		}
		return interceptors[0](ctx, req, info, curr)
	}
}

func ChainStreamServer(interceptors ...grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	sz := len(interceptors)
	if sz == 0 {
		return func(srv interface{}, ss grpc.ServerStream, _info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			return handler(srv, ss)
		}
	}
	if sz == 1 {
		return interceptors[0]
	}
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		curr := handler
		for i := sz - 1; i > 0; i-- {
			temp, i := curr, i
			curr = func(mysrv interface{}, myss grpc.ServerStream) error {
				return interceptors[i](mysrv, myss, info, temp)
			}
		}
		return interceptors[0](srv, ss, info, curr)
	}
}

func ChainUnaryClient(interceptors ...grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor {
	sz := len(interceptors)
	if sz == 0 {
		return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
	}
	if sz == 1 {
		return interceptors[0]
	}
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		curr := invoker
		for i := sz - 1; i > 0; i-- {
			temp, i := curr, i
			curr = func(myctx context.Context, mymethod string, myreq, myreply interface{}, mycc *grpc.ClientConn, myopts ...grpc.CallOption) error {
				return interceptors[i](myctx, mymethod, myreq, myreply, mycc, temp, myopts...)
			}
		}
		return interceptors[0](ctx, method, req, reply, cc, invoker, opts...)
	}
}

func ChainStreamClient(interceptors ...grpc.StreamClientInterceptor) grpc.StreamClientInterceptor {
	sz := len(interceptors)
	if sz == 0 {
		return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			return streamer(ctx, desc, cc, method, opts...)
		}
	}
	if sz == 1 {
		return interceptors[0]
	}
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		curr := streamer
		for i := sz - 1; i > 0; i-- {
			temp, i := curr, i
			curr = func(myctx context.Context, mydesc *grpc.StreamDesc, mycc *grpc.ClientConn, mymethod string, myopts ...grpc.CallOption) (grpc.ClientStream, error) {
				return interceptors[i](myctx, mydesc, mycc, mymethod, temp, opts...)
			}
		}
		return interceptors[0](ctx, desc, cc, method, streamer, opts...)
	}
}

func WithUnaryServerChain(interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(interceptors...)
}

func WithStreamServerChain(interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	return grpc.ChainStreamInterceptor(interceptors...)
}
