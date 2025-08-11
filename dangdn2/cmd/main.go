package main

import (
	"context"
	"os"
	"log"

	"test/internal/adapter/module/AuthHandler"
	"test/internal/core/module/AuthService"

	"github.com/gin-gonic/gin"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// lay secret Ky
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		secretKey = "AreUGay"
	}
	dbURL:= os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL =  "postgres://postgres:1@localhost:5433/auth?sslmode=disable"
	}

	//connect db
	pool,err := pgxpool.New(context.Background(),dbURL)
	if err != nil {
		log.Fatal("error connect gay:", err)
	}
	defer pool.Close()

	svc := AuthService.NewAuthService(secretKey, pool)
	h := AuthHandler.NewAuthHandler(svc)
	r := gin.Default()
	r.POST("/api/register", h.RegisterToken)
	r.Run(":8080")
}
