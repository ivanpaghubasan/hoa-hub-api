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
	FirstName  string  `json:"first_name" binding:"reqiured"`
	LastName   string  `json:"last_name" binding:"required"`
	MiddleName *string `json:"middle_name"`
	Email      string  `json:"email" binding:"required,email"`
	Password   string  `json:"password" binding:"required,min=8"`
	UserType   string  `json:"user_type" binding:"required"`
}

type CreatUserResponse struct {
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	MiddleName *string `json:"middle_name"`
	Email      string  `json:"email"`
	UserType   string  `json:"user_type"`
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		UserService: NewUserService(repos.UserRepository),
	}
}
