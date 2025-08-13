package AuthService

import (
	"crypto/sha256"
	"encoding/hex"
)

func hastoken(tokenstring string) string {
	hash := sha256.Sum256([]byte(tokenstring))
	return hex.EncodeToString(hash[:])
}
