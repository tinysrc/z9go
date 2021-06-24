package utils

import (
	"context"

	"google.golang.org/grpc"
)

type WrappedServerStream struct {
	grpc.ServerStream
	WrappedCtx context.Context
}

func (w *WrappedServerStream) Context() context.Context {
	return w.WrappedCtx
}

func WrapServerStream(s grpc.ServerStream) *WrappedServerStream {
	if w, ok := s.(*WrappedServerStream); ok {
		return w
	}
	return &WrappedServerStream{
		ServerStream: s,
		WrappedCtx:   s.Context(),
	}
}
