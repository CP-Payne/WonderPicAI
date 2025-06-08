package gorm

import (
	"fmt"

	"github.com/CP-Payne/wonderpicai/internal/domain"
	"github.com/CP-Payne/wonderpicai/internal/port"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type gormPromptRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewGormPromptRepository(db *gorm.DB, logger *zap.Logger) port.PromptRepository {
	return &gormPromptRepository{db: db, logger: logger.With(zap.String("component", "PromptRepoGORM"))}
}

func (r *gormPromptRepository) Create(prompt *domain.Prompt) error {
	result := r.db.Create(prompt)
	if result.Error != nil {
		r.logger.Error("Failed to create prompt in database", zap.Error(result.Error))
		return fmt.Errorf("failed to create prompt in database: %w", result.Error)
	}
	return nil
}
