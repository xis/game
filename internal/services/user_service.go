package service

import (
	"context"
	"errors"

	"game/internal/domain"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrUsernameExists  = errors.New("username exists")
)

//go:generate mockery --name UserService --structname MockUserService --outpkg mocks --filename user_service_mock.go --output ./mocks/. --with-expecter
type UserService interface {
	Login(ctx context.Context, username string, password string) (LoginResult, error)
	Register(ctx context.Context, username string, password string) (domain.User, error)
}

type LoginResult struct {
	UserID   string
	UserName string
	Token    string
}

type UserServiceDependencies struct {
	UserRepository      domain.UserRepository
	UserScoreRepository domain.UserScoreRepository
	TokenManager        domain.TokenManager
	PasswordHasher      domain.PasswordHasher
}

type userService struct {
	userRepository      domain.UserRepository
	userScoreRepository domain.UserScoreRepository
	tokenManager        domain.TokenManager
	passwordHasher      domain.PasswordHasher
}

func NewUserService(
	deps UserServiceDependencies,
) *userService {
	return &userService{
		userRepository:      deps.UserRepository,
		userScoreRepository: deps.UserScoreRepository,
		tokenManager:        deps.TokenManager,
		passwordHasher:      deps.PasswordHasher,
	}
}

func (service *userService) Login(ctx context.Context, username string, password string) (LoginResult, error) {
	user, err := service.userRepository.GetByName(ctx, username)
	if err != nil {
		return LoginResult{}, err
	}

	isMatch, err := service.passwordHasher.ComparePasswordAndHash(password, user.PasswordHash)
	if err != nil {
		return LoginResult{}, err
	}

	if !isMatch {
		return LoginResult{}, ErrInvalidPassword
	}

	token, err := service.tokenManager.Create(ctx, user.ID)
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{
		UserID:   user.ID,
		UserName: user.Name,
		Token:    token,
	}, nil
}

func (service *userService) Register(ctx context.Context, username string, password string) (domain.User, error) {
	exists, err := service.userRepository.CheckExistsByName(ctx, username)
	if err != nil {
		return domain.User{}, err
	}

	if exists {
		return domain.User{}, ErrUsernameExists
	}

	hashedPassword, err := service.passwordHasher.HashPassword(password)
	if err != nil {
		return domain.User{}, err
	}

	user, err := service.userRepository.Create(ctx, username, hashedPassword)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
