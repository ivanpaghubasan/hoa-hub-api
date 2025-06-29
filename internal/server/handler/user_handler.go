package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/auth"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/service"
)

type UserHandler struct {
	userService service.UserService
	auth        auth.IJWTAuth
}

func NewUserHandler(service service.UserService, auth auth.IJWTAuth) *UserHandler {
	return &UserHandler{
		userService: service,
		auth:        auth,
	}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var request service.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.userService.CreateUser(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": response})
}

func (h *UserHandler) Login(c *gin.Context) {
	var request service.LoginUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userResp, err := h.userService.LoginUser(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	tokenPairs, err := h.auth.GenerateToken(userResp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshCookie := h.auth.GetRefreshCookie(tokenPairs.RefreshToken)
	http.SetCookie(c.Writer, refreshCookie)

	c.JSON(http.StatusOK, gin.H{
		"access_token": tokenPairs.AccessToken,
		"user": gin.H{
			"id":         userResp.ID.String(),
			"first_name": userResp.FirstName,
			"last_name":  userResp.LastName,
			"email":      userResp.Email,
		},
	})
}
