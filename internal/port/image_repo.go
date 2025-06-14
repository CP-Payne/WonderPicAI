package port

import (
	"context"

	"github.com/CP-Payne/wonderpicai/internal/domain"
	"github.com/google/uuid"
)

type ImageRepository interface {
	GetByID(ctx context.Context, userID uuid.UUID, imageID uuid.UUID) (*domain.Image, error)
	Delete(ctx context.Context, userID uuid.UUID, imageID uuid.UUID) error
	DeleteFailed(ctx context.Context, userID uuid.UUID) error
	ContainsFailedImages(ctx context.Context, userID uuid.UUID) (bool, error)
}
