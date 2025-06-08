package port

import "github.com/CP-Payne/wonderpicai/internal/domain"

type ImageRepository interface {
	Create(image *domain.Image) error
}
