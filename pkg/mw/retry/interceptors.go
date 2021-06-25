package retry

import (
	"context"
	"strconv"

	"github.com/tinysrc/z9go/pkg/mw/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const AttemptMetadataKey = "z9-retry-attempt"

func UnaryClientInterceptor(cos ...CallOption) grpc.UnaryClientInterceptor {
	intOpts := reuseOrNewWithCallOptions(defaultOptions, cos)
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		grpcOpts, retryOpts := filterCallOptions(opts)
		callOpts := reuseOrNewWithCallOptions(intOpts, retryOpts)
		if callOpts.max == 0 {
			return invoker(ctx, method, req, reply, cc, grpcOpts...)
		}
		var err error
		for i := uint(0); i < callOpts.max; i++ {
			if err = waitRetryBackoff(i, ctx, callOpts); err != nil {
				return err
			}
			newCtx := ctx
			var cancel context.CancelFunc
			if i > 0 && callOpts.incHeader {
				md := utils.ExtractOutgoing(ctx).Clone().Set(AttemptMetadataKey, strconv.Itoa(int(i)))
				newCtx = md.ToOutgoing(newCtx)
			}
			if callOpts.timeout != 0 {
				newCtx, cancel = context.WithTimeout(newCtx, callOpts.timeout)
			}
			defer func() {
				if cancel != nil {
					cancel()
				}
			}()
			err = invoker(newCtx, method, req, reply, cc, grpcOpts...)
			if err == nil {
				return nil
			}
			logTrace(ctx, "grpc retry attempt=%d error=%v", i, err)
			if isCtxErr(err) {
				if ctx.Err() != nil {
					logTrace(ctx, "grpc retry attempt=%d context error=%v", i, ctx.Err())
					return err
				} else if callOpts.timeout != 0 {
					logTrace(ctx, "grpc retry attempt=%d context error from retry call", i)
					continue
				}
			}
			if !isRetriable(err, callOpts) {
				return err
			}
		}
		return err
	}
}

func StreamClientInterceptor(cos ...CallOption) grpc.StreamClientInterceptor {
	intOpts := reuseOrNewWithCallOptions(defaultOptions, cos)
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		grpcOpts, retryOpts := filterCallOptions(opts)
		callOpts := reuseOrNewWithCallOptions(intOpts, retryOpts)
		if callOpts.max == 0 {
			return streamer(ctx, desc, cc, method, grpcOpts...)
		}
		if desc.ClientStreams {
			return nil, status.Errorf(codes.Unimplemented, "grpc retry can't work on ClientStreams, set retry.Disable()")
		}
		var err error
		for i := uint(0); i < callOpts.max; i++ {
			if err = waitRetryBackoff(i, ctx, callOpts); err != nil {
				return nil, err
			}
			newCtx := ctx
			var cancel context.CancelFunc
			if callOpts.timeout != 0 {
				newCtx, cancel = context.WithTimeout(newCtx, callOpts.timeout)
			}
			defer func() {
				if cancel != nil {
					cancel()
				}
			}()
			var cs grpc.ClientStream
			cs, err = streamer(newCtx, desc, cc, method, grpcOpts...)
			if err == nil {
				s := &wrappedClientStream{
					ClientStream: cs,
					callOpts:     callOpts,
					parentCtx:    ctx,
					streamer: func(myctx context.Context) (grpc.ClientStream, error) {
						return streamer(myctx, desc, cc, method, grpcOpts...)
					},
				}
				return s, nil
			}
			logTrace(ctx, "grpc retry attempt=%d error=%v", i, err)
			if isCtxErr(err) {
				if ctx.Err() != nil {
					logTrace(ctx, "grpc retry attempt=%d context error=%v", i, ctx.Err())
					return nil, err
				} else if callOpts.timeout != 0 {
					logTrace(ctx, "grpc retry attempt=%d context error from retry call", i)
					continue
				}
			}
			if !isRetriable(err, callOpts) {
				return nil, err
			}
		}
		return nil, err
	}
}
