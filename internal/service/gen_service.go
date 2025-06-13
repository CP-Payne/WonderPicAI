package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CP-Payne/wonderpicai/internal/domain"
	"github.com/CP-Payne/wonderpicai/internal/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// TODO: Pass in context from handler
type GenService interface {
	GenerateImage(ctx context.Context, userID uuid.UUID, promptData *PromptData) (*domain.Prompt, error)
	CalculateCost(ctx context.Context, promptData *PromptData) int
	GetAllPrompts(ctx context.Context, userID uuid.UUID) ([]domain.Prompt, error)
	UpdatePlaceholderImages(externalPromptID uuid.UUID, images [][]byte) (*domain.Prompt, error)
	GetImageByID(ctx context.Context, userID uuid.UUID, imageID uuid.UUID) (image *domain.Image, err error)
	DeleteImageByID(ctx context.Context, userID uuid.UUID, imageID uuid.UUID) error
	DeleteFailedImages(ctx context.Context, userID uuid.UUID) error
}

type PromptData struct {
	Prompt     string
	ImageCount int
}

const (
	IMAGE_WIDTH     = 500
	IMAGE_HEIGHT    = 500
	GENERATION_COST = 2
)

type genService struct {
	logger         *zap.Logger
	imageGenClient port.ImageGeneration
	walletService  WalletService
	promptRepo     port.PromptRepository
	imageRepo      port.ImageRepository
}

func NewGenService(logger *zap.Logger, genClient port.ImageGeneration, promptRepo port.PromptRepository, imageRepo port.ImageRepository, walletService WalletService) GenService {
	return &genService{
		logger:         logger.With(zap.String("component", "GenService")),
		imageGenClient: genClient,
		promptRepo:     promptRepo,
		imageRepo:      imageRepo,
		walletService:  walletService,
	}
}

func (s *genService) GenerateImage(ctx context.Context, userID uuid.UUID, data *PromptData) (*domain.Prompt, error) {

	clientReqData := port.ImageGenerationInput{
		Prompt:     data.Prompt,
		ImageCount: data.ImageCount,
		Width:      IMAGE_WIDTH,
		Height:     IMAGE_HEIGHT,
	}

	totalCost := s.CalculateCost(ctx, data)

	err := s.walletService.DeductForImageGeneration(ctx, userID, totalCost)
	if err != nil {
		if errors.Is(err, domain.ErrInsufficientFunds) {
			s.logger.Error("Failed to deduct credits for image generation", zap.String("userID", userID.String()), zap.Int("totalCost", totalCost), zap.Error(err))
			return nil, err
		}

		s.logger.Error("Failed to deduct credits for image generation", zap.String("userID", userID.String()), zap.Int("totalCost", totalCost), zap.Error(err))
		return nil, fmt.Errorf("failed to deduct credits for image generation: %w", err)
	}

	externalPromptID, err := s.imageGenClient.GenerateImage(&clientReqData)
	if err != nil {
		s.logger.Error("Failed to send image generation request", zap.Error(err))
		s.logger.Info("Refunding credits", zap.Error(err))
		err = s.walletService.RefundCredits(ctx, userID, totalCost)
		if err != nil {
			s.logger.Error("CRITICAL: Failed to refund credits", zap.String("userID", userID.String()), zap.Error(err), zap.Int("totalCost", totalCost))
			return nil, fmt.Errorf("failed refunding credits after failed generation request: %w", err)
		}

		return nil, fmt.Errorf("failed to initiate image generation - credits refunded: %w", err)
	}

	s.logger.Debug("External Prompt Received", zap.String("ExternalPromptID", externalPromptID.String()))

	prompt := domain.Prompt{
		BaseModel: domain.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserID:           userID,
		ExternalPromptID: externalPromptID,
		Cost:             totalCost,
		ImageCount:       data.ImageCount,
		Width:            IMAGE_WIDTH,
		Height:           IMAGE_HEIGHT,
		Status:           domain.Pending,
		LastChecked:      time.Now(),
	}

	promptCreated, err := s.promptRepo.Create(ctx, &prompt)
	if err != nil {
		s.logger.Error("Prompt creation failed via repository", zap.Error(err))
		return nil, fmt.Errorf("failed to complete prompt image generation due to an internal issue: %w", err)
	}

	return promptCreated, nil

}

func (s *genService) CalculateCost(ctx context.Context, data *PromptData) int {
	if data == nil {
		return 0
	}

	return data.ImageCount * GENERATION_COST
}

func (s *genService) GetAllPrompts(ctx context.Context, userID uuid.UUID) ([]domain.Prompt, error) {
	prompts, err := s.promptRepo.FindAllByUser(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to retrieve user prompts from repository", zap.String("userID", userID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to retrieve user prompts from prompt repository: %w", err)
	}

	return prompts, nil
}

func (s *genService) UpdatePlaceholderImages(externalPromptID uuid.UUID, images [][]byte) (*domain.Prompt, error) {

	prompt, err := s.promptRepo.UpdatePlaceholderImages(context.Background(), externalPromptID, images)
	if err != nil {
		s.logger.Error("Failed to update image placeholders", zap.String("ExternalPromptID", externalPromptID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to update image placeholders using prompt repository: %w", err)
	}

	return prompt, nil
}

func (s *genService) GetImageByID(ctx context.Context, userID uuid.UUID, imageID uuid.UUID) (image *domain.Image, err error) {
	image, err = s.imageRepo.GetByID(ctx, userID, imageID)
	if err != nil {
		if errors.Is(err, domain.ErrImageNotFound) {
			return nil, err
		}
		s.logger.Error("Repository failed to get image by ID",
			zap.String("imageID", imageID.String()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to retrieve image: %w", err)
	}

	return image, nil
}

func (s *genService) DeleteImageByID(ctx context.Context, userID uuid.UUID, imageID uuid.UUID) error {
	s.logger.Info("Attempting to delete image",
		zap.String("userID", userID.String()),
		zap.String("imageID", imageID.String()),
	)

	err := s.imageRepo.Delete(ctx, userID, imageID)
	if err != nil {
		if errors.Is(err, domain.ErrRecordNotFound) {
			s.logger.Warn("Delete request for a non-existent or already deleted image",
				zap.String("userID", userID.String()),
				zap.String("imageID", imageID.String()),
			)
			return nil
		}

		s.logger.Error("Repository failed to delete image",
			zap.String("userID", userID.String()),
			zap.String("imageID", imageID.String()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete image: %w", err)
	}

	s.logger.Info("Successfully deleted image",
		zap.String("userID", userID.String()),
		zap.String("imageID", imageID.String()),
	)
	return nil
}

func (s *genService) DeleteFailedImages(ctx context.Context, userID uuid.UUID) error {
	s.logger.Info("Attempting to delete all failed images for user",
		zap.String("userID", userID.String()),
	)

	if err := s.imageRepo.DeleteFailed(ctx, userID); err != nil {
		s.logger.Error("Repository failed to delete failed images",
			zap.String("userID", userID.String()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete failed images: %w", err)
	}

	s.logger.Info("Successfully processed request to delete failed images for user",
		zap.String("userID", userID.String()),
	)
	return nil
}
