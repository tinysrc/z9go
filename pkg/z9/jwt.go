package z9

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

var (
	errTokenExpired     = errors.New("token expired")
	errTokenNotValidYet = errors.New("token not valid yet")
	errTokenMalformed   = errors.New("token malformed")
	errTokenInvalid     = errors.New("token invalid")
)

// CustomClaims struct
type CustomClaims struct {
	Userid string
	Orgid  string
	jwt.StandardClaims
}

// JWT struct
type JWT struct {
	Sign []byte
}

// NewJWT impl
func NewJWT(sign string) *JWT {
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
