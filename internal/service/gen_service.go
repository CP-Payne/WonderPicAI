package service

import (
	"context"
	"fmt"
	"time"

	"github.com/CP-Payne/wonderpicai/internal/domain"
	"github.com/CP-Payne/wonderpicai/internal/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type GenService interface {
	GenerateImage(*PromptData) (*domain.Prompt, error)
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
	// prompt and image repo here
	// comfyLite client here
}

func NewGenService(logger *zap.Logger, genClient port.ImageGeneration, promptRepo port.PromptRepository) GenService {
	return &genService{
		logger:         logger.With(zap.String("component", "GenService")),
		imageGenClient: genClient,
		promptRepo:     promptRepo,
	}
}

func (s *genService) GenerateImage(data *PromptData) (*domain.Prompt, error) {

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
