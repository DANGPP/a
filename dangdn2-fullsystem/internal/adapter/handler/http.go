package handler

import (
	"net/http"
	"test/internal/core/service"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	svc service.HelmService
}

func NewHttpHandler(svc service.HelmService) *HTTPHandler {
	return &HTTPHandler{svc: svc}
}

func (h *HTTPHandler) Deploy(c *gin.Context) {
	namespace := c.Param("namespace")
	releaseName := c.Param("release")

	result, err := h.svc.DeployRelease(namespace, releaseName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *HTTPHandler) Update(c *gin.Context) {
	namespace := c.Param("namespace")
	releaseName := c.Param("release")

	result, err := h.svc.UpdateRelease(namespace, releaseName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *HTTPHandler) Delete(c *gin.Context) {
	namespace := c.Param("namespace")
	releaseName := c.Param("release")

	result, err := h.svc.DeleteRelease(namespace, releaseName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *HTTPHandler) GetRelease(c *gin.Context) {
	namespace := c.Param("namespace")
	releaseName := c.Param("release")

	result, err := h.svc.GetRelease(namespace, releaseName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
