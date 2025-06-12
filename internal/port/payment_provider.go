package port

type UserData struct {
	Email string
}

type ProductData struct {
	Name     string
	Price    int
	Quantity int
}

type PaymentProvider interface {
	CreateCheckoutSession(UserData, ProductData) (string, error)
	HandleEvents(msg []byte)
}
