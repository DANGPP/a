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
func (t *AuthHandler) GetSecretKey(c *gin.Context) {
	var boddi struct {
		UUID string `json:"uuid"`
	}
	if err := c.ShouldBindJSON(&boddi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	secretKey, err := t.svc.Checksecretkey(boddi.UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get secret key", "details": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"uuid":      boddi.UUID,
		"secretKey": secretKey,
	})
}

// 1. tạo token mới
func (t *AuthHandler) RegisterToken(c *gin.Context) {
	var body AuthService.Bodi

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	token, secretkey, err := t.svc.CreateToken(body, c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token":     token,
		"secretKey": secretkey,
	})
}

// 2. thu hồi token uuid
func (t *AuthHandler) RevokeToken(c *gin.Context) {
	var uuid struct {
		UUID string `json:"uuid"`
	}
	if err := c.ShouldBindJSON(&uuid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	status, err := t.svc.RevokeToken(uuid.UUID, c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke token", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": status})

}

// 3. thu hồi full token
func (t *AuthHandler) RevokeTokenFull(c *gin.Context) {
	status, err := t.svc.RevokeTokenFull(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke all tokens", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": status})
}

// 4. Active toàn bộ token
func (t *AuthHandler) ActiveTokenFull(c *gin.Context) {
	status, err := t.svc.ActiveTokenFull(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate all tokens", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": status})
}

// 5. Lấy toàn bộ token
func (t *AuthHandler) GetAllToken(c *gin.Context) {
	token, err := t.svc.GetAllToken(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all tokens", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
