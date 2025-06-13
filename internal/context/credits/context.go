package credits

import (
	"context"
	"fmt"
)

type contextKey string

const creditsKey = contextKey("remainingCredits")

func NewContextWithCredits(ctx context.Context, remainingCredits int) context.Context {
	return context.WithValue(ctx, creditsKey, remainingCredits)
}

func creditsFromContext(ctx context.Context) (int, bool) {
	val, ok := ctx.Value(creditsKey).(int)
	if !ok {
		return 0, false
	}

	return val, true
}

func RemainingCredits(ctx context.Context) (int, error) {
	remainingCredits, ok := creditsFromContext(ctx)
	if !ok {
		return 0, fmt.Errorf("credits not found")
	}

	return remainingCredits, nil
}
