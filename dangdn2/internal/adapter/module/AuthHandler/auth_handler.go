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

// 1. tạo token mới
func (t *AuthHandler) RegisterToken(c *gin.Context) {
	var body AuthService.Bodi

	c.ShouldBindJSON(&body)
	token, _ := t.svc.CreateToken(body, c.Request.Context())

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// 2. thu hồi token uuid
func (t *AuthHandler) RevokeToken(c *gin.Context) {
	var uuid struct {
		UUID string `json:"uuid"`
	}
	c.ShouldBindJSON(&uuid)
	status, _ := t.svc.RevokeToken(uuid.UUID, c.Request.Context())

	c.JSON(http.StatusOK, gin.H{"status": status})

}

// 3. thu hồi full token
func (t *AuthHandler) RevokeTokenFull(c *gin.Context) {
	status, _ := t.svc.RevokeTokenFull(c.Request.Context())

	c.JSON(http.StatusOK, gin.H{"status": status})
}

// 4. active all token
func (t *AuthHandler) ActiveTokenFull(c *gin.Context) {
	status, _ := t.svc.ActiveTokenFull(c.Request.Context())

	c.JSON(http.StatusOK, gin.H{"status": status})
}

// 5. xem full token
func (t *AuthHandler) GetAllToken(c *gin.Context) {
	token, _ := t.svc.GetAllToken(c.Request.Context())

	c.JSON(http.StatusOK, gin.H{"token": token})
}
