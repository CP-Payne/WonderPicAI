package gorm

import (
	"context"
	"fmt"
	"time"

	"github.com/CP-Payne/wonderpicai/internal/domain"
	"github.com/CP-Payne/wonderpicai/internal/port"
	"github.com/google/uuid"
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

func (r *gormPromptRepository) Create(ctx context.Context, prompt *domain.Prompt) (*domain.Prompt, error) {

	var images []domain.Image

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&prompt).Error; err != nil {
			return err
		}

		images = make([]domain.Image, prompt.ImageCount)
		for i := 0; i < prompt.ImageCount; i++ {
			images[i] = domain.Image{
				BaseModel: domain.BaseModel{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				PromptID:  prompt.ID,
				ImageData: nil,
				Status:    domain.Pending,
			}
		}
		if err := tx.Create(&images).Error; err != nil {
			return err
		}

		return nil

	})

	if err != nil {
		r.logger.Error("Failed to create prompt and placeholder images in database", zap.Error(err))
		return nil, fmt.Errorf("failed to create prompt and placeholder images in database: %w", err)
	}

	prompt.Images = images

	return prompt, nil
}
