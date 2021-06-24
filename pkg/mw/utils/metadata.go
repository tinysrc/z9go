package utils

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

type WrappedMD metadata.MD

func ExtractIncoming(ctx context.Context) WrappedMD {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return WrappedMD(metadata.Pairs())
	}
	return WrappedMD(md)
}

func ExtractOutgoing(ctx context.Context) WrappedMD {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return WrappedMD(metadata.Pairs())
	}
	return WrappedMD(md)
}

func (w WrappedMD) Clone(keys ...string) WrappedMD {
	md := WrappedMD(metadata.Pairs())
	for k, v := range w {
		found := false
		if len(keys) == 0 {
			found = true
		} else {
			for _, key := range keys {
				if strings.EqualFold(key, k) {
					found = true
					break
				}
			}
		}
		if !found {
			continue
		}
		md[k] = make([]string, len(v))
		copy(md[k], v)
	}
	return md
}

func (w WrappedMD) ToIncoming(ctx context.Context) context.Context {
	return metadata.NewIncomingContext(ctx, metadata.MD(w))
}

func (w WrappedMD) ToOutgoing(ctx context.Context) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.MD(w))
}

func (w WrappedMD) Get(key string) string {
	key = strings.ToLower(key)
	values, ok := w[key]
	if ok {
		return values[0]
	}
	return ""
}

func (w WrappedMD) Set(key, value string) WrappedMD {
	key = strings.ToLower(key)
	w[key] = []string{value}
	return w
}

func (w WrappedMD) Add(key, value string) WrappedMD {
	key = strings.ToLower(key)
	w[key] = append(w[key], value)
	return w
}

func (w WrappedMD) Del(key string) WrappedMD {
	key = strings.ToLower(key)
	delete(w, key)
	return w
}
