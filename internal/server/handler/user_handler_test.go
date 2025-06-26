package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/service"
)

type MockUserService struct {
	CreateUserFn func(ctx context.Context, req *service.CreateUserRequest) (*service.CreatUserResponse, error)
}

func (m *MockUserService) CreateUser(ctx context.Context, req *service.CreateUserRequest) (*service.CreatUserResponse, error) {
	return m.CreateUserFn(ctx, req)
}

func TestRegisterUser(t *testing.T) {
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

			mockSvc := tc.mockService()
			h := NewUserHandler(mockSvc)

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
