package stripe

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/CP-Payne/wonderpicai/internal/domain"
	"github.com/CP-Payne/wonderpicai/internal/port"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/product"
	"github.com/stripe/stripe-go/v82/webhook"
	"go.uber.org/zap"
)

type StripeProvider struct {
	logger         *zap.Logger
	endpointSecret string
	successURL     string
	cancelURL      string
}

func NewProvider(logger *zap.Logger, privateKey, endpointSecret string, successURL, cancelURL string) port.PaymentProvider {
	stripe.Key = privateKey
	return &StripeProvider{
		logger:         logger,
		successURL:     successURL,
		cancelURL:      cancelURL,
		endpointSecret: endpointSecret,
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
						Metadata: map[string]string{
							"option": product.Option,
						},
					},
					UnitAmount: stripe.Int64(int64(product.Price) * 100),
				},
				Quantity: stripe.Int64(int64(product.Quantity)),
			},
		},
		Expand: []*string{
			stripe.String("line_items.data.price.product"),
		},
		SuccessURL: stripe.String(p.successURL),
		CancelURL:  stripe.String(p.cancelURL),
	}
	params.AddExpand("line_items")
	s, err := session.New(params)
	if err != nil {
		return "", err
	}

	return s.URL, nil
}

func (p *StripeProvider) HandleEvent(r *http.Request, data []byte) (*port.SessionSuccess, error) {

	event := stripe.Event{}

	if err := json.Unmarshal(data, &event); err != nil {
		p.logger.Error("Failed parsing basic request into event", zap.Error(err))
		return nil, fmt.Errorf("failed parsing data into stripe event: %w", err)
	}

	endpointSecret := p.endpointSecret
	signatureHeader := r.Header.Get("Stripe-Signature")
	event, err := webhook.ConstructEvent(data, signatureHeader, endpointSecret)
	if err != nil {
		p.logger.Error("Stripe webhook signature verification failed", zap.Error(err))
		return nil, fmt.Errorf("stripe webhook signature verification failed: %w", err)
	}
	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	case "checkout.session.completed":
		var sessionEvent stripe.CheckoutSession

		err := json.Unmarshal(event.Data.Raw, &sessionEvent)
		if err != nil {
			p.logger.Error("Failed parsing stripe webhook JSON", zap.Error(err))
			return nil, fmt.Errorf("failed parsing stripe webhook json into checkout session: %w", err)
		}

		// Get Line items
		params := &stripe.CheckoutSessionListLineItemsParams{
			Session: stripe.String(sessionEvent.ID),
		}

		var option string

		result := session.ListLineItems(params)

		for result.Next() {
			item := result.LineItem()

			if item.Price == nil || item.Price.Product == nil {
				p.logger.Error("Price or Product field is nil in line item")
				return nil, fmt.Errorf("failed accessing price or product field in line item")
			}

			productID := item.Price.Product.ID

			stripeProduct, err := product.Get(productID, nil)
			if err != nil {
				p.logger.Error("Failed to retrieve product from Stripe", zap.Error(err))
				return nil, fmt.Errorf("failed to retrieve product from stripe: %w", err)
			}

			option = stripeProduct.Metadata["option"]
		}

		sessionSuccess := port.SessionSuccess{
			UserEmail: sessionEvent.CustomerDetails.Email,
			Option:    option,
		}

		// p.logger.Debug("Session data after successfull purchase", zap.Any("session", sessionEvent))
		return &sessionSuccess, nil

	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
		p.logger.Warn("Unhandled event type", zap.Any("type", event.Type))
	}

	return nil, domain.ErrUnhandledEvent

}
