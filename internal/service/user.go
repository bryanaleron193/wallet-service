package service

import (
	"context"
	"fmt"

	"github.com/bryanaleron193/wallet-service/internal/model"
	"github.com/bryanaleron193/wallet-service/internal/repository"
	"github.com/bryanaleron193/wallet-service/pkg/apperror"
)

type UserService interface {
	CreateUser(ctx context.Context, username, fullName, email string) (*model.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(ur repository.UserRepository) UserService {
	return &userService{
		userRepo: ur,
	}
}

func (s *userService) CreateUser(ctx context.Context, username, fullName, email string) (*model.User, error) {
	if username == "" || email == "" {
		return nil, fmt.Errorf("username and email are required: %w", apperror.ErrInvalidInput)
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
