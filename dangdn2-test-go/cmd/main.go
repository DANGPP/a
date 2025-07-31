package main

import (
	"dangdn2-test-go/internal/adapter/handler"
	"dangdn2-test-go/internal/adapter/repo"
	"dangdn2-test-go/internal/core/domain"
	"dangdn2-test-go/internal/core/service"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=1 dbname=testdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Lỗi kết nối DB:", err)
	}

	// Tạo bảng nếu chưa có
	db.AutoMigrate(&domain.User{}) // hoặc migrate bằng domain.User nếu cần

	// Setup repository và service
	userRepo := &repo.GormUserRepo{DB: db}
	userService := &service.UserService{Repo: userRepo}
	userHandler := &handler.UserHandler{Service: userService}

	// Setup routes
	r := gin.Default()
	r.GET("/users", userHandler.GetUsers)
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	r.Run(":8080")
}
