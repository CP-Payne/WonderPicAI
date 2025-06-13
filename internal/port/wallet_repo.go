package port

import (
	"context"

	"github.com/CP-Payne/wonderpicai/internal/domain"
	"github.com/google/uuid"
)

type WalletRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Wallet, error)
	SubtractCredits(ctx context.Context, userID uuid.UUID, amount int) error
	AddCredits(ctx context.Context, userID uuid.UUID, amount int) error
	AddCreditsToEmail(ctx context.Context, email string, amount int) error
}
