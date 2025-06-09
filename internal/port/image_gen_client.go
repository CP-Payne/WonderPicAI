package port

import "github.com/google/uuid"

type ImageGenerationInput struct {
	Prompt     string
	ImageCount int
	Width      int
	Height     int
}

type ImageGeneration interface {
	GenerateImage(*ImageGenerationInput) (promptID uuid.UUID, err error)
}
