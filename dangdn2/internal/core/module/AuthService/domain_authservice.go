package AuthService

import (
	"gorm.io/gorm"
)

type AuthService struct {
	secretKey string
	db        *gorm.DB
}

// body gửi trong payload
type Bodi struct {
	ServiceName string `json:"serviceName"`
	Ttl         int64  `json:"ttl"`
}
