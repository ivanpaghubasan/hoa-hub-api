package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/ivanpaghubasan/hoa-hub/internal/model"
	"github.com/ivanpaghubasan/hoa-hub/internal/repository"
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

func TestCreateUser(t *testing.T) {
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
						return nil, repository.RecordNotFoundError
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
			fmt.Println(resp)
			if tc.expectErr {
				if err == nil {
					t.Errorf("expected error, got none")
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
