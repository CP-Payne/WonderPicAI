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

// func (r *gormPromptRepository) UpdatePlaceholderImages(ctx context.Context, externalPromptID uuid.UUID, images [][]byte) (*domain.Prompt, error) {

// 	var prompt domain.Prompt

// 	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		if err := tx.Where("external_prompt_id = ?", externalPromptID).First(&prompt).Error; err != nil {
// 			// No such prompt
// 			r.logger.Warn("external prompt ID does not exist in database", zap.String("externalPromptID", externalPromptID.String()))
// 			return err
// 		}

// 		var placeholderImages []domain.Image
// 		if err := tx.Where("prompt_id = ? AND status = ?", prompt.ID, domain.Pending).Find(&placeholderImages).Error; err != nil {
// 			return err
// 		}

// 		if len(placeholderImages) != len(images) {
// 			prompt.Status = domain.Failed
// 			// Save prompt first before exiting
// 			tx.Save(&prompt)
// 			return fmt.Errorf("mismatch: expected %d images, but %d was provided", len(placeholderImages), len(images))
// 		}

// 		// Update placeholders
// 		for i, imageData := range images {
// 			imageToUpdate := &placeholderImages[i]
// 			imageToUpdate.ImageData = imageData
// 			imageToUpdate.UpdatedAt = time.Now()
// 			imageToUpdate.Status = domain.Completed

// 			if err := tx.Save(imageToUpdate).Error; err != nil {
// 				// rollback failure
// 				return err
// 			}
// 		}

// 		// all images updated
// 		prompt.Status = domain.Completed
// 		if err := tx.Save(&prompt).Error; err != nil {
// 			return err
// 		}
// 		prompt.Images = placeholderImages
// 		return nil
// 	})

// 	if err != nil {
// 		r.logger.Error("Failed to update placeholder images", zap.Error(err))
// 		return nil, fmt.Errorf("failed updating placeholder images in database: %w", err)
// 	}
// 	return &prompt, nil
// }

func (r *gormPromptRepository) UpdatePlaceholderImages(ctx context.Context, externalPromptID uuid.UUID, images [][]byte) (*domain.Prompt, error) {
	var prompt domain.Prompt
	var finalErr error // This will hold the business logic error to return.

	txErr := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Find the parent prompt. If it fails, we must abort.
		if err := tx.Where("external_prompt_id = ?", externalPromptID).First(&prompt).Error; err != nil {
			r.logger.Warn("External prompt ID does not exist in database", zap.String("externalPromptID", externalPromptID.String()))
			return err // Abort transaction immediately.
		}

		// 2. Find placeholder images.
		var placeholderImages []domain.Image
		if err := tx.Where("prompt_id = ? AND status = ?", prompt.ID, domain.Pending).Find(&placeholderImages).Error; err != nil {
			finalErr = fmt.Errorf("failed to find placeholder images: %w", err)
			prompt.Status = domain.Failed
		}

		// 3. If we are still in a good state, perform business logic.
		if finalErr == nil {
			if len(placeholderImages) != len(images) {
				finalErr = fmt.Errorf("mismatch: expected %d images, but %d was provided", len(placeholderImages), len(images))
				prompt.Status = domain.Failed
			} else {
				// Try to update the images.
				for i, imageData := range images {
					imageToUpdate := &placeholderImages[i]
					imageToUpdate.ImageData = imageData
					imageToUpdate.UpdatedAt = time.Now()
					imageToUpdate.Status = domain.Completed

					if err := tx.Save(imageToUpdate).Error; err != nil {
						finalErr = fmt.Errorf("failed to save image index %d: %w", i, err)
						prompt.Status = domain.Failed
						break // Stop processing on the first error
					}
				}
			}
		}

		// 4. If no errors occurred during the loop, mark the prompt as completed.
		if finalErr == nil {
			prompt.Status = domain.Completed
			prompt.Images = placeholderImages // Attach successful images for the return value
		}

		// 5. CRITICAL STEP: Now, based on the final prompt status, update the children.
		// This is the simplification you identified.
		if prompt.Status == domain.Failed {
			// An error occurred. Update all images for this prompt to Failed.
			if err := tx.Model(&domain.Image{}).Where("prompt_id = ?", prompt.ID).Update("status", domain.Failed).Error; err != nil {
				// This is a critical failure. We failed to save the "failed" state. We must roll back.
				r.logger.Error("CRITICAL: Failed to update images to FAILED status", zap.Error(err))
				return err // Abort transaction!
			}
		}

		// 6. Save the final state of the parent prompt.
		if err := tx.Save(&prompt).Error; err != nil {
			// If we can't save the final prompt status, we must abort.
			r.logger.Error("CRITICAL: Failed to save final prompt status", zap.Error(err))
			return err // Abort transaction!
		}

		// Return nil to COMMIT the transaction with its final state (either Completed or Failed).
		return nil
	})

	// 7. Handle the results.
	if txErr != nil {
		r.logger.Error("Transaction rolled back due to fatal error", zap.Error(txErr))
		return nil, fmt.Errorf("database transaction failed: %w", txErr)
	}

	if finalErr != nil {
		r.logger.Error("Image update operation failed", zap.Error(finalErr))
		// The business logic failed, but we successfully committed the "Failed" state.
		// Return the business error to the caller.
		return nil, finalErr
	}

	// Success!
	return &prompt, nil
}
