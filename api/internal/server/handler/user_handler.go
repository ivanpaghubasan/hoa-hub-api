package handler

import "github.com/ivanpaghubasan/hoa-hub/internal/service"

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{userService: service}
}


