package service

import (
	"context"

	"github.com/ivanpaghubasan/hoa-hub/internal/repository"
)

type UserServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{repo: repo}
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, user *CreateUserRequest) (*CreatUserResponse, error) {
	return nil, nil
}
