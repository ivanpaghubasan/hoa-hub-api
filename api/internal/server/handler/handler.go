package handler

import "github.com/ivanpaghubasan/hoa-hub/internal/service"

type Handler struct {
	UserHandler *UserHandler
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		UserHandler: NewUserHandler(services.UserService),
	}
}
