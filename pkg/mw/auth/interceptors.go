package auth

import (
	"context"

	"github.com/tinysrc/z9go/pkg/mw/utils"
	"google.golang.org/grpc"
)

type AuthFunc func(ctx context.Context) (context.Context, error)

type serviceAuth interface {
	AuthFunc(ctx context.Context, fullMethodName string) (context.Context, error)
}

func UnaryServerInterceptor(authFunc AuthFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var newCtx context.Context
		var err error
		if sa, ok := info.Server.(serviceAuth); ok {
			newCtx, err = sa.AuthFunc(ctx, info.FullMethod)
		} else {
			newCtx, err = authFunc(ctx)
		}
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

func StreamServerInterceptor(authFunc AuthFunc) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var newCtx context.Context
		var err error
		if sa, ok := srv.(serviceAuth); ok {
			newCtx, err = sa.AuthFunc(ss.Context(), info.FullMethod)
		} else {
			newCtx, err = authFunc(ss.Context())
		}
		if err != nil {
			return err
		}
		w := utils.WrapServerStream(ss)
		w.WrappedCtx = newCtx
		return handler(srv, w)
	}
}

func UnaryClientInterceptor(authFunc AuthFunc) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx, err := authFunc(ctx)
		if err != nil {
			return err
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func StreamClientInterceptor(authFunc AuthFunc) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx, err := authFunc(ctx)
		if err != nil {
			return nil, err
		}
		s, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			return nil, err
		}
		w := utils.WrapClientStream(s)
		w.WrappedCtx = ctx
		return w, nil
	}
}
