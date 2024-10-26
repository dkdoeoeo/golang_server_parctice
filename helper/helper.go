package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// GetSha256
// 生成 sha256
func GetSha256(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

var myKey = []byte("post-paltform")

// GenerateToken
// 生成 token
func GenerateAccessToken(email string) (string, error) {
	hashedEmail := sha256.Sum256([]byte(email))
	return hex.EncodeToString(hashedEmail[:]), nil
}

// AnalyseToken
// 解析 token
func AnalyseToken(tokenString string) (*UserClaims, error) {
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("analyse Token Error:%v", err)
	}
	return userClaim, nil
}
