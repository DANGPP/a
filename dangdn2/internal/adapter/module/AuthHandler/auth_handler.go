package AuthHandler

import (
	"net/http"
	"test/internal/core/module/AuthService"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *AuthService.AuthService
}

func NewAuthHandler(svc *AuthService.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (t *AuthHandler) RegisterToken(c *gin.Context) {
	var body AuthService.Bodi

	c.ShouldBindJSON(&body)
	token, _ := t.svc.CreateToken(body, c.Request.Context())

	c.JSON(http.StatusOK, gin.H{"token": token})
}
