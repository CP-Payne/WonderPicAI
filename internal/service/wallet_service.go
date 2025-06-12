package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/CP-Payne/wonderpicai/internal/domain"
	"github.com/CP-Payne/wonderpicai/internal/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type WalletService interface {
	DeductForImageGeneration(ctx context.Context, userID uuid.UUID, amount int) error
	GetWallet(ctx context.Context, userID uuid.UUID) (*domain.Wallet, error)
	RefundCredits(ctx context.Context, userID uuid.UUID, amount int) error
	AddCredits(ctx context.Context, email string, amount int) error
}

type walletService struct {
	logger     *zap.Logger
	walletRepo port.WalletRepository
}

func NewWalletService(logger *zap.Logger, walletRepo port.WalletRepository) WalletService {
	return &walletService{
		logger:     logger.With(zap.String("component", "WalletService")),
		walletRepo: walletRepo,
	}
}

func (s *walletService) DeductForImageGeneration(ctx context.Context, userID uuid.UUID, amount int) error {

	err := s.walletRepo.SubtractCredits(ctx, userID, amount)
	if err != nil {
		if errors.Is(err, domain.ErrInsufficientFunds) {
			return err
		}
		return fmt.Errorf("failed to deduct credits using wallet repository: %w", err)
	}

	return nil
}

func (s *walletService) GetWallet(ctx context.Context, userID uuid.UUID) (wallet *domain.Wallet, err error) {

	wallet, err = s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Repository failed to get wallet by user ID", zap.String("userID", userID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to retrieve wallet: %w", err)
	}

	return wallet, nil
}

func (s *walletService) RefundCredits(ctx context.Context, userID uuid.UUID, amount int) error {

	err := s.walletRepo.AddCredits(ctx, userID, amount)
	if err != nil {
		s.logger.Error("Failed adding credits to wallet using wallet repository", zap.String("userID", userID.String()), zap.Int("amount", amount))
		return fmt.Errorf("repostiory failed to add credits to wallet: %w", err)
	}

	return nil
}

func (s *walletService) AddCredits(ctx context.Context, email string, amount int) error {
	err := s.walletRepo.AddCreditsToEmail(ctx, email, amount)
	if err != nil {
		if errors.Is(err, domain.ErrRecordNotFound) {
			s.logger.Error("Failed adding credits - user does not exist", zap.String("email", email), zap.Error(err))
		}
		s.logger.Error("Failed adding credits to wallet using repository", zap.String("email", email), zap.Error(err))
		return fmt.Errorf("failed adding credits to wallet using repository: %w", err)
	}

	return nil
}
