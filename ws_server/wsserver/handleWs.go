package wsserver

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var CHECK_INTERVAL = 10 // 每 10 秒检查一次连接
// 处理 WebSocket 连接
func (server *WebSocketServerStruct) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "token is required", http.StatusBadRequest)
		logrus.Error("token is required")
		return
	}
	runMode := os.Getenv("RUN_MODE")
	var userId string
	if runMode == "test" {
		userId = token
	} else if runMode == "dev" {
		var err error
		userId, err = ParseToken(token)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unauthorized: %v", err), http.StatusUnauthorized)
			logrus.Error("token error:", err)
			return
		}
	} else {
		var err error
		userId, err = ParseToken(token)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unauthorized: %v", err), http.StatusUnauthorized)
			logrus.Error("token error:", err)
			return
		}
	}

	ws, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Error("Failed to upgrade to WebSocket:", err)
		return
	}

	logrus.Infof("User %s connected", userId)
	defer ws.Close()
	// 为该用户创建并存储一个写入通道
	messageChannel := make(chan []byte)
	server.addUserToChannel(userId, messageChannel)
	addUserToRedis(userId, server.nodeId)
	defer server.removeUserFromChannel(userId)
	// 启动一个 goroutine 处理消息发送
	go handleWrites(ws, messageChannel)
	//处理数据库存储的消息
	go server.handleDBmsg(userId)
	//处理队列中的消息
	go server.consumeNodeMessage()
	// // 启动检查 WebSocket 连接的 goroutine
	// go server.checkConnection(ws, userId)

	// 处理 WebSocket 消息
	server.handleMessages(ws)

}

func handleWrites(ws *websocket.Conn, messageChannel chan []byte) {
	for {
		message, ok := <-messageChannel
		if !ok {
			logrus.Warn("Message channel is closed")
			return
		}
		// 发送消息给 WebSocket 客户端
		err := ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			logrus.Error("Failed to write message:", err)
			return
		}
	}
}

// 新的检查 WebSocket 连接的逻辑
func (server *WebSocketServerStruct) checkConnection(wsConn *websocket.Conn, userId string) {
	// 定时检查连接是否正常
	ticker := time.NewTicker(time.Duration(CHECK_INTERVAL) * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		// 设置读取超时时间，防止阻塞
		wsConn.SetReadDeadline(time.Now().Add(time.Duration(CHECK_INTERVAL) * time.Second))

		// 尝试读取一个小的数据包
		err := wsConn.WriteMessage(websocket.PingMessage, []byte("ping"))
		if err != nil {
			logrus.Infof("User %s WebSocket ping failed: %v, removing from Redis.", userId, err)
			// 移除 Redis 中的用户信息
			removeUserFromRedis(userId)
			server.removeUserFromChannel(userId)
			return
		}

		// 更新 Redis 中的过期时间
		updateUserExpiration(userId)
	}
}
