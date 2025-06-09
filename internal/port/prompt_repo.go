package port

import (
	"context"

	"github.com/CP-Payne/wonderpicai/internal/domain"
)

type PromptRepository interface {
	Create(ctx context.Context, prompt *domain.Prompt) (*domain.Prompt, error)
}
