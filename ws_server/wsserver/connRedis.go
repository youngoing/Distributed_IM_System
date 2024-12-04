package wsserver

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// 定义了发送 redis存储用户的过期时间, 单位为秒
var REDIS_USER_KV_INTERVAL = 15
var ctx = context.Background()

func connectRedis(redisURI string) error {
	// 解析 Redis URI
	parsedURI, err := url.Parse(redisURI)
	if err != nil {
		logrus.Error("Failed to parse Redis URI:", err)
		return err
	}

	// 提取主机和端口
	addr := parsedURI.Host

	// 提取密码（如果有）
	var password string
	if parsedURI.User != nil {
		password, _ = parsedURI.User.Password()
	}

	// 提取数据库编号（如果有）
	dbChoice := 0                // 默认数据库编号为 0
	if len(parsedURI.Path) > 1 { // Path 格式为 "/0"
		db, err := strconv.Atoi(parsedURI.Path[1:])
		if err != nil {
			logrus.Error("Invalid database number in URI:", err)
			return err
		}
		dbChoice = db
	}

	// 创建 Redis 客户端
	RedisConn = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       dbChoice,
		PoolSize: 5, // 设置连接池大小
	})

	// 测试连接
	ctx := context.Background()
	_, err = RedisConn.Ping(ctx).Result()
	if err != nil {
		logrus.Error("Failed to connect to Redis:", err)
		return err
	}
	logrus.Info("Connected to Redis successfully")
	return nil
}

func addUserToRedis(userId, nodeId string) error {
	err := RedisConn.Set(ctx, userId, nodeId, 0).Err() // 设置过期时间为 0，表示没有过期时间
	if err != nil {
		logrus.Error("Failed to add user to Redis:", err)
		return err
	}
	return nil
}

func addUserToRedisWithExpiration(userId, nodeId string) error {
	err := RedisConn.SetEX(ctx, userId, nodeId, time.Duration(REDIS_USER_KV_INTERVAL)*time.Second).Err()
	if err != nil {
		logrus.Error("Failed to add user to Redis:", err)
		return err
	}
	return nil
}

func removeUserFromRedis(userId string) error {
	err := RedisConn.Del(ctx, userId).Err()
	if err != nil {
		logrus.Error("Failed to remove user from Redis:", err)
		return err
	}
	return nil
}

func updateUserExpiration(userId string) error {
	// 设置过期时间
	cmd := RedisConn.Expire(ctx, userId, time.Duration(REDIS_USER_KV_INTERVAL)*time.Second)

	// 判断是否成功
	if err := cmd.Err(); err != nil {
		logrus.Error("Failed to update user expiration in Redis:", err)
		return err
	}
	logrus.Infof("User %s expiration updated successfully.", userId)
	return nil
}

func appendNodeToRedis(nodeId, value string) error {
	err := RedisConn.HSet(ctx, "nodes", nodeId, value).Err()
	if err != nil {
		logrus.Error("Failed to add node to Redis:", err)
		return err
	}
	logrus.Infof("Node %s has been added to Redis", nodeId)
	return nil
}

func removeNodeFromRedis(nodeId string) error {
	err := RedisConn.HDel(ctx, "nodes", nodeId).Err()
	if err != nil {
		logrus.Error("Failed to remove node from Redis:", err)
		return err
	}
	logrus.Infof("Node %s has been removed from Redis", nodeId)
	return nil
}
