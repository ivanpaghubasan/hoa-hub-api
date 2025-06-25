package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/server/handler"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/service"
)

func NewRouter(services *service.Service) *gin.Engine {
	r := gin.Default()

	handler := handler.NewHandler(services)
	_ = handler
	v1 := r.Group("/v1")
	{
		v1.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
		})

	}
	return r
}
