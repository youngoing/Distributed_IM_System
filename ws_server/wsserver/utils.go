package wsserver

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// confirmEnv 检查多个环境变量是否存在，并返回错误信息
func ConfirmEnv() error {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		return err
	}

	// 需要检查需要配置的环境变量列表
	requiredEnvs := []string{
		"MYSQL_DATABASE_URL",
		"RABBITMQ_URL",
		"REDIS_URL",
		"WS_PORT",
		"HTTP_PORT",
		"WS_PATH",
		"NODE_ID",
		"BACKEND_URL",
	}

	// 遍历检查每个环境变量
	for _, env := range requiredEnvs {
		if value := os.Getenv(env); value == "" {
			return fmt.Errorf("%s is not set", env) // 返回更详细的错误信息
		}
	}

	return nil
}
