package wsserver

import (
	"encoding/json"
	"ws_server/shared"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func (server *WebSocketServerStruct) handleMessages(ws *websocket.Conn) {
	for {
		_, msg, err := ws.ReadMessage()
		// 错误处理
		if err != nil {
			// 检查 WebSocket 是否正常关闭
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Info("WebSocket is closed") // WebSocket连接被关闭，记录并退出
			} else {
				// 其他错误处理
				logrus.Errorf("Failed to read message: %v, closing connection", err) // 读取消息失败的错误日志
			}
			break // 错误发生时退出循环
		}
		// 将 JSON 消息解析为 WsMessage
		var wsMessage shared.WsMsg
		err = json.Unmarshal(msg, &wsMessage)
		if err != nil {
			logrus.Warn("parse json message error:", err)
			continue
		}
		switch wsMessage.MsgType {

		case shared.MsgTypeUser:
			server.handleUserMessage(wsMessage)
		case shared.MsgTypeGroup:
			server.handleGroupMessage(wsMessage)
		case shared.MsgTypeInvition:
			server.handleInvitionMessage(wsMessage)
		}
	}
}

func (server *WebSocketServerStruct) handleUserMessage(wsMessage shared.WsMsg) {
	server.userManager.mu.RLock() // 使用读锁
	logrus.Info("channels: ", server.userManager.userChannels)
	server.userManager.mu.RUnlock() // 释放读锁
	// 输出解析后的消息
	logrus.Infof("received private message: %+v", wsMessage)
	// 记录需要移除的接收者 ID
	var remainingReceivers []string
	for _, receiverId := range wsMessage.ReceiverId {
		server.userManager.mu.RLock()                        // 使用读锁
		_, exists := server.userManager.userList[receiverId] // 返回的是 string
		server.userManager.mu.RUnlock()                      // 释放读锁

		if exists {
			// 如果接收者在当前节点，直接发送消息
			logrus.Infof("receiver %s is in current node,forwarding msg", receiverId)
			server.sendMessageToUser(receiverId, wsMessage)
		} else {
			// 如果接收者不在当前节点，记录下来
			logrus.Infof("receiver %s is not in current node,forwarding msg to queue", receiverId)
			remainingReceivers = append(remainingReceivers, receiverId)
		}
	}
	if len(remainingReceivers) > 0 {
		// 更新消息中的接收者 ID
		wsMessage.ReceiverId = remainingReceivers
		server.forwardMessageToMQ(wsMessage, "input.msg")
	}
}

func (server *WebSocketServerStruct) handleGroupMessage(wsMessage shared.WsMsg) {
	// 输出解析后的消息
	logrus.Infof("received group message: %+v", wsMessage)
	// 记录需要移除的接收者 ID
	var remainingReceivers []string
	for _, receiverId := range wsMessage.ReceiverId {
		server.userManager.mu.RLock()
		_, exists := server.userManager.userList[receiverId]
		server.userManager.mu.RUnlock()
		if exists {
			// 如果接收者在当前节点，直接发送消息
			logrus.Infof("receiver %s is in current node,forwarding msg", receiverId)
			server.sendMessageToUser(receiverId, wsMessage)
		} else {
			// 如果接收者不在当前节点，记录下来
			logrus.Infof("receiver %s is not in current node,forwarding msg to queue", receiverId)
			remainingReceivers = append(remainingReceivers, receiverId)
		}
	}

	// 更新消息中的接收者 ID
	wsMessage.ReceiverId = remainingReceivers

	// 如果消息中还有接收者 ID，说明接收者不在当前节点，转发消息到
	if len(wsMessage.ReceiverId) > 0 {
		// TODO: 转发到 RabbitMQ
		logrus.Info("接收者不在当前节点，消息将被转发")
		server.forwardMessageToMQ(wsMessage, "input.msg")
	}
}

func (server *WebSocketServerStruct) handleInvitionMessage(wsMessage shared.WsMsg) {
	logrus.Infof("received intition message: %+v", wsMessage)
	var remainingReceivers []string
	for _, receiverId := range wsMessage.ReceiverId {
		server.userManager.mu.RLock()
		_, exists := server.userManager.userList[receiverId]
		server.userManager.mu.RUnlock()
		if exists {
			logrus.Infof("receiver %s is in current node,forwarding msg", receiverId)
			server.sendMessageToUser(receiverId, wsMessage)
		} else {
			logrus.Infof("receiver %s is not in current node,forwarding msg to queue", receiverId)
			remainingReceivers = append(remainingReceivers, receiverId)
		}
	}
	if len(remainingReceivers) > 0 {
		wsMessage.ReceiverId = remainingReceivers
		server.forwardMessageToMQ(wsMessage, "input.msg")
	}
}

func (server *WebSocketServerStruct) sendMessageToUser(userId string, message shared.WsMsg) {
	server.userManager.mu.RLock() // 使用读锁
	logrus.Info("channels: ", server.userManager.userChannels)
	for k, v := range server.userManager.userChannels {
		logrus.Info("k: ", k, "v: ", v)
	}
	channel, exists := server.userManager.userChannels[userId]
	logrus.Info("channel: ", channel)
	server.userManager.mu.RUnlock() // 释放读锁
	// 获取用户通道
	if exists && channel != nil {
		messageBytes, err := json.Marshal(message)
		if err != nil {
			logrus.Errorf("Failed to serialize WsMsg to []byte: %v", err)
			return
		}
		// 直接发送消息给通道
		channel <- messageBytes
		logrus.Infof("Message sent to user %s", userId)
	} else {
		// 用户不在，放入消息队列
		logrus.Infof("User %s is not connected, message sent to MQ", userId)
		server.forwardMessageToMQ(message, "input.msg")
	}
}
