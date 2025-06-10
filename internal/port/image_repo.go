package port

import (
	"github.com/CP-Payne/wonderpicai/internal/domain"
	"github.com/google/uuid"
)

type ImageRepository interface {
	GetByID(imageID uuid.UUID) (*domain.Image, error)
}
