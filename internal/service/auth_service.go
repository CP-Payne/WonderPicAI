package service

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/CP-Payne/wonderpicai/internal/config"
	"github.com/CP-Payne/wonderpicai/internal/domain"
	"github.com/CP-Payne/wonderpicai/internal/port"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthService interface {
	Register(username, email, password string) (*domain.User, string, error)
	Login(email, password string) (*domain.User, string, error)
	HandleExternalAuthCallback(r *http.Request) (*domain.User, string, error)
}

type authServiceImpl struct {
	logger       *zap.Logger
	userRepo     port.UserRepository
	tokenService port.TokenService
	externalAuth port.ExternalAuthService
}

func NewAuthService(userRepo port.UserRepository, tokenService port.TokenService, logger *zap.Logger, externalAuth port.ExternalAuthService) AuthService {
	return &authServiceImpl{
		userRepo:     userRepo,
		logger:       logger.With(zap.String("component", "AuthService")),
		tokenService: tokenService,
		externalAuth: externalAuth,
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

	s.logger.Info("User creation successfull", zap.String("email", userToCreate.Email), zap.String("UserID", userToCreate.ID.String()))

	claims := jwt.MapClaims{
		"sub": userToCreate.ID,
		"exp": time.Now().Add(time.Duration(config.Cfg.JWT.ExpiryMinutes) * time.Minute).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": config.Cfg.JWT.Issuer,
		// Same as issuer in this implementation
		"aud": config.Cfg.JWT.Issuer,
	}

	token, err := s.tokenService.GenerateToken(claims)
	if err != nil {
		s.logger.Error("Failed to create JWT token after successfull user registration",
			zap.String("email", userToCreate.Email),
			zap.String("UserID", userToCreate.ID.String()),
			zap.Error(err),
		)
		return userToCreate, "", fmt.Errorf("failed to generate jwt token via local token service: %w", err)
	}

	return userToCreate, token, nil
}

func (s *authServiceImpl) Login(email, password string) (*domain.User, string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, "", domain.ErrInvalidCredentials
		}
		return nil, "", fmt.Errorf("failed authenticating user: %w", err)
	}

	if ok := checkPasswordHash(password, user.Password); !ok {
		return nil, "", domain.ErrInvalidCredentials
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Duration(config.Cfg.JWT.ExpiryMinutes) * time.Minute).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": config.Cfg.JWT.Issuer,
		// Same as issuer in this implementation
		"aud": config.Cfg.JWT.Issuer,
	}

	token, err := s.tokenService.GenerateToken(claims)
	if err != nil {
		s.logger.Error("Failed to create JWT token after successfull user authentication",
			zap.String("email", user.Email),
			zap.String("UserID", user.ID.String()),
			zap.Error(err),
		)
		return user, "", fmt.Errorf("failed to generate jwt token via local token service: %w", err)
	}

	return user, token, nil
}

// TODO: Implement Email verification

// --- Developer Note ---
// Without email verification, a user can signup with any email.
// If a another user owns this email or creates this email (e.g. on gmail) in the future
// They would be able to authenticate using google without the password.
// As such, it is important that email verification is implemented to ensure that the user that signsup
// owns the email.
func (s *authServiceImpl) HandleExternalAuthCallback(r *http.Request) (*domain.User, string, error) {
	externalUser, err := s.externalAuth.HandleCallback(r)
	if err != nil {
		return nil, "", fmt.Errorf("failed to authenticate user through external provider: %w", err)
	}

	if externalUser.Email == "" || externalUser.Name == "" {
		return nil, "", fmt.Errorf("failed to handle external auth: %w", errors.New("empty username or password"))
	}

	user, err := s.userRepo.GetByEmail(externalUser.Email)

	if errors.Is(err, domain.ErrUserNotFound) {
		user = &domain.User{
			BaseModel: domain.BaseModel{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Username: externalUser.Name,
			Email:    externalUser.Email,
		}

		if err := s.userRepo.Create(user); err != nil {
			s.logger.Error("Failed to create user via repository", zap.Error(err))
			return nil, "", fmt.Errorf("failed to create user using repository: %w", err)
		}
	} else if err != nil {
		return nil, "", fmt.Errorf("failed authenticating user: %w", err)
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Duration(config.Cfg.JWT.ExpiryMinutes) * time.Minute).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": config.Cfg.JWT.Issuer,
		// Same as issuer in this implementation
		"aud": config.Cfg.JWT.Issuer,
	}

	token, err := s.tokenService.GenerateToken(claims)
	if err != nil {
		s.logger.Error("Failed to create JWT token after successfull user authentication",
			zap.String("email", user.Email),
			zap.String("UserID", user.ID.String()),
			zap.Error(err),
		)
		return user, "", fmt.Errorf("failed to generate jwt token via local token service: %w", err)
	}

	return user, token, nil

}
