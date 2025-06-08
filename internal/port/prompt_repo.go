package port

import "github.com/CP-Payne/wonderpicai/internal/domain"

type PromptRepository interface {
	Create(prompt *domain.Prompt) error
}
