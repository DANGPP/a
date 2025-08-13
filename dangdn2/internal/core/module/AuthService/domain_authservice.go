package AuthService

import (
	"test/internal/core/module/anotherService"

	"gorm.io/gorm"
)

type AuthService struct {
	ano *anotherService.Another
	db  *gorm.DB
}

// body gửi trong payload
type Bodi struct {
	UUID        string `json:"uuid"`
	ServiceName string `json:"serviceName"`
	Ttl         int64  `json:"ttl"`
}
