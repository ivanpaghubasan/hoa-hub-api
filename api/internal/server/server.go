package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ivanpaghubasan/hoa-hub/internal/config"
)

type Server struct {
	Addr   string
	Engine *gin.Engine
}

func New(services interface{}, cfg *config.Config) *Server {
	router := NewRouter()
	return &Server{

		Engine: router,
	}
}

func (s *Server) Run() error {
	fmt.Printf("Starting server at %s\n", s.Addr)
	return s.Engine.Run(":" + s.Addr)
}
