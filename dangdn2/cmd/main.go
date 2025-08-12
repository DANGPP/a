package main

import (
	// "context"
	"log"
	"os"

	"test/internal/adapter/module/AuthHandler"
	"test/internal/adapter/module/anotherAdapter"
	"test/internal/core/module/AuthService"
	"test/internal/core/module/anotherService"

	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	//load .env
	godotenv.Load(".env")
	// lay secret Ky
	// secretKey := os.Getenv("SECRET_KEY")
	// if secretKey == "" {
	// 	secretKey = "AreUggggGay"
	// }
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "host=localhost user=postgres password=1 dbname=auth port=5433 sslmode=disable"
	}
	vaultAddr := os.Getenv("VAULT_ADDR")
	if vaultAddr == "" {
		vaultAddr = "http://127.0.0.1:8205"
	}

	vaultToken := os.Getenv("VAULT_TOKEN")
	if vaultToken == "" {
		// log.Fatal("VAULT_TOKEN không được để trống")
		vaultToken = ""
	}

	vaultPath := os.Getenv("VAULT_PATH")
	if vaultPath == "" {
		vaultPath = "mon_deployment_secretkey"
	}

	//connect db
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal("error connect gay:", err)
	}

	svcA := anotherService.NewAnotherService(vaultAddr, vaultToken, vaultPath)
	hA := anotherAdapter.NewAnotherAdapter(svcA)
	uuid := "1b01cf34-0cdb-4e07-aac5-488ff894d050"
	secretKey, _ := svcA.GetSecretKey(uuid)

	svc := AuthService.NewAuthService(secretKey, db)
	h := AuthHandler.NewAuthHandler(svc)
	r := gin.Default()
	r.POST("/api/register", h.RegisterToken)
	r.GET("/api/fulltoken", h.GetAllToken)
	r.PUT("/api/revoketoke", h.RevokeToken)
	r.PUT("/api/revoketokenfull", h.RevokeTokenFull)
	r.PUT("/api/activetokenfull", h.ActiveTokenFull)
	r.GET("api/secretkey", hA.GetSecretKey)
	r.GET("api/secretkey2", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"secretKey": secretKey}) })
	r.GET("api/dburl", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"dburl": dbURL}) })

	r.POST("/api/secretkey", hA.GenSecretKey)
	r.Run(":8080")
}
