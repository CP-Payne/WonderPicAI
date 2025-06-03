package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/CP-Payne/wonderpicai/internal/domain"
	"github.com/CP-Payne/wonderpicai/internal/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthService interface {
	Register(username, email, password string) (*domain.User, string, error)
	Login(username, password string) (*domain.User, string, error)
}

type authServiceImpl struct {
	logger   *zap.Logger
	userRepo port.UserRepository
}

func NewAuthService(userRepo port.UserRepository /*, tokenService port.TokenService */, logger *zap.Logger) AuthService {
	return &authServiceImpl{
		userRepo: userRepo,
		logger:   logger.With(zap.String("component", "AuthService")),
		// tokenService: tokenService
	}
}

func (s *authServiceImpl) Register(username, email, password string) (*domain.User, string, error) {
	existingUser, err := s.userRepo.GetByEmail(email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		s.logger.Error("Error checking existing user by email", zap.Error(err))
		return nil, "", fmt.Errorf("registration process failed: %w", err)
	}

	if existingUser != nil {
		return nil, "", domain.ErrEmailAlreadyExists
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		s.logger.Error("Password hashing failed", zap.Error(err))
		return nil, "", fmt.Errorf("user creation failed: %w", err)
	}

	userToCreate := &domain.User{
		BaseModel: domain.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	if err := s.userRepo.Create(userToCreate); err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) || errors.Is(err, domain.ErrDuplicateEntry) {
			s.logger.Warn("User creation conflict", zap.Error(err))
			return nil, "", err
		}
		s.logger.Error("Failed to create user via repository", zap.Error(err))

		return nil, "", fmt.Errorf("failed to complete registration due to an internal issue: %w", err)
	}

	// Don't return password, even if it is hashed
	userToCreate.Password = ""

	// TODO: Generate token, Update return to also return the token

	return userToCreate, "", nil
}

func (s *authServiceImpl) Login(email, password string) (*domain.User, string, error) {
	if email == "" || password == "" {
		return nil, "", errors.New("username and password are required")
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, "", fmt.Errorf("error fetching user: %w", err)
	}
	if user == nil {
		return nil, "", errors.New("invalid username or password")
	}

	// TODO: Compare hashed password

	if user.Password != password {
		return nil, "", errors.New("invalid username or password")
	}

	// TODO: generate jwt token using tokenservice

	token := "some dummy token until tokenservice implemented"

	return user, token, nil
}
