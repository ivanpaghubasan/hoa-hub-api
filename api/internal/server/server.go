package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ivanpaghubasan/hoa-hub/internal/config"
)

type Server struct {
	Port   string
	Engine *gin.Engine
}

func New(services interface{}, cfg *config.Config) *Server {
	router := NewRouter()
	return &Server{
		Port:   cfg.Port,
		Engine: router,
	}
}

func (s *Server) Run() error {
	fmt.Printf("Starting server at port %s\n", s.Port)
	return s.Engine.Run(":" + s.Port)
}
