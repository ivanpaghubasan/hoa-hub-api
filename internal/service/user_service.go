package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ivanpaghubasan/hoa-hub-api/internal/auth"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/constants"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/model"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/repository"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/util"
)

type UserServiceImpl struct {
	userRepo repository.UserRepository
	jwt      auth.IJWTAuth
}

func NewUserService(repo repository.UserRepository, jwt auth.IJWTAuth) UserService {
	return &UserServiceImpl{
		userRepo: repo,
		jwt:      jwt,
	}
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreatUserResponse, error) {
	result, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if result != nil {
		return nil, constants.ErrRecordExists
	}

	var dateOfBirth *time.Time
	if req.DateOfBirth != "" {
		t, err := time.Parse(constants.DateFormat, req.DateOfBirth)
		if err != nil {
			return nil, fmt.Errorf("invalid date of birth format: %w", err)
		}
		dateOfBirth = &t
	}

	hashPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	user := &model.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		MiddleName:   req.MiddleName,
		DateOfBirth:  dateOfBirth,
		MobileNumber: req.MobileNumber,
		Gender:       req.Gender,
		Email:        req.Email,
		PasswordHash: hashPassword,
		Status:       constants.ActiveStatus,
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

func (s *UserServiceImpl) LoginUser(ctx context.Context, req *LoginUserRequest) (*LoginUserResponse, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, constants.ErrRecordNotFound
	}

	isValid := util.VerifyPasswordHash(req.Password, user.PasswordHash)
	if !isValid {
		return nil, constants.ErrInvalidPassword
	}

	// Generate token
	token, err := s.jwt.GenerateToken(user)
	if err != nil {
		return nil, constants.ErrInternalServer
	}
	_ = token
	return &LoginUserResponse{}, nil
}
