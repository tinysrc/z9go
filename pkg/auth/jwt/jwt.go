package jwt

import (
	"context"
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/tinysrc/z9go/pkg/conf"
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
	UUID uuid.UUID
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

// ToToken impl
func (j *JWT) ToToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.Sign)
}

// FromToken impl
func (j *JWT) FromToken(ts string) (*CustomClaims, error) {
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

// Auth impl
func Auth(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	j := NewJWT()
	_, err = j.FromToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token error=%v", err)
	}
	return ctx, nil
}