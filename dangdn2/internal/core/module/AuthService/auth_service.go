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
		`INSERT INTO tokens (uuid, service, exp, iat, hash_token, shape) 
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		tokenUUID, bd.ServiceName, now+bd.Ttl, now, hashHex, "default",
	)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
