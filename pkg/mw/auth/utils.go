package auth

import (
	"context"

	"github.com/tinysrc/z9go/pkg/mw/tags"
	"github.com/tinysrc/z9go/pkg/z9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// JwtAuth impl
func JwtAuth(ctx context.Context, sign string) (context.Context, error) {
	token, err := AuthFromMD(ctx, "Basic")
	if err != nil {
		return nil, err
	}
	j := z9.NewJWT(sign)
	claims, err := j.ParseToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token error=%v", err)
	}
	tags.Extract(ctx).Set("userid", claims.Userid)
	return ctx, nil
}
