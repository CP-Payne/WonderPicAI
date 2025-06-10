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

func (r *gormPromptRepository) FindAllByUser(ctx context.Context, userID uuid.UUID) ([]domain.Prompt, error) {
	var prompts []domain.Prompt

	err := r.db.WithContext(ctx).Preload("Images").Order("created_at desc").Where("user_id = ?", userID).Find(&prompts).Error
	if err != nil {
		return nil, fmt.Errorf("failed retrieving prompts from repo: %w", err)
	}

	return prompts, nil
}

func (r *gormPromptRepository) UpdatePlaceholderImages(ctx context.Context, externalPromptID uuid.UUID, images [][]byte) (*domain.Prompt, error) {
	var prompt domain.Prompt
	var finalErr error

	txErr := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("external_prompt_id = ?", externalPromptID).First(&prompt).Error; err != nil {
			r.logger.Warn("External prompt ID does not exist in database", zap.String("externalPromptID", externalPromptID.String()))
			return err
		}

		var placeholderImages []domain.Image
		if err := tx.Where("prompt_id = ? AND status = ?", prompt.ID, domain.Pending).Find(&placeholderImages).Error; err != nil {
			finalErr = fmt.Errorf("failed to find placeholder images: %w", err)
			prompt.Status = domain.Failed
		}

		if finalErr == nil {
			if len(placeholderImages) != len(images) {
				finalErr = fmt.Errorf("mismatch: expected %d images, but %d was provided", len(placeholderImages), len(images))
				prompt.Status = domain.Failed
			} else {
				for i, imageData := range images {
					imageToUpdate := &placeholderImages[i]
					imageToUpdate.ImageData = imageData
					imageToUpdate.UpdatedAt = time.Now()
					imageToUpdate.Status = domain.Completed

					if err := tx.Save(imageToUpdate).Error; err != nil {
						finalErr = fmt.Errorf("failed to save image index %d: %w", i, err)
						prompt.Status = domain.Failed
						break
					}
				}
			}
		}

		if finalErr == nil {
			prompt.Status = domain.Completed
			prompt.Images = placeholderImages
		}

		if prompt.Status == domain.Failed {
			if err := tx.Model(&domain.Image{}).Where("prompt_id = ?", prompt.ID).Update("status", domain.Failed).Error; err != nil {
				r.logger.Error("CRITICAL: Failed to update images to FAILED status", zap.Error(err))
				return err
			}
		}

		if err := tx.Save(&prompt).Error; err != nil {
			r.logger.Error("CRITICAL: Failed to save final prompt status", zap.Error(err))
			return err // Abort transaction!
		}

		return nil
	})

	if txErr != nil {
		r.logger.Error("Transaction rolled back due to fatal error", zap.Error(txErr))
		return nil, fmt.Errorf("database transaction failed: %w", txErr)
	}

	if finalErr != nil {
		r.logger.Error("Image update operation failed", zap.Error(finalErr))
		return nil, finalErr
	}

	// Success!
	return &prompt, nil
}
