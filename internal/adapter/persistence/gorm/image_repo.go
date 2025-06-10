package gorm

import (
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

func (r *gormImageRepository) GetByID(imageID uuid.UUID) (*domain.Image, error) {
	var image domain.Image
	result := r.db.Where("id = ?", imageID).First(&image)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrImageNotFound
		}
		r.logger.Error("Failed to get image by ID", zap.String("imageID", imageID.String()), zap.Error(result.Error))
		return nil, fmt.Errorf("database error fetching image by id: %w", result.Error)
	}
	return &image, nil
}

func (r *gormImageRepository) Delete(imageID uuid.UUID) error {
	err := r.db.Delete(&domain.Image{}, imageID).Error
	if err != nil {
		r.logger.Error("failed to delete image from repository")
		return fmt.Errorf("failed to delete image")
	}

	return nil
}

func (r *gormImageRepository) DeleteFailed() error {
	err := r.db.Where("status=?", domain.Failed).Delete(&domain.Image{}).Error
	if err != nil {
		r.logger.Error("failed to delete images with status failed from repository")
		return fmt.Errorf("failed to delete images with status failed")
	}

	return nil
}
