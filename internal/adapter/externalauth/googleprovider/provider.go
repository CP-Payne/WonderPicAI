package googleprovider

import (
	"fmt"
	"net/http"

	"github.com/CP-Payne/wonderpicai/internal/domain"
	"github.com/CP-Payne/wonderpicai/internal/port"
	"go.uber.org/zap"
	"google.golang.org/api/idtoken"
)

type GoogleAuthProvider struct {
	logger   *zap.Logger
	clientID string
}

func NewAuth(logger *zap.Logger, clientID string) port.ExternalAuthService {
	return &GoogleAuthProvider{
		logger:   logger,
		clientID: clientID,
	}
}

func (p *GoogleAuthProvider) HandleCallback(r *http.Request) (*port.ExternalUserData, error) {
	idToken := r.FormValue("credential")
	if idToken == "" {
		p.logger.Warn("idToken not found in form")
		return nil, fmt.Errorf("idToken missing from form: %w", domain.ErrInvalidCredentials)
	}

	payload, err := idtoken.Validate(r.Context(), idToken, p.clientID)
	if err != nil {
		p.logger.Warn("ID token validation failed", zap.Error(err))
		return nil, fmt.Errorf("idToken validation failed: %w", err)
	}

	exUserData := &port.ExternalUserData{
		Email: payload.Claims["email"].(string),
		Name:  payload.Claims["name"].(string),
	}

	p.logger.Info("User authentication via Google Auth", zap.String("email", exUserData.Email))

	return exUserData, nil
}
