package port

import (
	"github.com/CP-Payne/wonderpicai/internal/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetByEmail(email string) (*domain.User, error)
	GetByID(userID uuid.UUID) (*domain.User, error)
}
