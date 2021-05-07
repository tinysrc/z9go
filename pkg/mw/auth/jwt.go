package auth

import (
	"context"
	"errors"

	"github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/tinysrc/z9go/pkg/conf"
	"github.com/tinysrc/z9go/pkg/mw/tags"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errTokenExpired     = errors.New("token expired")
	errTokenNotValidYet = errors.New("token not valid yet")
	errTokenMalformed   = errors.New("token malformed")
	errTokenInvalid     = errors.New("token invalid")
)

// CustomClaims struct
type CustomClaims struct {
	UUID string
	jwt.StandardClaims
}

// JWT struct
type JWT struct {
	Sign []byte
}

// NewJWT impl
func NewJWT() *JWT {
	sign := conf.Global.GetString("service.jwt.sign")
	return &JWT{[]byte(sign)}
}

// MakeToken impl
func (j *JWT) MakeToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.Sign)
}

// ParseToken impl
func (j *JWT) ParseToken(ts string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(ts, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.Sign, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errTokenNotValidYet
			} else if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errTokenMalformed
			} else {
				return nil, errTokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
	}
	return nil, errTokenInvalid
}

// AuthFunc impl
func AuthFunc(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "Basic")
	if err != nil {
		return nil, err
	}
	j := NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token error=%v", err)
	}
	tags.Extract(ctx).Set("userid", claims.UUID)
	return ctx, nil
}
