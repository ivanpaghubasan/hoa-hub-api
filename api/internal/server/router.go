package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.GET("health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
		})
	}
	return r
}
