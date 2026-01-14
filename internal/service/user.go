package service

import (
	"context"
	"fmt"

	"github.com/bryanaleron193/wallet-service/internal/config"
	"github.com/bryanaleron193/wallet-service/internal/model"
	"github.com/bryanaleron193/wallet-service/internal/repository"
	"github.com/bryanaleron193/wallet-service/pkg/response"
	"github.com/bryanaleron193/wallet-service/pkg/util"
)

type UserService interface {
	CreateUser(ctx context.Context, username, fullName, email string) (*model.User, error)
	Login(ctx context.Context, username string) (string, error)
}

type userService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewUserService(ur repository.UserRepository, cfg *config.Config) UserService {
	return &userService{
		userRepo: ur,
		cfg:      cfg,
	}
}

func (s *userService) CreateUser(ctx context.Context, username, fullName, email string) (*model.User, error) {
	if username == "" || email == "" {
		return nil, fmt.Errorf("username and email are required: %w", response.ErrInvalidInput)
	}

	user := &model.User{
		Username: username,
		FullName: fullName,
		Email:    email,
	}

	err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("service.CreateUser: %w", err)
	}

	return user, nil
}

func (s *userService) Login(ctx context.Context, username string) (string, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	token, err := util.GenerateToken(user.ID, user.Username, s.cfg.JWT.Secret)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}
