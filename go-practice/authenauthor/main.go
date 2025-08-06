package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var user = map[string]string{
	"d":  "1",
	"s2": "1",
	"d3": "1",
}

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

func init() {
	var err error
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048) //random để sinh ra struct của khóa bất đối xứng độ dài 2048 bit có cả publickey bên trong
	if err != nil {
		log.Fatalf("Lỗi tạo private key: %v", err)
	}
	publicKey = &privateKey.PublicKey
}

func loginHandler(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := c.ShouldBindJSON(&loginData) // maping phần body sang struct
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"mess": "Dữ liệu đầu vào không hợp lệ"})
		return
	}

	pass, ok := user[loginData.Username]
	if !ok || pass != loginData.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"mess": "Thông tin đăng nhập sai"})
		return
	}

	// Tạo claims(payload)
	claims := jwt.MapClaims{
		"username": loginData.Username,
		"role":     "admin",
		"exp":      time.Now().Add(time.Minute * 5).Unix(), // Thời hạn token
	}

	// Ký token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(privateKey) //Mã hóa base64url phần header và payload.
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"mess": "Không tạo được token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Đăng nhập thành công",
		"token":   signedToken,
	})
}

// Middleware kiểm tra JWT
func jwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Thiếu Authorization header"})
			c.Abort()
			return
		}

		const prefix = "Bearer "
		if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Sai định dạng Bearer token"})
			c.Abort()
			return
		}

		tokenStr := authHeader[len(prefix):]

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) { //Parse token: kiểm tra chữ ký, giải mã header + payload.
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Sai phương thức ký")
			}
			return publicKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token không hợp lệ"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("claims", claims)
		}

		c.Next()
	}
}

// API yêu cầu xác thực token
func readHandler(c *gin.Context) {
	claims := c.MustGet("claims").(jwt.MapClaims)
	username := claims["username"]
	role := claims["role"]

	c.JSON(http.StatusOK, gin.H{
		"message":  "Bạn đã truy cập thành công vào tài nguyên protected",
		"username": username,
		"role":     role,
	})
}

// Có thể mở rộng nếu cần
func writeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Write thành công"})
}

func main() {
	r := gin.Default()

	r.POST("/login", loginHandler)

	// Các route cần xác thực JWT
	protected := r.Group("/")
	protected.Use(jwtMiddleware())
	protected.GET("/read", readHandler)
	protected.POST("/write", writeHandler)

	r.Run(":8080")
}
