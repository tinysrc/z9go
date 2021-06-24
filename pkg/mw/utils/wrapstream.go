package utils

import (
	"context"

	"google.golang.org/grpc"
)

type WrappedClientStream struct {
	grpc.ClientStream
	WrappedCtx context.Context
}

func (w *WrappedClientStream) Context() context.Context {
	return w.WrappedCtx
}

func WrapClientStream(s grpc.ClientStream) *WrappedClientStream {
	if w, ok := s.(*WrappedClientStream); ok {
		return w
	}
	return &WrappedClientStream{
		ClientStream: s,
		WrappedCtx:   s.Context(),
	}
}

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
