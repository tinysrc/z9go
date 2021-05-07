package tags

import (
	"context"

	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
)

func Extract(ctx context.Context) grpc_ctxtags.Tags {
	return grpc_ctxtags.Extract(ctx)
}

func SetInContext(ctx context.Context, tags grpc_ctxtags.Tags) context.Context {
	return grpc_ctxtags.SetInContext(ctx, tags)
}

func NewTags() grpc_ctxtags.Tags {
	return grpc_ctxtags.NewTags()
}
