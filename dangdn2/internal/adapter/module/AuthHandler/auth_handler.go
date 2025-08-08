package aut_http

import (
	"net/http"
	// "time"

	"test/internal/core/module/AuthService"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	svc *auth_service.TokenService
}

func NewHTTPHandler(svc *auth_service.TokenService) *HTTPHandler {
	return &HTTPHandler{svc: svc}
}

func (h *HTTPHandler) IssueToken(c *gin.Context) {
	// TODO: Check auth của admin ở đây

	var req struct {
		Service string `json:"service"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service is required"})
		return
	}

	token, err := h.svc.GenToken(req.Service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
