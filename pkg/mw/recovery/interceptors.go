package recovery

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RecoveryHandlerFunc func(p interface{}) (err error)
type RecoveryHandlerContextFunc func(ctx context.Context, p interface{}) (err error)

func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	o := evalOptions(opts)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		panic := true
		defer func() {
			if r := recover(); r != nil || panic {
				recoverFrom(ctx, r, o.recoveryHandlerFunc)
			}
		}()
		resp, err := handler(ctx, req)
		panic = false
		return resp, err
	}
}

func StreamServerInterceptor(opts ...Option) grpc.StreamServerInterceptor {
	o := evalOptions(opts)
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		panic := true
		defer func() {
			if r := recover(); r != nil || panic {
				recoverFrom(ss.Context(), r, o.recoveryHandlerFunc)
			}
		}()
		err := handler(srv, ss)
		panic = false
		return err
	}
}

func recoverFrom(ctx context.Context, p interface{}, r RecoveryHandlerContextFunc) error {
	if r == nil {
		return status.Errorf(codes.Internal, "%v", p)
	}
	return r(ctx, p)
}
