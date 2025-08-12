package anotherAdapter

import (
	"net/http"
	"test/internal/core/module/anotherService"

	"github.com/gin-gonic/gin"
)

type AnotherAdapter struct {
	svc *anotherService.Another
}

func NewAnotherAdapter(svc *anotherService.Another) *AnotherAdapter {
	return &AnotherAdapter{svc: svc}
}

func (a *AnotherAdapter) GenSecretKey(c *gin.Context) {
	var bodi struct {
		Length int `json:"length"`
	}

	c.ShouldBindJSON(&bodi)

	IDkey, SecretKey, _ := a.svc.GenSecretKey(bodi.Length)
	c.JSON(http.StatusOK, gin.H{
		"KeyID":  IDkey,
		"key ne": SecretKey})
}

func (a *AnotherAdapter) GetSecretKey(c *gin.Context) {
	var bodi struct {
		UUIDKeySecret string `json:"uuidsecret"`
	}

	c.ShouldBindJSON(&bodi)

	SecretKey, _ := a.svc.GetSecretKey(bodi.UUIDKeySecret)
	c.JSON(http.StatusOK, gin.H{
		// "KeyID":  IDkey,
		"key ne": SecretKey})
}
