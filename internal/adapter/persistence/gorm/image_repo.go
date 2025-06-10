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
