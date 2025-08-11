package AuthService

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthService struct {
	secretKey string
	db        *pgxpool.Pool
}

type Bodi struct {
	ServiceName string `json:"serviceName"`
	Ttl         int64  `json:"ttl"`
}

func NewAuthService(secretKey string, db *pgxpool.Pool) *AuthService {
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
		"exp":     now + bd.Ttl, // Token hết hạn sau 15 phút
		"iat":     now,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(a.secretKey))

	// Hash token
	hash := sha256.Sum256([]byte(tokenString))
	hashHex := hex.EncodeToString(hash[:])

	// Save to DB
	_, err := a.db.Exec(ctx,
		`INSERT INTO tokens (uuid, service, exp, iat, hash_token, status) 
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		tokenUUID, bd.ServiceName, now+bd.Ttl, now, hashHex, "active",
	)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// 2. thu hồi token
func (a *AuthService) RevokeToken(uuid string, ctx context.Context) (string, error) {
	_, err := a.db.Exec(ctx,
		`update tokens set status = $1 where uuid = $2 `, "revoke", uuid,
	)
	if err != nil {
		return "lỗi ở auth_service.go", err
	}
	return "revoke", nil
}

// 3. thu hồi token full
func (a *AuthService) RevokeTokenFull(ctx context.Context) (string, error) {
	_, err := a.db.Exec(ctx,
		`update tokens set status = $1 `, "revoke",
	)
	if err != nil {
		return "lỗi ở auth_service.go", err
	}
	return "revokeAll", nil
}

// 3. active token
func (a *AuthService) ActiveTokenFull(ctx context.Context) (string, error) {
	_, err := a.db.Exec(ctx,
		`update tokens set status = $1 `, "active",
	)
	if err != nil {
		return "lỗi ở auth_service.go", err
	}
	return "ActiveAll", nil
}
