package AuthService

import (
	"context"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	// "gorm.io/driver/postgres"
	"github.com/google/uuid"
	// "github.com/jackc/pgx/v5/pgxpool"
)

func NewAuthService(secretKey string, db *gorm.DB) *AuthService {
	return &AuthService{secretKey: secretKey,
		db: db,
	}
}


// 1. sinh token
func (a *AuthService) CreateToken(bd Bodi, ctx context.Context) (string, error) {
	now := time.Now().Unix()
	tokenUUID := uuid.New().String()

	claims := jwt.MapClaims{
		"uuid":    tokenUUID,
		"service": bd.ServiceName,
		"exp":     now + bd.Ttl,
		"iat":     now,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(a.secretKey))

	// Hash token
	// hash := sha256.Sum256([]byte(tokenString))
	// hashHex := hex.EncodeToString(hash[:])

	hashHex := hastoken(tokenString)
	// Save to DB
	newToken := Token{
		UUID:      tokenUUID,
		Service:   bd.ServiceName,
		Exp:       now + bd.Ttl,
		Iat:       now,
		HashToken: hashHex,
		Status:    "active",
	}

	if err := a.db.WithContext(ctx).Create(&newToken).Error; err != nil {
		return "", err
	}

	return tokenString, nil
}

// 2. thu hồi token
func (a *AuthService) RevokeToken(uuid string, ctx context.Context) (string, error) {
	if err := a.db.WithContext(ctx).
		Model(&Token{}).
		Where("uuid = ?", uuid).
		Update("status", "revoke").Error; err != nil {
		return "lỗi ở auth_service.go", err
	}
	return "revoke", nil
}

// 3. thu hồi token full
func (a *AuthService) RevokeTokenFull(ctx context.Context) (string, error) {
	if err := a.db.WithContext(ctx).
		Model(&Token{}).
		Update("status", "revoke").Error; err != nil {
		return "lỗi ở auth_service.go", err
	}
	return "revokeAll", nil
}

// 3. active token
func (a *AuthService) ActiveTokenFull(ctx context.Context) (string, error) {
	if err := a.db.WithContext(ctx).
		Model(&Token{}).
		Update("status", "active").Error; err != nil {
		return "lỗi ở auth_service.go", err
	}
	return "ActiveAll", nil
}

// 4. lấy full token
func (a *AuthService) GetAllToken(ctx context.Context) ([]Token, error) {
	var tokens []Token
	if err := a.db.WithContext(ctx).Find(&tokens).Error; err != nil {
		return nil, err
	}
	return tokens, nil
}
