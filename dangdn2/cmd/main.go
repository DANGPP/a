package main

import (
	// "context"
	"fmt"
	"log"
	"os"

	"test/internal/adapter/module/AuthHandler"
	"test/internal/adapter/module/anotherAdapter"
	"test/internal/core/module/AuthService"
	"test/internal/core/module/anotherService"

	"github.com/gin-gonic/gin"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("gay vl")
	}
	dbURL := os.Getenv("DATABASE_URL") //"host=localhost user=postgres password=1 dbname=auth port=5433 sslmode=disable"

	vaultAddr := os.Getenv("VAULT_ADDR") // "http://127.0.0.1:8205"

	vaultToken := os.Getenv("VAULT_TOKEN") // 

	vaultPath := os.Getenv("VAULT_PATH") //"fmon_deployment_secretkey"

	//connect db
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal("error connect gay:", err)
	}

	svcA := anotherService.NewAnotherService(vaultAddr, vaultToken, vaultPath)
	hA := anotherAdapter.NewAnotherAdapter(svcA)
	// uuid := "1b01cf34-0cdb-4e07-aac5-488ff894d050"
	// secretKey, _ := svcA.GetSecretKey(uuid)

	svc := AuthService.NewAuthService(svcA, db)
	h := AuthHandler.NewAuthHandler(svc)
	r := gin.Default()
	r.POST("/api/register", h.RegisterToken)
	r.GET("/api/fulltoken", h.GetAllToken)
	r.PUT("/api/revoketoke", h.RevokeToken)
	r.PUT("/api/revoketokenfull", h.RevokeTokenFull)
	r.PUT("/api/activetokenfull", h.ActiveTokenFull)
	r.GET("api/secretkey", hA.GetSecretKey)
	r.GET("api/secretkey2", h.GetSecretKey)
	r.GET("api/dburl", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"dburl": dbURL}) })

	r.POST("/api/secretkey", hA.GenSecretKey)
	r.Run(":8080")
}
