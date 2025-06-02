package port

import "github.com/CP-Payne/wonderpicai/internal/domain"

type UserRepository interface {
	Create(user *domain.User) error
	GetByEmail(email string) (*domain.User, error)
}
