package port

import "net/http"

type UserData struct {
	Email string
}

type ProductData struct {
	Name     string
	Price    int
	Quantity int
	Option   string
}

type SessionSuccess struct {
	UserEmail string
	Option    string
}

type PaymentProvider interface {
	CreateCheckoutSession(UserData, ProductData) (string, error)
	HandleEvent(r *http.Request, data []byte) (*SessionSuccess, error)
}
