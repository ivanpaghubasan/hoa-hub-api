package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/auth"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/config"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/service"
)

type Server struct {
	Port   string
	Engine *gin.Engine
}

func New(services *service.Service, cfg *config.Config, jwt auth.IJWTAuth) *Server {
	router := NewRouter(services, jwt)

	return &Server{
		Port:   cfg.Port,
		Engine: router,
	}
}

func (s *Server) Run() error {
	fmt.Printf("Starting server at port %s\n", s.Port)
	return s.Engine.Run(":" + s.Port)
}
