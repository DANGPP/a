package handler

import (
	"dangdn2-test-go/internal/core/domain"
	"dangdn2-test-go/internal/core/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service *service.UserService
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.Service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy danh sách người dùng"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) Register(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	err := h.Service.Register(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể đăng ký"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đăng ký thành công"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var input struct {
		Name string `json:"name"`
		Pass string `json:"pass"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu name hoặc pass"})
		return
	}

	user, err := h.Service.Login(input.Name, input.Pass)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sai tài khoản hoặc mật khẩu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đăng nhập thành công", "user": user})
}
