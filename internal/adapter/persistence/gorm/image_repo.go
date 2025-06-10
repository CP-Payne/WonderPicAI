package gorm

import (
	"context"
	"errors"
	"fmt"

	"github.com/CP-Payne/wonderpicai/internal/domain"
	"github.com/CP-Payne/wonderpicai/internal/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type gormImageRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewGormImageRepository(db *gorm.DB, logger *zap.Logger) port.ImageRepository {
	return &gormImageRepository{db: db, logger: logger.With(zap.String("component", "ImageRepoGORM"))}
}

func (r *gormImageRepository) GetByID(ctx context.Context, userID uuid.UUID, imageID uuid.UUID) (*domain.Image, error) {
	var image domain.Image

	err := r.db.WithContext(ctx).
		Select("images.*").
		Joins("JOIN prompts ON prompts.id = images.prompt_id").
		Where("images.id = ?", imageID).
		Where("prompts.user_id = ?", userID).
		First(&image).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Warn("Record not found for user", zap.String("imageID", imageID.String()), zap.String("userID", userID.String()))
			return nil, domain.ErrImageNotFound
		}
		r.logger.Error("Failed to get image by ID due to database error", zap.String("imageID", imageID.String()), zap.String("userID", userID.String()), zap.Error(err))
		return nil, fmt.Errorf("database error fetching image: %w", err)
	}
	return &image, nil
}

func (r *gormImageRepository) Delete(ctx context.Context, userID uuid.UUID, imageID uuid.UUID) error {

	subQuery := r.db.Model(&domain.Prompt{}).Select("id").Where("user_id = ?", userID)

	result := r.db.WithContext(ctx).
		Where("id = ?", imageID).
		Where("prompt_id IN (?)", subQuery).
		Delete(&domain.Image{})
	if result.Error != nil {
		r.logger.Error("failed to execute delete query", zap.Error(result.Error))
		return fmt.Errorf("database error while deleting image")
	}

	if result.RowsAffected == 0 {
		r.logger.Warn("delete operation affected 0 rows", zap.String("imageID", imageID.String()), zap.String("userID", userID.String()))
		return domain.ErrRecordNotFound
	}

	r.logger.Info("successfully deleted image", zap.String("imageID", imageID.String()), zap.String("userID", userID.String()))

	return nil
}

func (r *gormImageRepository) DeleteFailed(ctx context.Context, userID uuid.UUID) error {

	subQuery := r.db.Model(&domain.Prompt{}).Select("id").Where("user_id = ?", userID)

	result := r.db.WithContext(ctx).
		Where("status = ?", domain.Failed).
		Where("prompt_id IN (?)", subQuery). // Use the subquery here
		Delete(&domain.Image{})

	if result.Error != nil {
		r.logger.Error("Failed to execute delete query for failed images",
			zap.String("userID", userID.String()),
			zap.Error(result.Error),
		)
		return fmt.Errorf("database error while deleting failed images")
	}

	r.logger.Info("Successfully processed request to delete failed images",
		zap.String("userID", userID.String()),
		zap.Int64("images_deleted", result.RowsAffected),
	)

	return nil
}
