package service

import (
	"context"

	"github.com/ivanpaghubasan/hoa-hub/internal/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, user *CreateUserRequest) (*CreatUserResponse, error)
}

type Service struct {
	UserService UserService
}

type CreateUserRequest struct {
}

type CreatUserResponse struct {
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		UserService: NewUserService(repos.UserRepository),
	}
}
