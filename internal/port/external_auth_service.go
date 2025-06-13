package port

import "net/http"

type ExternalUserData struct {
	Email string
	Name  string
}

type ExternalAuthService interface {
	HandleCallback(r *http.Request) (*ExternalUserData, error)
}
