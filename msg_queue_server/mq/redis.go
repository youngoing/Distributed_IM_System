package mq

import (
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

// 初始化 Redis 连接
func connectRedis() (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}
	log.Println("connect to redis success")
	return redisClient, nil
}
func getOnlineUsers(redisClient *redis.Client) (map[string]struct{}, error) {
	onlineUsers := make(map[string]struct{})

	nodes, err := redisClient.HGetAll(ctx, "nodes").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes from Redis: %v", err)
	}

	for nodeId := range nodes {
		userKey := fmt.Sprintf("user:%s", nodeId)
		users, err := redisClient.SMembers(ctx, userKey).Result()
		if err != nil {
			log.Printf("Failed to get users for node %s: %v", nodeId, err)
			continue
		}

		for _, user := range users {
			onlineUsers[user] = struct{}{}
		}
	}

	return onlineUsers, nil
}

func getOnlineNodes(redisClient *redis.Client) (map[string]string, error) {
	//从redis的nodes hash中获取所有节点
	nodes, err := redisClient.HGetAll(ctx, "nodes").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes from Redis: %v", err)
	}
	return nodes, nil
}

// 根据用户 ID 获取所在的节点
func getNodeByUser(redisClient *redis.Client, userID string) (string, error) {
	//判断用户是否在线，判断是否有key=userID,获取值为节点ID
	nodeID, err := redisClient.Get(ctx, userID).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get node for user %s: %v", userID, err)
	}
	return nodeID, nil
}

// // 获取所有在线用户及其对应的节点
// func getAllOnlineUsers(redisClient *redis.Client) (map[string]string, error) {
// 	return redisClient.HGetAll(ctx, "online_users").Result()
// }
