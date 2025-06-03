package domain

import "errors"

var (
	ErrUserNotFound            = errors.New("user not found")
	ErrEmailAlreadyExists      = errors.New("email already exists")
	ErrInvalidCredentials      = errors.New("invalid credentials") // TODO: Consider moving to service layer
	ErrInsufficientPermissions = errors.New("insufficient perfmissions")
)
