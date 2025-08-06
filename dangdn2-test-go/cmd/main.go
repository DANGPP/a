// cmd/main.go
package main

import (
	"dangdn2-test-go/internal/adapter/handler"
	"dangdn2-test-go/internal/core/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	svc, err := service.NewHelmService()
	if err != nil {
		log.Fatalf("failed to initialize Helm service: %v", err)
	}

	h := handler.NewHTTPHandler(svc)

	r := gin.Default()
	r.POST("/api/deploy/:namespace/:release", h.Deploy)
	r.POST("/api/update/:namespace/:release", h.Update)
	r.DELETE("/api/delete/:namespace/:release", h.Delete)

	r.Run(":8080")
}
