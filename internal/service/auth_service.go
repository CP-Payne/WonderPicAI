package service

import (
	"errors"
	"fmt"

	"github.com/CP-Payne/wonderpicai/internal/domain"
	"github.com/CP-Payne/wonderpicai/internal/port"
)

type AuthService interface {
	Register(username, password string) (*domain.User, error)
	Login(username, password string) (*domain.User, string, error)
}

type authServiceImpl struct {
	userRepo port.UserRepository
}

func NewAuthService(userRepo port.UserRepository /*, tokenService port.TokenService */) AuthService {
	return &authServiceImpl{
		userRepo: userRepo,
		// tokenService: tokenService
	}
}

func (s *authServiceImpl) Register(username, password string) (*domain.User, error) {

	if username == "" || password == "" {
		return nil, errors.New("username and password are required")
	}

	existingUser, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("error checking for existing user: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("user with this username already exists") // TODO: Create custom domain error domain.ErrUserAlreadyExists
	}

	// TODO: Hash password

	user := &domain.User{
		Username: username,
		Password: password, // TODO: Store hashed password
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Don't return password, even if it is hashed
	user.Password = ""
	return user, nil
}

func (s *authServiceImpl) Login(username, password string) (*domain.User, string, error) {
	if username == "" || password == "" {
		return nil, "", errors.New("username and password are required")
	}

	user, err := s.userRepo.GetByUsername((username))
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
