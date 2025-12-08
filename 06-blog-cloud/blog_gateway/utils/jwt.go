package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWT_SECRET 密钥
const JWT_SECRET = "moon_zhang"

type JWTClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID int64, username string) (string, error) {
	// 设置过期时间
	expirationTime := time.Now().Add(24 * time.Hour)
	// 设置secret key
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Subject:   username,
			// 设置签发人
			Issuer: "blog",
		},
	}
	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 签名令牌
	tokenString, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*JWTClaims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET), nil
	})
	if err != nil {
		return nil, err
	}
	// 验证令牌
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
