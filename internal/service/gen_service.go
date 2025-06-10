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
	GenerateImage(userID uuid.UUID, promptData *PromptData) (*domain.Prompt, error)
	GetAllPrompts(userID uuid.UUID) ([]domain.Prompt, error)
	UpdatePlaceholderImages(externalPromptID uuid.UUID, images [][]byte) (*domain.Prompt, error)
	GetImageByID(imageID uuid.UUID) (image *domain.Image, err error)
	DeleteImageByID(imageID uuid.UUID) error
	DeleteFailedImages() error
}

type PromptData struct {
	Prompt     string
	ImageCount int
}

const (
	IMAGE_WIDTH  = 500
	IMAGE_HEIGHT = 500
)

type genService struct {
	logger         *zap.Logger
	imageGenClient port.ImageGeneration
	promptRepo     port.PromptRepository
	imageRepo      port.ImageRepository
}

func NewGenService(logger *zap.Logger, genClient port.ImageGeneration, promptRepo port.PromptRepository, imageRepo port.ImageRepository) GenService {
	return &genService{
		logger:         logger.With(zap.String("component", "GenService")),
		imageGenClient: genClient,
		promptRepo:     promptRepo,
		imageRepo:      imageRepo,
	}
}

func (s *genService) GenerateImage(userID uuid.UUID, data *PromptData) (*domain.Prompt, error) {

	clientReqData := port.ImageGenerationInput{
		Prompt:     data.Prompt,
		ImageCount: data.ImageCount,
		Width:      IMAGE_WIDTH,
		Height:     IMAGE_HEIGHT,
	}

	externalPromptID, err := s.imageGenClient.GenerateImage(&clientReqData)
	if err != nil {
		s.logger.Error("Failed to send image generation request to image generation server", zap.Error(err))
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
		Cost:             1,
		ImageCount:       data.ImageCount,
		Width:            IMAGE_WIDTH,
		Height:           IMAGE_HEIGHT,
		Status:           domain.Pending,
		LastChecked:      time.Now(),
	}

	promptCreated, err := s.promptRepo.Create(context.Background(), &prompt)
	if err != nil {
		s.logger.Error("Prompt creation failed via repository", zap.Error(err))
		return nil, fmt.Errorf("failed to complete prompt image generation due to an internal issue: %w", err)
	}

	return promptCreated, nil

}

func (s *genService) GetAllPrompts(userID uuid.UUID) ([]domain.Prompt, error) {
	prompts, err := s.promptRepo.FindAllByUser(context.Background(), userID)
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

func (s *genService) GetImageByID(imageID uuid.UUID) (image *domain.Image, err error) {
	image, err = s.imageRepo.GetByID(imageID)
	if err != nil {
		if errors.Is(err, domain.ErrImageNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to retrieve image from repository: %w", err)
	}

	return image, nil
}

func (s *genService) DeleteImageByID(imageID uuid.UUID) error {
	err := s.imageRepo.Delete(imageID)
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}
	return nil
}

func (s *genService) DeleteFailedImages() error {
	err := s.imageRepo.DeleteFailed()

	if err != nil {
		return fmt.Errorf("failed to delete images with status failed: %w", err)
	}
	return nil
}
