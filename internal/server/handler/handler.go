package handler

import (
	"github.com/ivanpaghubasan/hoa-hub-api/internal/auth"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/service"
)

type Handler struct {
	UserHandler *UserHandler
	Auth        auth.IJWTAuth
}

func NewHandler(services *service.Service, auth auth.IJWTAuth) *Handler {
	return &Handler{
		UserHandler: NewUserHandler(services.UserService, auth),
	}
}
