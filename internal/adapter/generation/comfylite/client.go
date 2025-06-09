package comfylite

import (
	"github.com/CP-Payne/wonderpicai/internal/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ComfyLiteClient struct {
	logger *zap.Logger
}

func NewClient(logger *zap.Logger) port.ImageGeneration {
	return &ComfyLiteClient{
		logger: logger,
	}
}

func (c *ComfyLiteClient) GenerateImage(input *port.ImageGenerationInput) (uuid.UUID, error) {
	// temporarily return random UUID for testing
	return uuid.New(), nil
}
