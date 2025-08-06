// adapter/handler/http.go
package handler

import (
	"net/http"

	"dangdn2-test-go/internal/core/service"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	svc service.HelmService
}

func NewHTTPHandler(svc service.HelmService) *HTTPHandler {
	return &HTTPHandler{svc: svc}
}

func (h *HTTPHandler) Deploy(c *gin.Context) {
	namespace := c.Param("namespace")
	releaseName := c.Param("release")
	result, err := h.svc.Deploy(namespace, releaseName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *HTTPHandler) Update(c *gin.Context) {
	namespace := c.Param("namespace")
	release := c.Param("release")
	result, err := h.svc.Update(namespace, release)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *HTTPHandler) Delete(c *gin.Context) {
	namespace := c.Param("namespace")
	release := c.Param("release")
	result, err := h.svc.Delete(namespace, release)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
