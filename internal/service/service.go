package service

import (
	"context"

	"github.com/ivanpaghubasan/hoa-hub-api/internal/auth"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, user *CreateUserRequest) (*CreatUserResponse, error)
}

type Service struct {
	UserService UserService
}

type CreateUserRequest struct {
	FirstName    string  `json:"firstName" binding:"required"`
	LastName     string  `json:"lastName" binding:"required"`
	MiddleName   *string `json:"middleName"`
	Email        string  `json:"email" binding:"required,email"`
	Password     string  `json:"password" binding:"required,min=8"`
	DateOfBirth  string  `json:"dateOfBirth"`
	MobileNumber string  `json:"mobileNumber"`
	Gender       string  `json:"gender"`
	Status       string  `json:"status"`
	//RoleID       string  `json:"roleId" binding:"required"`
}

type CreatUserResponse struct {
	FirstName  string  `json:"firstName"`
	LastName   string  `json:"lastName"`
	MiddleName *string `json:"middleName"`
	Email      string  `json:"email"`
	UserType   string  `json:"userType"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserResponse struct {
	TokenClaims *auth.Claims
}

func NewService(repos *repository.Repository, jwt auth.IJWTAuth) *Service {
	return &Service{
		UserService: NewUserService(repos.UserRepository, jwt),
	}
}
