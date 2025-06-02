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
	Register(username, email, password string) (*domain.User, error)
	Login(username, password string) (*domain.User, string, error)
}

type authServiceImpl struct {
	appLogger *zap.Logger
	userRepo  port.UserRepository
}

func NewAuthService(userRepo port.UserRepository /*, tokenService port.TokenService */, logger *zap.Logger) AuthService {
	return &authServiceImpl{
		userRepo:  userRepo,
		appLogger: logger,
		// tokenService: tokenService
	}
}

func (s *authServiceImpl) Register(username, email, password string) (*domain.User, error) {

	if username == "" || email == "" || password == "" {
		return nil, errors.New("username, email and password are required")
	}

	existingUser, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error checking for existing user: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists") // TODO: Create custom domain error domain.ErrUserAlreadyExists
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("error checking for existing user: %w", err)
	}

	user := &domain.User{
		BaseModel: domain.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Don't return password, even if it is hashed
	user.Password = ""

	// TODO: Generate token, Update return to also return the token

	return user, nil
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
