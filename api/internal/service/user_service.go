package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ivanpaghubasan/hoa-hub/internal/model"
	"github.com/ivanpaghubasan/hoa-hub/internal/repository"
	"github.com/ivanpaghubasan/hoa-hub/internal/util"
)

type UserServiceImpl struct {
	userRepo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{userRepo: repo}
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreatUserResponse, error) {
	result, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if result != nil {
		return nil, repository.RecordExistsError
	}

	var dateOfBirth *time.Time
	if req.DateOfBirth != "" {
		t, err := time.Parse(repository.DateFormat, req.DateOfBirth)
		if err != nil {
			return nil, fmt.Errorf("invalid date of birth format: %w", err)
		}
		dateOfBirth = &t
	}

	hashPassword, err := util.HashPassword(req.Password)
	if err != nil {

	}

	user := &model.User{
		ID:           uuid.New(),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		MiddleName:   req.MiddleName,
		DateOfBirth:  dateOfBirth,
		MobileNumber: req.MobileNumber,
		Gender:       req.Gender,
		Email:        req.Email,
		PasswordHash: hashPassword,
		Status:       repository.ActiveStatus,
		CreatedAt:    time.Now(),
	}

	resp, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &CreatUserResponse{
		FirstName:  resp.FirstName,
		LastName:   resp.LastName,
		MiddleName: resp.MiddleName,
		Email:      resp.Email,
	}, nil
}
