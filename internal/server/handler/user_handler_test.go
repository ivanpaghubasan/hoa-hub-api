package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/auth"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/model"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/service"
)

type MockUserService struct {
	CreateUserFn func(ctx context.Context, req *service.CreateUserRequest) (*service.CreatUserResponse, error)
	LoginUserFn  func(ctx context.Context, req *service.LoginUserRequest) (*model.User, error)
}

func (m *MockUserService) CreateUser(ctx context.Context, req *service.CreateUserRequest) (*service.CreatUserResponse, error) {
	return m.CreateUserFn(ctx, req)
}

func (m *MockUserService) LoginUser(ctx context.Context, req *service.LoginUserRequest) (*model.User, error) {
	return m.LoginUserFn(ctx, req)
}

type MockJWTAuth struct {
	GenerateTokenFn        func(user *model.User) (auth.TokenPairs, error)
	GenerateTokenCalled    bool
	ParseAccessTokenFn     func(tokenStr string) (*auth.Claims, error)
	ParseAccessTokenCalled bool
}

func (m *MockJWTAuth) GenerateToken(user *model.User) (auth.TokenPairs, error) {
	m.GenerateTokenCalled = true
	return m.GenerateTokenFn(user)
}

func (m *MockJWTAuth) ParseAccessToken(tokenStr string) (*auth.Claims, error) {
	m.ParseAccessTokenCalled = true
	return m.ParseAccessTokenFn(tokenStr)
}

func (m *MockJWTAuth) GetRefreshCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{}
}

func (m *MockJWTAuth) GetExpiredRefreshCookie() *http.Cookie {
	return &http.Cookie{}
}

func setupJWTManagerMock() *MockJWTAuth {
	return &MockJWTAuth{
		GenerateTokenFn: func(user *model.User) (auth.TokenPairs, error) {
			return auth.TokenPairs{
				AccessToken:  "token12345",
				RefreshToken: "refresh12345",
			}, nil
		},
		ParseAccessTokenFn: func(tokenStr string) (*auth.Claims, error) {
			return &auth.Claims{
				UserID: "12345",
				Name:   "John Doe",
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			}, nil
		},
	}
}

func TestUserHandler_RegisterUser(t *testing.T) {
	tests := []struct {
		name           string
		payload        map[string]interface{}
		mockService    func() *MockUserService
		expectedStatus int
	}{
		{
			name: "success",
			payload: map[string]interface{}{
				"firstName":    "John",
				"lastName":     "Doe",
				"middleName":   "Test",
				"email":        "john.doe@example.com",
				"password":     "password123",
				"dateOfBirth":  "2000-01-23",
				"mobileNumber": "09123456789",
				"gender":       "Male",
			},
			mockService: func() *MockUserService {
				return &MockUserService{
					CreateUserFn: func(ctx context.Context, req *service.CreateUserRequest) (*service.CreatUserResponse, error) {
						return &service.CreatUserResponse{Email: req.Email}, nil
					},
				}
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "bad request body",
			payload: map[string]interface{}{
				"firstName": 123,
			},
			mockService: func() *MockUserService {
				return &MockUserService{
					CreateUserFn: func(ctx context.Context, req *service.CreateUserRequest) (*service.CreatUserResponse, error) {
						return nil, errors.New("some error")
					},
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			payload: map[string]interface{}{
				"firstName":    "John",
				"lastName":     "Doe",
				"email":        "john.doe@example.com",
				"password":     "password123",
				"mobileNumber": "09123456789",
				"gender":       "Male",
			},
			mockService: func() *MockUserService {
				return &MockUserService{
					CreateUserFn: func(ctx context.Context, req *service.CreateUserRequest) (*service.CreatUserResponse, error) {
						return nil, errors.New("some error")
					},
				}
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			r := gin.Default()
			mockJwtManager := setupJWTManagerMock()

			mockSvc := tc.mockService()
			h := NewUserHandler(mockSvc, mockJwtManager)

			r.POST("/register", h.RegisterUser)

			body, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, rec.Code)
			}
		})
	}
}
