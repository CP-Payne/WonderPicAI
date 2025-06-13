package domain

import "errors"

var (
	ErrUserNotFound            = errors.New("user not found")
	ErrEmailAlreadyExists      = errors.New("email already exists")
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrDuplicateEntry          = errors.New("entry with unique field already exists")
	ErrImageNotFound           = errors.New("image not found")
	ErrRecordNotFound          = errors.New("record not found")
	ErrInsufficientFunds       = errors.New("insufficient funds")
	ErrInvalidPurchaseOption   = errors.New("invalid purchase option")

	ErrUnhandledEvent = errors.New("unhandled event")
)
