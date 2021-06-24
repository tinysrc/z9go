package auth

import (
	"context"
	"strings"

	"github.com/tinysrc/z9go/pkg/mw/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthFromMD impl
func AuthFromMD(ctx context.Context, scheme string) (string, error) {
	val := utils.ExtractIncoming(ctx).Get("Authorization")
	if val == "" {
		return "", status.Errorf(codes.Unauthenticated, "Request unauthenticated scheme=%s", scheme)
	}
	splits := strings.SplitN(val, " ", 2)
	if len(splits) < 2 {
		return "", status.Errorf(codes.Unauthenticated, "Bad authorization string")
	}
	if !strings.EqualFold(splits[0], scheme) {
		return "", status.Errorf(codes.Unauthenticated, "Request unauthenticated scheme=%s", scheme)
	}
	return splits[1], nil
}
