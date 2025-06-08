package gorm

import (
	"fmt"

	"github.com/CP-Payne/wonderpicai/internal/domain"
	"github.com/CP-Payne/wonderpicai/internal/port"
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

func (r *gormImageRepository) Create(image *domain.Image) error {
	result := r.db.Create(image)
	if result.Error != nil {
		r.logger.Error("Failed to create image in database", zap.Error(result.Error))
		return fmt.Errorf("failed to create image in database: %w", result.Error)
	}
	return nil
}
