package gorm

import (
	"context"
	"fmt"

	"github.com/CP-Payne/wonderpicai/internal/domain"
	"github.com/CP-Payne/wonderpicai/internal/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type gormWalletRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewGormWalletRepository(db *gorm.DB, logger *zap.Logger) port.WalletRepository {
	return &gormWalletRepository{db: db, logger: logger.With(zap.String("component", "walletRepoGORM"))}
}

func (r *gormWalletRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Wallet, error) {

	var wallet domain.Wallet

	err := r.db.WithContext(ctx).
		Select("*").
		Where("user_id = ?", userID).
		First(&wallet).Error

	if err != nil {
		r.logger.Error("Failed to get wallet by user ID due to database error", zap.String("userID", userID.String()), zap.Error(err))

		return nil, fmt.Errorf("database error fetching wallet: %w", err)
	}

	return &wallet, nil

}

func (r *gormWalletRepository) SubtractCredits(ctx context.Context, userID uuid.UUID, amount int) error {

	expression := gorm.Expr("credits - ?", amount)

	result := r.db.Model(&domain.Wallet{}).
		Where("user_id = ? AND credits >= ?", userID, amount).
		Update("credits", expression)

	if result.Error != nil {
		r.logger.Error("Database error during credit subtraction",
			zap.String("userID", userID.String()),
			zap.Int("amount", amount),
			zap.Error(result.Error))
		return fmt.Errorf("db error subtracting credits: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		r.logger.Warn("Failed to subtract credits: insufficient funds or user not found",
			zap.String("userID", userID.String()),
			zap.Int("amount", amount))
		return domain.ErrInsufficientFunds
	}

	return nil

}

func (r *gormWalletRepository) AddCredits(ctx context.Context, userID uuid.UUID, amount int) error {
	expression := gorm.Expr("credits + ?", amount)

	result := r.db.Model(&domain.Wallet{}).
		Where("user_id = ?", userID).
		Update("credits", expression)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		// This would mean the user_id doesn't exist
		return domain.ErrRecordNotFound
	}

	return nil
}
