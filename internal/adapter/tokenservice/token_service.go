package tokenservice

import (
	"fmt"

	"github.com/CP-Payne/wonderpicai/internal/port"
	"github.com/golang-jwt/jwt/v5"
)

type tokenService struct {
	secret string
	iss    string
}

func NewTokenService(secret, iss string) port.TokenService {
	return &tokenService{
		secret: secret,
		iss:    iss,
	}
}

func (ts *tokenService) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(ts.secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (ts *tokenService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(ts.secret), nil
	}, jwt.WithExpirationRequired(),
		// Audience and issuer in this implementation are the same
		jwt.WithAudience(ts.iss),
		jwt.WithIssuer(ts.iss),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
}
