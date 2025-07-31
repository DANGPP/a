package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Struct ánh xạ với bảng users
type User struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	Age  int    `json:"age"`
	Pass string `json:"pass"`
}

var db *gorm.DB

func main() {
	// Kết nối tới PostgreSQL
	dsn := "host=localhost user=postgres password = 1 dbname=testdb port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Lỗi kết nối DB: ", err)
	}

	// Tạo bảng nếu chưa có
	db.AutoMigrate(&User{})

	// // Thêm bản ghi mẫu nếu cần
	// db.Create(&User{Name: "dang1", Age: 21})
	// db.Create(&User{Name: "dang2", Age: 22})

	// Gin setup
	r := gin.Default()
	r.GET("/users", GetUsers)
	r.POST("/register", Register)
	r.POST("/login", Login)
	r.Run(":8080")
}

// Lấy danh sách người dùng từ DB
func GetUsers(c *gin.Context) {
	var users []User
	db.Find(&users)
	c.JSON(http.StatusOK, users)
}

// Đăng ký người dùng mới
func Register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	// Lưu vào DB
	result := db.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đăng ký thành công"})
}

// Đăng nhập người dùng
func Login(c *gin.Context) {
	var input struct {
		Name string `json:"name"`
		Pass string `json:"pass"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu Name hoặc Pass"})
		return
	}

	var user User
	result := db.Where("Name = ? AND Pass = ?", input.Name, input.Pass).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sai tên đăng nhập hoặc mật khẩu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đăng nhập thành công", "user": user})
}
