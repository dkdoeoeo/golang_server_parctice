package helper

import (
	"crypto/sha256"
	"encoding/hex"
)

var myKey = []byte("post-paltform")

// GetSha256
// 生成 sha256
func GenerateAccessToken(email string) (string, error) {
	hashedEmail := sha256.Sum256([]byte(email))
	return hex.EncodeToString(hashedEmail[:]), nil
}
