package xjwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type Jwt interface {
	Parse(token string) (jwt.MapClaims, error)
}
