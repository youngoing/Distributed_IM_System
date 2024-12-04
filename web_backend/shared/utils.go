package shared

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

type Claims struct {
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.StandardClaims
}

func GenerateToken(userId, userName string) (string, error) {
	// 获取USER_EXPIRATION_HOURS环境变量
	expirationHoursStr := os.Getenv("USER_EXPIRATION_HOURS")

	// 如果获取不到或解析失败，则使用默认值24
	if expirationHoursStr == "" {
		expirationHoursStr = "24" // 默认值为24小时
	}

	// 尝试将字符串转换为整数
	expirationHours, err := strconv.Atoi(expirationHoursStr)
	if err != nil {
		// 如果转换失败，记录日志并设置默认值为24小时
		log.Println("Invalid USER_EXPIRATION_HOURS, using default value of 24 hours.")
		expirationHours = 24
	}

	// 设置过期时间
	duration := time.Duration(expirationHours) * time.Hour

	// 创建JWT声明
	claims := Claims{
		UserId:   userId,
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			// 设置过期时间
			ExpiresAt: time.Now().Add(duration).Unix(),
			// JWT标准中的一个声明，用于标识令牌的发行者
			Issuer: "youngo_backend",
		},
	}

	// 生成访问令牌
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return token, nil
}

// func ParseToken(tokenString string) (*Claims, error) {
// 	// 解析JWT令牌
// 	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
// 		return []byte("secret"), nil // 使用密钥验证签名
// 	})
// 	if err != nil {
// 		if err == jwt.ErrSignatureInvalid {
// 			return nil, fmt.Errorf("invalid token signature: %w", err)
// 		}
// 		return nil, fmt.Errorf("failed to parse token: %w", err)
// 	}

// 	// 验证是否过期
// 	claims, ok := token.Claims.(*Claims)
// 	if !ok {
// 		return nil, fmt.Errorf("failed to parse claims")
// 	}

// 	// 检查token是否过期
// 	if claims.ExpiresAt < time.Now().Unix() {
// 		return nil, fmt.Errorf("token has expired")
// 	}

// 	return claims, nil
// }

// confirmEnv 检查多个环境变量是否存在，并返回错误信息
func ConfirmEnv() error {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// 需要检查的环境变量列表
	requiredEnvs := []string{
		"MYSQL_DATABASE_URL",
		"REDIS_URL",
		"USER_EXPIRATION_HOURS",
	}

	// 遍历检查每个环境变量
	for _, env := range requiredEnvs {
		if value := os.Getenv(env); value == "" {
			return fmt.Errorf("%s is not set", env) // 返回更详细的错误信息
		}
	}

	return nil
}
