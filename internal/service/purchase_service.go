package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"

	"github.com/CP-Payne/wonderpicai/internal/domain"
	"github.com/CP-Payne/wonderpicai/internal/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PurcaseService interface {
	GetOptions(ctx context.Context) []PurchaseOption
	OptionExists(ctx context.Context, option string) bool
	CreateCheckout(ctx context.Context, userID uuid.UUID, option string) (checkoutURL string, err error)
	HandleProviderEvents(r *http.Request, data []byte) error
}

type PurchaseOption struct {
	Credits   int
	Price     int
	ActionURL string
}

// Current purchase options
var options map[string]PurchaseOption = map[string]PurchaseOption{
	"100": {
		Credits:   100,
		Price:     5,
		ActionURL: "/purchase/100",
	},

	"250": {
		Credits:   250,
		Price:     10,
		ActionURL: "/purchase/250",
	},

	"700": {
		Credits:   700,
		Price:     25,
		ActionURL: "/purchase/700",
	},

	"1500": {
		Credits:   1500,
		Price:     50,
		ActionURL: "/purchase/1500",
	},
}

type purchaseService struct {
	logger        *zap.Logger
	walletService WalletService
	provider      port.PaymentProvider
	userRepo      port.UserRepository
}

func NewPurchaseService(logger *zap.Logger, walletService WalletService, provider port.PaymentProvider, userRepo port.UserRepository) PurcaseService {
	return &purchaseService{
		logger:        logger.With(zap.String("component", "PurchaseService")),
		walletService: walletService,
		provider:      provider,
		userRepo:      userRepo,
	}
}

func (s *purchaseService) GetOptions(ctx context.Context) []PurchaseOption {
	var purchaseOptions []PurchaseOption

	for _, option := range options {
		purchaseOptions = append(purchaseOptions, option)
	}

	sort.Slice(purchaseOptions, func(i, j int) bool {
		return purchaseOptions[i].Price < purchaseOptions[j].Price
	})

	return purchaseOptions
}

func (s *purchaseService) OptionExists(ctx context.Context, option string) bool {
	_, exists := options[option]
	return exists
}

func (s *purchaseService) CreateCheckout(ctx context.Context, userID uuid.UUID, option string) (checkoutURL string, err error) {

	opt, ok := options[option]
	if !ok {
		return "", domain.ErrInvalidPurchaseOption
	}

	productName := fmt.Sprintf("%d credits", opt.Credits)

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return "", err
	}
	userData := port.UserData{
		Email: user.Email,
	}

	productData := port.ProductData{
		Name:     productName,
		Price:    opt.Price,
		Quantity: 1,
		Option:   option,
	}
	checkoutURL, err = s.provider.CreateCheckoutSession(userData, productData)
	if err != nil {
		return "", fmt.Errorf("failed to create checkout session: %w", err)
	}

	return checkoutURL, nil
}

func (s *purchaseService) HandleProviderEvents(r *http.Request, data []byte) error {

	sessionData, err := s.provider.HandleEvent(r, data)
	if err != nil {
		if errors.Is(err, domain.ErrUnhandledEvent) {
			s.logger.Warn("skipping unhandled event from provided")
		}
		return err
	}

	purchasedOption, ok := options[sessionData.Option]
	if !ok {
		return domain.ErrInvalidPurchaseOption
	}

	err = s.walletService.AddCredits(r.Context(), sessionData.UserEmail, purchasedOption.Credits)
	if err != nil {
		s.logger.Error("CRITICAL - Failed adding credits to user account", zap.String("email", sessionData.UserEmail), zap.Int("amount", purchasedOption.Credits), zap.Error(err))
		return err
	}

	return nil

}
