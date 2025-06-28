package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/auth"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/constants"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/model"
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

type MockJWTManager struct {
	GenerateTokenFn     func(userID uuid.UUID) (string, error)
	GenerateTokenCalled bool
	ParseTokenFn        func(tokenStr string) (*auth.Claims, error)
	ParseTokenCalled    bool
}

func (m *MockJWTManager) GenerateToken(userID uuid.UUID) (string, error) {
	m.GenerateTokenCalled = true
	return m.GenerateTokenFn(userID)
}

func (m *MockJWTManager) ParseToken(tokenStr string) (*auth.Claims, error) {
	m.ParseTokenCalled = true
	return m.ParseTokenFn(tokenStr)
}

func setupJWTManagerMock() *MockJWTManager {
	return &MockJWTManager{
		GenerateTokenFn: func(userID uuid.UUID) (string, error) {
			return "token12345", nil
		},
		ParseTokenFn: func(tokenStr string) (*auth.Claims, error) {
			return &auth.Claims{
				UserID: "12345",
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			}, nil
		},
	}
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
			mockJwtManager := setupJWTManagerMock()
			service := NewUserService(mockRepo(), mockJwtManager)

			resp, err := service.CreateUser(context.Background(), tc.req)

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
