package auth_service

import (
	// "context"
	// "crypto/hmac"
	// "crypto/sha256"
	// "encoding/base64"
	// "fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	// "github.com/google/uuid"
	// "github.com/jmoiron/sqlx"
	// "github.com/redis/go-redis/v9"
)

// type Config struct {
// 	PostgresDSN   string
// 	RedisAddr     string
// 	JWTSigningKey string // base64
// 	HashSecret    string // base64
// 	KeyID         string
// }

// type TokenService struct {
// 	// cfg  Config
// 	// db   *sqlx.DB
// 	// rdb  *redis.Client
// 	keys map[string][]byte // kid -> secret
// }

// // func NewTokenService(cfg Config, db *sqlx.DB, rdb *redis.Client) (*TokenService, error) {
// func NewTokenService(jti, service string, exp string) (*TokenService, error) {
// 	rawKey, err := base64.StdEncoding.DecodeString(cfg.JWTSigningKey)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid base64 JWTSigningKey: %w", err)
// 	}
// 	return &TokenService{
// 		cfg:  cfg,
// 		db:   db,
// 		rdb:  rdb,
// 		keys: map[string][]byte{cfg.KeyID: rawKey},
// 	}, nil
// }

// type IssueTokenResult struct {
// 	Token     string
// 	JTI       string
// 	KID       string
// 	ExpiresAt time.Time
// }

// func (s *TokenService) IssueToken(ctx context.Context, service string, ttl time.Duration) (*IssueTokenResult, error) {
// 	jti := uuid.NewString()
// 	now := time.Now().UTC()
// 	expiresAt := now.Add(ttl)

// 	secret, ok := s.keys[s.cfg.KeyID]
// 	if !ok {
// 		return nil, fmt.Errorf("no signing key for kid=%s", s.cfg.KeyID)
// 	}

// 	claims := jwt.MapClaims{
// 		"jti":     jti,
// 		"service": service,
// 		"iat":     now.Unix(),
// 		"exp":     expiresAt.Unix(),
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	token.Header["kid"] = s.cfg.KeyID

// 	tokenStr, err := token.SignedString(secret)
// 	if err != nil {
// 		return nil, err
// 	}

// 	tokenHash, err := s.hashToken(tokenStr)
// 	if err != nil {
// 		return nil, err
// 	}

// 	_, err = s.db.ExecContext(ctx, `
// 		INSERT INTO tokens (jti, kid, token_hash, service, issued_at, expires_at, status)
// 		VALUES ($1,$2,$3,$4,$5,$6,'active')
// 	`, jti, s.cfg.KeyID, tokenHash, service, now, expiresAt)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ttlCache := time.Until(expiresAt)
// 	if err := s.rdb.Set(ctx, "active:"+jti, "1", ttlCache).Err(); err != nil {
// 		// log warning only
// 		fmt.Printf("warn: redis set failed: %v\n", err)
// 	}

// 	return &IssueTokenResult{
// 		Token:     tokenStr,
// 		JTI:       jti,
// 		KID:       s.cfg.KeyID,
// 		ExpiresAt: expiresAt,
// 	}, nil
// }

//	func (s *TokenService) hashToken(token string) (string, error) {
//		secretB, err := base64.StdEncoding.DecodeString(s.cfg.HashSecret)
//		if err != nil {
//			return "", fmt.Errorf("invalid HashSecret: %w", err)
//		}
//		mac := hmac.New(sha256.New, secretB)
//		mac.Write([]byte(token))
//		sum := mac.Sum(nil)
//		return base64.StdEncoding.EncodeToString(sum), nil
//	}
type TokenService struct {
	jwtSecret []byte
}

func NewTokenService(secret string) *TokenService {
	return &TokenService{jwtSecret: []byte(secret)}
}

func (t *TokenService) GenToken(service string) (string, error) {
	claims := jwt.MapClaims{
		"service": service,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(t.jwtSecret)
}

// func (t *TokenService) ValToken() (string, error){

// }
