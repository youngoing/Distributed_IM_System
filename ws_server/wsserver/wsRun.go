package wsserver

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

// 启动 WebSocket 服务器
func (server *WebSocketServerStruct) Run() {
	redisUrl := os.Getenv("REDIS_URL")
	mysqlUrl := os.Getenv("MYSQL_DATABASE_URL")

	if err := connectRedis(redisUrl); err != nil {
		logrus.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer RedisConn.Close()

	if err := server.createDBPool(mysqlUrl); err != nil {
		logrus.Fatalf("Failed to create DB pool: %v", err)
	}
	defer DbConn.Close()

	if err := server.connectRabbitMQ(); err != nil {
		logrus.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer server.resourceManager.MqConn.Close()

	go server.RunWsHttp()
	logrus.Infof("%s is running on %s:%s%s\n", server.nodeId, server.host, server.webSocketPort, server.WebSocketPath)
	value := fmt.Sprintf("%s:%s%s", server.host, server.webSocketPort, server.WebSocketPath)
	if err := appendNodeToRedis(server.nodeId, value); err != nil {
		logrus.Fatalf("Failed to append node to Redis: %v", err)
	}

	// 启动 HTTP 服务器
	go server.runHttpServer(server.host, server.httpPort)

	// 捕获终止信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	// 等待终止信号
	sig := <-sigChan
	logrus.Infof("Received: %s, going to close server", sig)
	// 清理操作
	server.cleanup()
}

func (server *WebSocketServerStruct) cleanup() {
	server.userManager.mu.Lock()
	defer server.userManager.mu.Unlock()

	// 清理 Redis 中的用户数据
	for userId := range server.userManager.userList {
		removeUserFromRedis(userId)
	}
	// 清理 Redis 中的节点数据
	removeNodeFromRedis(server.nodeId)
	// 使用 make 函数创建新的空映射
	server.userManager.userList = make(map[string]string)
	server.userManager.userChannels = make(map[string]chan []byte)

	logrus.Infof("Cleanup completed")
}

func (server *WebSocketServerStruct) RunWsHttp() {
	http.HandleFunc(server.WebSocketPath, server.handleWebSocket)
	addr := fmt.Sprintf("%s:%s", server.host, server.webSocketPort)
	if err := http.ListenAndServe(addr, nil); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}

func StartWsServer() {
	// 创建 WebSocket 服务器
	host := os.Getenv("HOST")
	webSocketPort := os.Getenv("WS_PORT")
	httpPort := os.Getenv("HTTP_PORT")
	webSocketPath := os.Getenv("WS_PATH")
	nodeId := os.Getenv("NODE_ID")
	server := NewWebSocketServer(host, webSocketPort, httpPort, webSocketPath, nodeId)
	server.Run()
}
