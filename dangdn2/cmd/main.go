package main

import (
	// "context"
	"fmt"
	// "log"
	"os"
	"time"

	"test/internal/adapter/module/AuthHandler"
	"test/internal/adapter/module/anotherAdapter"
	"test/internal/core/module/AuthService"
	"test/internal/core/module/anotherService"

	"test/internal/infra/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

func LoggerMiddle(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		log.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"start":  start,
		}).Info("Bắt đầu gửi")
		c.Next()
		latency := time.Since(start)
		log.WithFields(logrus.Fields{
			"status":   c.Writer.Status(),
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"latency":  latency,
			"clientIP": c.ClientIP(),
		}).Info("kết thúc gửi")

	}
}

func main() {
	log := logger.Init(logrus.DebugLevel)
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

	svcA, _ := anotherService.NewAnotherService(vaultAddr, vaultToken, vaultPath) // khởi tạo service keyvault
	hA := anotherAdapter.NewAnotherAdapter(svcA)

	svc := AuthService.NewAuthService(svcA, db) //khởi tạo service của module authenauthor
	h := AuthHandler.NewAuthHandler(svc)

	r := gin.New()
	r.Use(LoggerMiddle(log), gin.Recovery())
	//1, Tạo secret Key
	r.POST("/api/secretkey", hA.GenSecretKey)

	//2 Tạo token với secret key
	r.POST("/api/register", h.RegisterToken)

	//3 Xem toàn bộ token
	r.GET("/api/fulltoken", h.GetAllToken)

	//4 thu hồi token
	r.PUT("/api/revoketoke", h.RevokeToken)

	//5 thu hồi toàn bộ token
	r.PUT("/api/revoketokenfull", h.RevokeTokenFull)

	//6 Active toàn bộ token
	r.PUT("/api/activetokenfull", h.ActiveTokenFull)

	//7 xem uuid secretkey
	r.GET("/api/uuidsecretkey", hA.GetFullUUIDkey)

	// back door
	r.GET("api/secretkey", hA.GetSecretKey)
	r.GET("api/secretkey2", h.GetSecretKey)
	r.GET("api/dburl", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"dburl": dbURL}) })

	r.Run(":8080")
}
