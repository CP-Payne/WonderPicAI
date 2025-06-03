package viewmodel

type ToastType string

const (
	ToastError   ToastType = "error"
	ToastSuccess ToastType = "success"
	ToastInfo    ToastType = "info"
	ToastWarning ToastType = "warning"
)

type ToastComponentData struct {
	Message string
	Type    ToastType
	ToastID string
}
