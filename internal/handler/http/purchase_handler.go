package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/CP-Payne/wonderpicai/internal/context/auth"
	"github.com/CP-Payne/wonderpicai/internal/handler/http/response"
	"github.com/CP-Payne/wonderpicai/internal/service"
	creditPages "github.com/CP-Payne/wonderpicai/web/template/pages/credits"
	"github.com/CP-Payne/wonderpicai/web/template/viewmodel"
)

type PurchaseHandler struct {
	logger          *zap.Logger
	validate        *validator.Validate
	purchaseService service.PurcaseService
}

func NewPurchaseHandler(logger *zap.Logger, validate *validator.Validate, purchaseService service.PurcaseService) *PurchaseHandler {
	return &PurchaseHandler{
		logger:          logger.With(zap.String("component", "PurchaseHandler")),
		validate:        validate,
		purchaseService: purchaseService,
	}
}

func (h *PurchaseHandler) ShowPurchasePage(w http.ResponseWriter, r *http.Request) {

	// Building viewmodel options for display

	availableOptions := h.purchaseService.GetOptions(r.Context())

	viewOptions := make([]viewmodel.PurchaseOption, len(availableOptions))

	for i, availOpt := range availableOptions {
		viewOptions[i] = viewmodel.PurchaseOption{
			Credits:   availOpt.Credits,
			Price:     availOpt.Price,
			ActionURL: availOpt.ActionURL,
		}
	}

	purchasePageData := viewmodel.PurchaseViewData{
		Options: viewOptions,
	}

	err := creditPages.PurchasePage(purchasePageData).Render(r.Context(), w)
	if err != nil {
		h.logger.Error("Failed to render Purchase page", zap.Error(err))
		response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
		return
	}
}

func (h *PurchaseHandler) ShowSuccessPage(w http.ResponseWriter, r *http.Request) {

	err := creditPages.SuccessPage().Render(r.Context(), w)
	if err != nil {
		h.logger.Error("Failed to render success page", zap.Error(err))
		response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
		return
	}
}

func (h *PurchaseHandler) ShowCancelPage(w http.ResponseWriter, r *http.Request) {

	err := creditPages.CancelPage().Render(r.Context(), w)
	if err != nil {
		h.logger.Error("Failed to render success page", zap.Error(err))
		response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
		return
	}
}

func (h *PurchaseHandler) HandlePurchaseOption(w http.ResponseWriter, r *http.Request) {
	option := chi.URLParam(r, "option")

	exist := h.purchaseService.OptionExists(r.Context(), option)
	if !exist {
		toastID, loadErr := response.LoadErrorToast(w, r, h.logger, "invalid option")
		if loadErr != nil {
			h.logger.Error("failed loading ErrorToast", zap.String("toastID", toastID), zap.Error(loadErr))
			response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
			return
		}
		// Rerender page with purchase options including the added toast
		availableOptions := h.purchaseService.GetOptions(r.Context())

		viewOptions := make([]viewmodel.PurchaseOption, len(availableOptions))

		for i, availOpt := range availableOptions {
			viewOptions[i] = viewmodel.PurchaseOption{
				Credits:   availOpt.Credits,
				Price:     availOpt.Price,
				ActionURL: availOpt.ActionURL,
			}
		}

		purchasePageData := viewmodel.PurchaseViewData{
			Options: viewOptions,
		}

		response.LoadPurchasePage(w, r, h.logger, purchasePageData)
		// End request processing
		return

	}

	userID, err := auth.UserID(r.Context())
	if err != nil {
		h.logger.Error("Failed to get userID from context", zap.String("userID", userID.String()), zap.Error(err))
		response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
		return
	}

	checkoutURL, err := h.purchaseService.CreateCheckout(r.Context(), userID, option)
	if err != nil {
		h.logger.Error("Failed to get checkout URL", zap.String("option", option+" credits"), zap.String("userID", userID.String()), zap.Error(err))
		response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
		return
	}

	response.HxRedirect(w, r, checkoutURL)
}
