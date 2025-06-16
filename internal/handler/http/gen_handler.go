package http

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/CP-Payne/wonderpicai/internal/context/auth"
	"github.com/CP-Payne/wonderpicai/internal/context/credits"
	"github.com/CP-Payne/wonderpicai/internal/domain"
	"github.com/CP-Payne/wonderpicai/internal/handler/http/response"
	"github.com/CP-Payne/wonderpicai/internal/service"
	"github.com/CP-Payne/wonderpicai/internal/validation"
	genPages "github.com/CP-Payne/wonderpicai/web/template/pages/gen"
	"github.com/CP-Payne/wonderpicai/web/template/viewmodel"
)

type GenHandler struct {
	logger     *zap.Logger
	validate   *validator.Validate
	genService service.GenService
}

type GenRequest struct {
	Prompt     string `validate:"required,min=3"`
	ImageCount int    `validate:"required,number,gte=1,lte=10"`
}

type ImageUpdateWebhookRequest struct {
	Status   string   `json:"status"`
	PromptID string   `json:"prompt_id"`
	Images   []string `json:"images"`
	Error    string   `json:"error"`
}

func NewGenHandler(logger *zap.Logger, validate *validator.Validate, genService service.GenService) *GenHandler {
	return &GenHandler{
		logger:     logger.With(zap.String("component", "GenHandler")),
		validate:   validate,
		genService: genService,
	}
}

func (h *GenHandler) ShowGenPage(w http.ResponseWriter, r *http.Request) {

	var images []viewmodel.Image

	userID, err := auth.UserID(r.Context())
	if err != nil {
		h.logger.Error("Failed to get UserID from context")
		response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
		return
	}

	userCredits, err := credits.RemainingCredits(r.Context())
	if err != nil {
		h.logger.Error("failed to retrieve user credits from context", zap.Error(err))
		response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
		return
	}

	userPrompts, err := h.genService.GetAllPrompts(r.Context(), userID)
	if err != nil {
		h.logger.Error("failed to retrieve images from genService", zap.Error(err), zap.String("userID", userID.String()))
		toastID, loadErr := response.LoadErrorToast(w, r, h.logger, "failed loading images")
		if loadErr != nil {
			response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, toastID, "")
			return
		}
		images = []viewmodel.Image{}
	}

	containsFailedImages := false

	// Extract images from prompts
	for _, prompt := range userPrompts {
		for _, img := range prompt.Images {
			if img.Status == domain.Failed {
				containsFailedImages = true
			}
			images = append(images, viewmodel.Image{
				ID:     img.ID.String(),
				Data:   base64.StdEncoding.EncodeToString(img.ImageData),
				Status: img.Status.String(),
			})
		}
	}

	// Calculate minimum cost
	minCost := h.genService.CalculateCost(r.Context(), &service.PromptData{
		ImageCount: 1,
	})

	genPageData := viewmodel.GenPageData{
		GalleryData: viewmodel.GalleryComponentData{
			Images: images,
		},
		GenFormData: viewmodel.GenFormComponentData{
			Form: viewmodel.GenFormData{
				MinCost:         minCost,
				Credits:         userCredits,
				MaxImagesPerGen: 10,
				HasFailedImages: containsFailedImages,
				ImageCount:      1,
			},
			Errors: map[string]string{},
			Error:  "",
		},
	}
	err = genPages.GenPage(genPageData).Render(r.Context(), w)
	if err != nil {
		h.logger.Error("Failed to render login page", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *GenHandler) HandleGenerationCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.logger.Error("Failed to parse form", zap.Error(err))
		response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, xid.New().String(), "")
		return
	}

	userID, err := auth.UserID(r.Context())
	if err != nil {
		h.logger.Error("Failed to get UserID from context")
		response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
		return
	}

	// Get user credits
	userCredits, err := credits.RemainingCredits(r.Context())
	if err != nil {
		h.logger.Error("failed to retrieve user credits from context", zap.Error(err))
		response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
		return
	}

	containsFailedImages, err := h.genService.ContainsFailedImages(r.Context(), userID)
	if err != nil {
		h.logger.Error("failed to determine if user has failed images defaulting to false", zap.Error(err))
		containsFailedImages = false
	}

	vm := viewmodel.GenFormComponentData{
		Form: viewmodel.GenFormData{
			Prompt:          r.FormValue("prompt"),
			Credits:         userCredits,
			MinCost:         h.genService.CalculateCost(r.Context(), &service.PromptData{ImageCount: 1}),
			MaxImagesPerGen: 10,
			HasFailedImages: containsFailedImages,
		},
		Errors: map[string]string{},
		Error:  "",
	}

	imageCountStr := r.FormValue("image_count")
	imageCount, err := strconv.Atoi(imageCountStr)
	if err != nil {
		h.logger.Error("Failed to parse image count", zap.Error(err))
		vm.Errors["imageCount"] = "invalid image count"
		loadErr := response.LoadGenForm(w, r, h.logger, vm)
		if loadErr != nil {
			response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
			return
		}
		// End request processing here
		return
	}

	req := GenRequest{
		Prompt:     r.FormValue("prompt"),
		ImageCount: imageCount,
	}
	vm.Form.ImageCount = req.ImageCount

	err = h.validate.Struct(req)
	if err != nil {
		fieldErrors, generalValError := validation.TranslateValidationErrors(err)
		vm.Errors = fieldErrors
		vm.Error = generalValError

		if vm.Error != "" {
			h.logger.Error("General generation validation error", zap.String("error", vm.Error), zap.Error(err))

			toastID, loadErr := response.LoadErrorToast(w, r, h.logger, vm.Error)
			if loadErr != nil {
				response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, toastID, "")
				return
			}
		}

		h.logger.Warn("generation validation errors", zap.Any("errors", vm.Errors))

		loadErr := response.LoadGenForm(w, r, h.logger, vm)
		if loadErr != nil {
			response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
			return
		}
		return
	}

	// -- Validation Passed ---

	// Clear validation errors
	vm.Errors = make(map[string]string)
	vm.Error = ""

	// Check funds
	// --- Note ---
	// Avaible funds is also checked in the service and repository layer
	// This check is to prevent unecessary calls to service and repository layer
	totalCost := h.genService.CalculateCost(r.Context(), &service.PromptData{ImageCount: req.ImageCount})

	if totalCost > userCredits {
		vm.Errors["credits"] = fmt.Sprintf("insufficient credits - you require %d more", totalCost-userCredits)
		loadErr := response.LoadGenForm(w, r, h.logger, vm)
		if loadErr != nil {
			response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
			return
		}
		return
	}

	// --- Check

	prompt, err := h.genService.GenerateImage(r.Context(), userID, &service.PromptData{
		Prompt:     req.Prompt,
		ImageCount: req.ImageCount,
	})

	if err != nil {
		h.logger.Error("failed to create prompt", zap.Error(err))
		response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
		return
	}

	// Generation request successfully sent show new credit balance
	vm.Form.Credits = vm.Form.Credits - totalCost
	// TODO: if the image generation API fails generating images, refund credits (events)

	// Load new pending images
	for _, image := range prompt.Images {
		loadErr := response.LoadOOBPendingImage(w, r, h.logger, viewmodel.Image{
			ID:     image.ID.String(),
			Data:   string(image.ImageData),
			Status: "Pending",
		})
		if loadErr != nil {

			toastID, loadErr := response.LoadErrorToast(w, r, h.logger, "failed to load pending image")
			if loadErr != nil {
				response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, toastID, "")
				return
			}
		}
	}

	// Images loaded

	loadErr := response.LoadGenForm(w, r, h.logger, vm)
	if loadErr != nil {
		response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
		return
	}

}

func (h *GenHandler) HandleImageCompletionWebhook(w http.ResponseWriter, r *http.Request) {

	request := ImageUpdateWebhookRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.logger.Error("failed to decode webhook request body", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	externalPromptID, err := uuid.Parse(request.PromptID)
	if err != nil {
		h.logger.Error("webhook received with invalid promptID format", zap.String("promptID", request.PromptID), zap.Error(err))
		http.Error(w, "invalid prompt_id format", http.StatusBadRequest)
		return
	}

	if request.Status == "failure" {
		h.logger.Warn("received failure webhook for prompt",
			zap.String("promptID", request.PromptID),
			zap.String("error", request.Error),
		)

		_, err := h.genService.UpdatePlaceholderImages(r.Context(), externalPromptID, [][]byte{}, domain.Failed)
		if err != nil {
			h.logger.Error("failed to update prompt status to failed",
				zap.String("promptID", request.PromptID),
				zap.Error(err),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	if len(request.Images) == 0 {
		h.logger.Error("success webhook received with no images", zap.String("promptID", request.PromptID))
		http.Error(w, "no images provided in success payload", http.StatusBadRequest)
		return
	}

	imagesDecoded := make([][]byte, 0, len(request.Images))

	for i, imageData := range request.Images {
		decoded, err := base64.StdEncoding.DecodeString(imageData)
		if err != nil {
			h.logger.Error("image base64 could not be decoded",
				zap.String("promptID", request.PromptID),
				zap.Int("imageIndex", i),
				zap.Error(err),
			)
			http.Error(w, fmt.Sprintf("invalid base64 data for image at index %d", i), http.StatusBadRequest)
			return
		}
		imagesDecoded = append(imagesDecoded, decoded)
	}

	_, err = h.genService.UpdatePlaceholderImages(r.Context(), externalPromptID, imagesDecoded, domain.Completed)
	if err != nil {
		h.logger.Error("failed to update placeholder images to status completed",
			zap.String("promptID", request.PromptID),
			zap.Error(err),
		)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Info("successfully processed completion webhook", zap.String("promptID", request.PromptID))
	w.WriteHeader(http.StatusOK)

}

func (h *GenHandler) HandleImageStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Warn("invalid image uuid provided", zap.Error(err), zap.String("id", idStr))
		return
	}

	userID, err := auth.UserID(r.Context())
	if err != nil {
		h.logger.Error("Failed to get UserID from context")
		response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
		return
	}

	image, err := h.genService.GetImageByID(r.Context(), userID, id)
	if err != nil {
		h.logger.Error("failed to retrieve image by ID", zap.Error(err))
		return
	}

	if image.Status == domain.Completed {
		vm := viewmodel.Image{
			ID:     image.ID.String(),
			Data:   base64.StdEncoding.EncodeToString(image.ImageData),
			Status: "completed",
		}

		loadErr := response.LoadCompletedImage(w, r, h.logger, vm)
		if loadErr != nil {
			response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
			return
		}
		return
	}

	if image.Status == domain.Failed {

		vm := viewmodel.Image{
			ID:     image.ID.String(),
			Data:   base64.StdEncoding.EncodeToString(image.ImageData),
			Status: "failed",
		}

		loadErr := response.LoadFailedImage(w, r, h.logger, vm)
		if loadErr != nil {
			response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
			return
		}
		return
	}

	// Don't do anything if still pending:
	w.WriteHeader(http.StatusNoContent)

}

func (h *GenHandler) HandleImageDelete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Warn("invalid image uuid provided", zap.Error(err), zap.String("id", idStr))
		return
	}

	userID, err := auth.UserID(r.Context())
	if err != nil {
		h.logger.Error("Failed to get UserID from context")
		response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
		return
	}

	err = h.genService.DeleteImageByID(r.Context(), userID, id)
	if err != nil {
		h.logger.Error("failed to delete image", zap.Error(err))

		toastID, loadErr := response.LoadErrorToast(w, r, h.logger, "deletion failed")
		if loadErr != nil {
			h.logger.Error("failed loading ErrorToast", zap.String("toastID", toastID), zap.Error(err))
			response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
			return
		}
		return
	}

	toastID, loadErr := response.LoadSuccessToast(w, r, h.logger, "image deleted")
	if loadErr != nil {
		h.logger.Error("failed loading SuccessToast", zap.String("toastID", toastID), zap.Error(err))
		response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
		return
	}

}

func (h *GenHandler) HandleFailedImagesDelete(w http.ResponseWriter, r *http.Request) {

	userID, err := auth.UserID(r.Context())
	if err != nil {
		h.logger.Error("Failed to get UserID from context")
		response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
		return
	}

	err = h.genService.DeleteFailedImages(r.Context(), userID)
	if err != nil {

		toastID, loadErr := response.LoadErrorToast(w, r, h.logger, "deletion failed")
		if loadErr != nil {
			h.logger.Error("failed loading ErrorToast", zap.String("toastID", toastID), zap.Error(err))
			response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
			return
		}
		return
	}

	response.HxRedirect(w, r, "/gen")

}
