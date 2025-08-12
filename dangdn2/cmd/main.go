package main

import (
	"context"
	"log"
	"os"

	"test/internal/adapter/module/AuthHandler"
	"test/internal/adapter/module/anotherAdapter"
	"test/internal/core/module/AuthService"
	"test/internal/core/module/anotherService"

	"github.com/gin-gonic/gin"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	//load .env
	godotenv.Load(".env")
	// lay secret Ky
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		secretKey = "AreUggggGay"
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "?????????postgres://postgres:1@localhost:5433/auth?sslmode=disable"
	}
	vaultAddr := os.Getenv("VAULT_ADDR")
	if vaultAddr == "" {
		vaultAddr = "http://127.0.0.1:8205"
	}

	vaultToken := os.Getenv("VAULT_TOKEN")
	if vaultToken == "" {
		// log.Fatal("VAULT_TOKEN không được để trống")
		vaultToken = "REDACTED"
	}

	vaultPath := os.Getenv("VAULT_PATH")
	if vaultPath == "" {
		vaultPath = "mon_deployment_secretkey"
	}

	//connect db
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("error connect gay:", err)
	}
	defer pool.Close()

	svcA := anotherService.NewAnotherService(vaultAddr, vaultToken, vaultPath)
	hA := anotherAdapter.NewAnotherAdapter(svcA)

	svc := AuthService.NewAuthService(secretKey, pool)
	h := AuthHandler.NewAuthHandler(svc)
	r := gin.Default()
	r.POST("/api/register", h.RegisterToken)
	r.PUT("/api/revoketoke", h.RevokeToken)
	r.PUT("/api/revoketokenfull", h.RevokeTokenFull)
	r.PUT("/api/activetokenfull", h.ActiveTokenFull)
	r.GET("api/secretkey", hA.GetSecretKey)
	r.GET("api/dburl", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"dburl": dbURL}) })

	r.POST("/api/secretkey", hA.GenSecretKey)
	r.Run(":8080")
}
