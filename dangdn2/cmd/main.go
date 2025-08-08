package main

import (
	// "log"
	// "os"

	auth_http "test/internal/adapter/module/AuthHandler"
	auth_service "test/internal/core/module/AuthService"

	"github.com/gin-gonic/gin"
	// "github.com/jmoiron/sqlx"
	// _ "github.com/lib/pq"
	// "github.com/redis/go-redis/v9"
)

func main() {
	r := gin.Default()

	// Tạo service & handler
	svc := auth_service.NewTokenService("my_super_secret_key")
	h := auth_http.NewHTTPHandler(svc)

	// Route
	r.POST("/api/register", h.IssueToken)

	r.Run(":8080")
}
