package port

import (
	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	GenerateToken(claims jwt.Claims) (string, error)
	// TODO: Decouple from jwt.Token by returning a custom claims struct (e.g., *CustomClaims).
	// This will make the port more independent of the underlying JWT library.
	// Will revisit once all necessary claims for middleware/services are finalized.
	ValidateToken(token string) (*jwt.Token, error)
}
