package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type contextKey string

const userIDKey = contextKey("userID")

var ErrUserNotAuthenticated = errors.New("no user ID found in context")

func NewContextWithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func userIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	val, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, false
	}

	if val == uuid.Nil {
		return uuid.Nil, false
	}

	return val, true
}

func UserID(ctx context.Context) (uuid.UUID, error) {
	userID, ok := userIDFromContext(ctx)
	if !ok {
		return uuid.Nil, ErrUserNotAuthenticated
	}

	return userID, nil
}

func IsAuthenticated(ctx context.Context) bool {
	_, ok := userIDFromContext(ctx)
	return ok
}
