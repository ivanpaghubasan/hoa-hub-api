package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ivanpaghubasan/hoa-hub-api/internal/constants"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/model"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/util"
)

type MockUserRepository struct {
	GetUserByEmailFn func(ctx context.Context, email string) (*model.User, error)
	CreateUserFn     func(ctx context.Context, user *model.User) (*model.User, error)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return m.GetUserByEmailFn(ctx, email)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	return m.CreateUserFn(ctx, user)
}

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name        string
		req         *CreateUserRequest
		setupMock   func() *MockUserRepository
		expectErr   bool
		expectEmail string
	}{
		{
			name: "success create user",
			req: func() *CreateUserRequest {
				middleName := "Test"
				return &CreateUserRequest{
					FirstName:    "John",
					LastName:     "Doe",
					MiddleName:   &middleName,
					Email:        "john.doe@example.com",
					Password:     "password12345",
					DateOfBirth:  "2000-01-01",
					MobileNumber: "09123456789",
					Gender:       "Male",
				}
			}(),
			setupMock: func() *MockUserRepository {
				return &MockUserRepository{
					GetUserByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
						return nil, nil
					},
					CreateUserFn: func(ctx context.Context, user *model.User) (*model.User, error) {
						return user, nil
					},
				}
			},
			expectErr:   false,
			expectEmail: "john.doe@example.com",
		},
		{
			name: "existing email error",
			req: func() *CreateUserRequest {
				middleName := "test"
				return &CreateUserRequest{
					FirstName:    "John",
					LastName:     "Doe",
					MiddleName:   &middleName,
					Email:        "john.doe@example.com",
					Password:     "password12345",
					DateOfBirth:  "2000-01-01",
					MobileNumber: "09123456789",
					Gender:       "Male",
				}
			}(),
			setupMock: func() *MockUserRepository {
				return &MockUserRepository{
					GetUserByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
						return &model.User{Email: email}, nil
					},
					CreateUserFn: func(ctx context.Context, user *model.User) (*model.User, error) {
						return nil, errors.New("should not be called")
					},
				}
			},
			expectErr: true,
		},
		{
			name: "invalid date of birth",
			req: func() *CreateUserRequest {
				middleName := "test"
				return &CreateUserRequest{
					FirstName:    "John",
					LastName:     "Doe",
					MiddleName:   &middleName,
					Email:        "john.doe@example.com",
					Password:     "password12345",
					DateOfBirth:  "invalid-date",
					MobileNumber: "09123456789",
					Gender:       "Male",
				}
			}(),
			setupMock: func() *MockUserRepository {
				return &MockUserRepository{
					GetUserByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
						return nil, constants.ErrRecordNotFound
					},
					CreateUserFn: func(ctx context.Context, user *model.User) (*model.User, error) {
						return user, nil
					},
				}
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := tc.setupMock
			service := NewUserService(mockRepo())

			resp, err := service.CreateUser(context.Background(), tc.req)
			if tc.expectErr {
				if err == nil {
					t.Error("expected error, got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if resp.Email != tc.expectEmail {
					t.Errorf("expected email %s, got %s", tc.expectEmail, resp.Email)
				}
			}

		})
	}
}

func TestUserService_LoginUser(t *testing.T) {
	password := "password12345"
	hashPassword, _ := util.HashPassword(password)

	tests := []struct {
		name        string
		req         *LoginUserRequest
		setupMock   func() *MockUserRepository
		expectErr   bool
		expectedErr error
	}{
		{
			name: "successful login",
			req: &LoginUserRequest{
				Email:    "john.doe@test.com",
				Password: password,
			},
			setupMock: func() *MockUserRepository {
				return &MockUserRepository{
					GetUserByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
						return &model.User{Email: email, PasswordHash: hashPassword}, nil
					},
				}
			},
			expectErr:   false,
			expectedErr: nil,
		},
		{
			name: "user not found",
			req: &LoginUserRequest{
				Email:    "john@test.com",
				Password: password,
			},
			setupMock: func() *MockUserRepository {
				return &MockUserRepository{
					GetUserByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
						return nil, constants.ErrRecordNotFound
					},
				}
			},
			expectErr:   true,
			expectedErr: constants.ErrRecordNotFound,
		},
		{
			name: "invalid password",
			req: &LoginUserRequest{
				Email:    "john@test.com",
				Password: "wron_password123",
			},
			setupMock: func() *MockUserRepository {
				return &MockUserRepository{
					GetUserByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
						return &model.User{Email: email, PasswordHash: hashPassword}, nil
					},
				}
			},
			expectErr:   true,
			expectedErr: constants.ErrInvalidPassword,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := tc.setupMock()
			service := NewUserService(mockRepo)
			user, err := service.LoginUser(context.Background(), tc.req)
			if tc.expectErr {
				if err == nil {
					t.Error("expected error, got none")
				} else if err != tc.expectedErr {
					t.Errorf("expected error %v, got %v", tc.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if user.Email != tc.req.Email {
					t.Errorf("expected email %s, got %s", tc.req.Email, user.Email)
				}
			}
		})
	}
}
