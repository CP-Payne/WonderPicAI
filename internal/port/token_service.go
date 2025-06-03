package port

import (
	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	GenerateToken(claims jwt.Claims) (string, error)
	ValidateToken(string) (*jwt.Token, error)
}
