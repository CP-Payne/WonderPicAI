package http

import (
	"net/http"
	"strconv"

	errorpages "github.com/CP-Payne/wonderpicai/web/template/pages/error"
	"github.com/CP-Payne/wonderpicai/web/template/viewmodel"
	"go.uber.org/zap"
)

type ErrorHandler struct {
	logger *zap.Logger
}

func NewErrorHandler(logger *zap.Logger) *ErrorHandler {
	return &ErrorHandler{logger: logger}
}

func (h *ErrorHandler) ServeGenericErrorPage(w http.ResponseWriter, r *http.Request) {
	statusCodeStr := r.URL.Query().Get("statusCode")
	errorID := r.URL.Query().Get("errorID")
	customMessage := r.URL.Query().Get("message")

	statusCode := http.StatusInternalServerError
	if sc, err := strconv.Atoi(statusCodeStr); err == nil && sc >= 400 && sc < 600 {
		statusCode = sc
	}

	pageData := viewmodel.ErrorPageData{
		StatusCode:  statusCode,
		Title:       "Oops! Something Went Wrong",
		ErrorID:     errorID,
		Message:     customMessage,
		ShowDetails: false,
	}

	if customMessage == "" {
		switch statusCode {
		case http.StatusNotFound:
			pageData.Title = "Page Not Found"
			pageData.Message = "Sorry, we couldn't find the page you were looking for."
		case http.StatusForbidden:
			pageData.Title = "Access Denied"
			pageData.Message = "Sorry, you don't have permission to access this page."
		default:
			pageData.Message = "We're sorry, but an unexpected error occured. Our team has been notified."
		}
	}

	w.WriteHeader(statusCode)
	component := errorpages.ServerErrorPage(pageData)
	err := component.Render(r.Context(), w)
	if err != nil {
		h.logger.Error("Failed to render error page template itself!", zap.Error(err))
		http.Error(w, "An error occured, and then another error occured while trying to display the first error.", http.StatusInternalServerError)
	}

}
