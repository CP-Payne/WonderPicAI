package comfylite

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/CP-Payne/wonderpicai/internal/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ComfyLiteClient struct {
	logger     *zap.Logger
	baseURL    string
	HttpClient *http.Client
}

func NewClient(logger *zap.Logger, baseUrl string) port.ImageGeneration {
	return &ComfyLiteClient{
		logger:     logger,
		HttpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    baseUrl,
	}
}

type ComfyGenRequest struct {
	Prompt     string `json:"prompt"`
	ImageCount int    `json:"batch_size"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
}

type ComfyGenResponse struct {
	PromptID string `json:"prompt_id"`
}

func (c *ComfyLiteClient) GenerateImage(input *port.ImageGenerationInput) (uuid.UUID, error) {
	// temporarily return random UUID for testing
	url := c.baseURL + "/gen"

	data, err := json.Marshal(ComfyGenRequest{
		Prompt:     input.Prompt,
		ImageCount: input.ImageCount,
		Width:      input.Width,
		Height:     input.Height,
	})
	if err != nil {
		c.logger.Error("Failed to marshal ComfyGenRequest", zap.Error(err))
		return uuid.UUID{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		c.logger.Error("failed to create request", zap.Error(err))
		return uuid.UUID{}, fmt.Errorf("failed to create request to ComfyLite: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		c.logger.Error("Error sending request", zap.String("destination", url+"/gen"), zap.Error(err))
		return uuid.UUID{}, fmt.Errorf("failed sending request to comfylite: %w", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Error reading response body", zap.Error(err))
		return uuid.UUID{}, fmt.Errorf("error reading response body from comfylite: %w", err)
	}

	comfyResp := ComfyGenResponse{}

	err = json.Unmarshal(respBody, &comfyResp)
	if err != nil {
		c.logger.Error("Failed to read ComfyLite response", zap.Error(err))
		return uuid.UUID{}, fmt.Errorf("failed to read comfyLite response: %w", err)
	}

	promptUUID, err := uuid.Parse(comfyResp.PromptID)
	if err != nil {
		c.logger.Error("failed to convert promptID string to uuid", zap.Error(err))
		return uuid.UUID{}, fmt.Errorf("invalid prompt id received: %w", err)
	}

	return promptUUID, nil
}
