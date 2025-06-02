package gorm

import (
	"errors"

	"github.com/CP-Payne/wonderpicai/internal/domain"
	"github.com/CP-Payne/wonderpicai/internal/port"
	"gorm.io/gorm"
)

type gormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) port.UserRepository {
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) Create(user *domain.User) error {
	result := r.db.Create(user)
	return result.Error
}

func (r *gormUserRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User

	result := r.db.Where("email= ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}
