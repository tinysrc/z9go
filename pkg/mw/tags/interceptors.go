package tags

import (
	"context"

	"github.com/tinysrc/z9go/pkg/mw/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	o := evalOptions(opts)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		newCtx := newTagsForCtx(ctx)
		if o.requestFieldsFunc != nil {
			setRequestFieldTags(newCtx, o.requestFieldsFunc, info.FullMethod, req)
		}
		return handler(newCtx, req)
	}
}

func StreamServerInterceptor(opts ...Option) grpc.StreamServerInterceptor {
	o := evalOptions(opts)
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		newCtx := newTagsForCtx(ss.Context())
		if o.requestFieldsFunc == nil {
			w := utils.WrapServerStream(ss)
			w.WrappedCtx = newCtx
			return handler(srv, w)
		}
		w := &wrappedStream{ss, info, o, newCtx, true}
		return handler(srv, w)
	}
}

type wrappedStream struct {
	grpc.ServerStream
	info       *grpc.StreamServerInfo
	opts       *options
	WrappedCtx context.Context
	init       bool
}

func (w *wrappedStream) Context() context.Context {
	return w.WrappedCtx
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	err := w.ServerStream.RecvMsg(m)
	if !w.info.IsClientStream || w.opts.requestFieldsFromInit && w.init {
		w.init = false
		setRequestFieldTags(w.Context(), w.opts.requestFieldsFunc, w.info.FullMethod, m)
	}
	return err
}

func newTagsForCtx(ctx context.Context) context.Context {
	tags := NewTags()
	if peer, ok := peer.FromContext(ctx); ok {
		tags.Set("peer.addr", peer.Addr.String())
	}
	return SetInContext(ctx, tags)
}

func setRequestFieldTags(ctx context.Context, f RequestFieldExtractorFunc, fullMethodName string, req interface{}) {
	if res := f(fullMethodName, req); res != nil {
		t := Extract(ctx)
		for k, v := range res {
			t.Set("rpc.req."+k, v)
		}
	}
}
