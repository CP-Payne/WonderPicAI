package port

import "github.com/CP-Payne/wonderpicai/internal/domain"

type UserRepository interface {
	Create(user *domain.User) error
	GetByUsername(username string) (*domain.User, error)
}
