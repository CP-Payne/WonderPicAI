package gorm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/CP-Payne/wonderpicai/internal/domain"
	"github.com/CP-Payne/wonderpicai/internal/port"
	"github.com/google/uuid"
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

	tx := r.db.Begin()

	defer func() {
		if rcv := recover(); rcv != nil {
			tx.Rollback()
			panic(rcv)
		}
	}()

	if err := tx.Error; err != nil {
		r.logger.Error("Failed to begin transaction", zap.Error(err))
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// --- Note ---
	// This method violates the Single Responsibility Principle by also handling Wallet creation. In a larger, production system, this would be implemented using a Unit Of Work pattern

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		r.logger.Debug("GORM create user failed within transaction", zap.Error(err))

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				r.logger.Info("PostgreSQL unique violation on user creation",
					zap.String("constraint", pgErr.ConstraintName),
					zap.String("detail", pgErr.Detail))

				if strings.Contains(pgErr.ConstraintName, "email") {
					return domain.ErrEmailAlreadyExists
				}
				return fmt.Errorf("%w: a unique field caused a conflict", domain.ErrDuplicateEntry)
			}
		}

		r.logger.Error("Failed to create user in database", zap.Error(err))
		return fmt.Errorf("failed to create user in database: %w", err)
	}

	wallet := &domain.Wallet{
		BaseModel: domain.BaseModel{ID: uuid.New()},
		UserID:    user.ID,
	}

	if err := tx.Create(wallet).Error; err != nil {
		tx.Rollback()
		r.logger.Error("Failed to create associated wallet for user",
			zap.String("userID", user.ID.String()),
			zap.Error(err))
		return fmt.Errorf("failed to create wallet for user %d: %w", user.ID, err)
	}

	if err := tx.Commit().Error; err != nil {
		r.logger.Error("Failed to commit transaction for user creation", zap.Error(err))
		return fmt.Errorf("failed to commit transaction: %w", err)
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

func (r *gormUserRepository) GetByID(userID uuid.UUID) (*domain.User, error) {

	var user domain.User

	result := r.db.First(&user, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		r.logger.Error("Failed to get user by ID", zap.String("userID", userID.String()), zap.Error(result.Error))
		return nil, fmt.Errorf("database error fetching user by ID: %w", result.Error)
	}
	return &user, nil
}
