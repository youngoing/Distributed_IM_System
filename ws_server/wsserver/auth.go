package wsserver

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.StandardClaims
}

func ParseToken(tokenString string) (string, error) {
	secret := os.Getenv("TOKEN_SECRET")
	// 解析JWT令牌
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil // 使用密钥验证签名
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", fmt.Errorf("invalid token signature: %w", err)
		}
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	// 验证是否过期
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", fmt.Errorf("failed to parse claims")
	}

	// 检查token是否过期
	if claims.ExpiresAt < time.Now().Unix() {
		return "", fmt.Errorf("token has expired")
	}

	return claims.UserId, nil
}
