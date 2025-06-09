package port

import (
	"context"

	"github.com/CP-Payne/wonderpicai/internal/domain"
	"github.com/google/uuid"
)

type PromptRepository interface {
	Create(ctx context.Context, prompt *domain.Prompt) (*domain.Prompt, error)
	FindAllByUser(ctx context.Context, userID uuid.UUID) ([]domain.Prompt, error)
}
