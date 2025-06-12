package stripe

import (
	"github.com/CP-Payne/wonderpicai/internal/port"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"go.uber.org/zap"
)

type StripeProvider struct {
	logger     *zap.Logger
	successURL string
	cancelURL  string
}

func NewProvider(logger *zap.Logger, privateKey string, successURL, cancelURL string) port.PaymentProvider {
	stripe.Key = privateKey
	return &StripeProvider{
		logger:     logger,
		successURL: successURL,
		cancelURL:  cancelURL,
	}
}

func (p *StripeProvider) CreateCheckoutSession(user port.UserData, product port.ProductData) (string, error) {
	params := &stripe.CheckoutSessionParams{
		CustomerEmail: stripe.String(user.Email),
		Mode:          stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(product.Name),
					},
					UnitAmount: stripe.Int64(int64(product.Price) * 100),
				},
				Quantity: stripe.Int64(int64(product.Quantity)),
			},
		},
		SuccessURL: stripe.String(p.successURL),
		CancelURL:  stripe.String(p.cancelURL),
	}
	s, err := session.New(params)
	if err != nil {
		return "", err
	}

	return s.URL, nil
}

func (p *StripeProvider) HandleEvents(msg []byte) {

}
