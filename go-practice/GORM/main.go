package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Token struct {
	UUID      string `gorm:"primaryKey;type:uuid"`
	Service   string
	Exp       int64
	Iat       int64
	HashToken string
	Status    string
}

func main() {
	dsn := "host=localhost user=postgres password=1 dbname=auth port=5433 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Lỗi kết nối DB:", err)
	}
	// Query tất cả token
	var tokens []Token
	if err := db.Find(&tokens).Error; err != nil {
		log.Fatal("Lỗi query tokens:", err)
	}

	// In kết quả
	for _, t := range tokens {
		fmt.Printf("UUID: %s, Service: %s, Status: %s\n", t.UUID, t.Service, t.Status)
	}

}
