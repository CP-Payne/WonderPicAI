package gorm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/CP-Payne/wonderpicai/internal/domain"
	"github.com/CP-Payne/wonderpicai/internal/port"
	"github.com/jackc/pgconn"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type gormUserRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewGormUserRepository(db *gorm.DB, logger *zap.Logger) port.UserRepository {
	return &gormUserRepository{db: db, logger: logger.With(zap.String("component", "UserRepoGORM"))}
}

func (r *gormUserRepository) Create(user *domain.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		r.logger.Debug("GORM create user failed", zap.Error(result.Error))

		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) {
			if pgErr.Code == "23505" {
				r.logger.Info("PostgreSQL unique violation",
					zap.String("constraint", pgErr.ConstraintName),
					zap.String("detail", pgErr.Detail))

				if strings.Contains(pgErr.ConstraintName, "email") {
					return domain.ErrEmailAlreadyExists
				}

				return fmt.Errorf("%w: a unique field caused a conflict", domain.ErrDuplicateEntry)
			}
		}

		r.logger.Error("Failed to create user in database", zap.Error(result.Error))
		return fmt.Errorf("failed to create user in database: %w", result.Error)
	}
	return nil
}

func (r *gormUserRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User

	result := r.db.Where("email= ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		r.logger.Error("Failed to get user by email", zap.String("email", email), zap.Error(result.Error))
		return nil, fmt.Errorf("database error fetching user by email: %w", result.Error)
	}
	return &user, nil
}
