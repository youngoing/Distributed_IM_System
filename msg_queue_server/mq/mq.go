package mq

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var ctx = context.Background()

// 定义常量以避免硬编码队列名称
const (
	InputQueueName = "input.msg"
	// ForwardQueueName = "forward.msg"
	StoreQueueName = "store.msg"
	Node1QueueName = "node1.msg"
	Node2QueueName = "node2.msg"
	Node3QueueName = "node3.msg"
	Node4QueueName = "node4.msg"
	Node5QueueName = "node5.msg"
)

type MessageHandler struct {
	Channel     *amqp.Channel
	RedisClient *redis.Client
}

func (h *MessageHandler) ConsumeInputMsg(queueName string) {
	msgs, err := h.Channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logrus.Fatalf("Failed to register a consumer: %v", err)
	}
	logrus.Infof("Start to consume messages from %s", queueName)

	for msg := range msgs {
		logrus.Infof("Received a message: %s", msg.Body)
		h.handleMessage(msg.Body)
	}
}

// 处理单条消息
func (h *MessageHandler) handleMessage(body []byte) {
	var wsMsg WsMsg
	err := json.Unmarshal(body, &wsMsg)
	if err != nil {
		logrus.Errorf("Failed to unmarshal message: %v", err)
		return
	}
	switch wsMsg.MsgType {
	// 用户私聊消息
	case MsgTypeUser:
		h.handlePrivateMsg(wsMsg)
	case MsgTypeGroup:
		// 群组消息
		h.handleGroupMsg(wsMsg)
	}

}

func (h *MessageHandler) handleGroupMsg(wsMsg WsMsg) {
	logrus.Infof("Received group message: %+v", wsMsg)
	// 获取在线用户列表
	if wsMsg.WsHeader.ReceiverId == nil {
		logrus.Warn("ReceiverId is nil")
		return
	}
	if len(wsMsg.WsHeader.ReceiverId) == 0 {
		logrus.Warn("ReceiverId is empty")
		return
	}
	if len(wsMsg.WsHeader.ReceiverId) == 1 {
		// 只有一个接收者，直接发送
		receiverId := wsMsg.WsHeader.ReceiverId[0]
		nodeId, err := getNodeByUser(h.RedisClient, receiverId)
		if err != nil {
			logrus.Errorf("Failed to get node for user %s: %v", receiverId, err)
			h.forwardMsgToMQ(wsMsg, StoreQueueName)
			return
		}
		// 转发消息到对应的节点
		NodeQueueName := nodeId + ".msg"
		h.forwardMsgToMQ(wsMsg, NodeQueueName)
		return
	} else {
		// 将多个接收者消息拆分为多个单聊消息
		for _, receiverId := range wsMsg.WsHeader.ReceiverId {
			msg := NewWsGroupMessage([]string{receiverId}, wsMsg.WsHeader.SenderId, wsMsg.WsBody.MsgContent, wsMsg.WsHeader.GroupId)
			nodeId, err := getNodeByUser(h.RedisClient, receiverId)
			if err != nil {
				logrus.Errorf("Failed to get node for user %s: %v", receiverId, err)
				h.forwardMsgToMQ(msg, StoreQueueName)
				return
			}

			// 转发消息到对应的节点
			NodeQueueName := nodeId + ".msg"
			h.forwardMsgToMQ(msg, NodeQueueName)
		}
	}
}

func (h *MessageHandler) handlePrivateMsg(wsMsg WsMsg) {
	logrus.Infof("Received private message: %+v", wsMsg)
	// 获取在线用户列表
	if wsMsg.WsHeader.ReceiverId == nil {
		logrus.Warn("ReceiverId is nil")
		return
	}
	if len(wsMsg.WsHeader.ReceiverId) == 0 {
		logrus.Warn("ReceiverId is empty")
		return
	}
	if len(wsMsg.WsHeader.ReceiverId) == 1 {
		// 单聊消息
		receiverId := wsMsg.WsHeader.ReceiverId[0]
		nodeId, err := getNodeByUser(h.RedisClient, receiverId)
		if err != nil {
			logrus.Errorf("Failed to get node for user %s: %v", receiverId, err)
			h.forwardMsgToMQ(wsMsg, StoreQueueName)
			return
		}
		NodeQueueName := nodeId + ".msg"
		h.forwardMsgToMQ(wsMsg, NodeQueueName)
		return
	} else {
		// 私聊多发消息
		for _, receiverId := range wsMsg.WsHeader.ReceiverId {
			nodeId, err := getNodeByUser(h.RedisClient, receiverId)
			if err != nil {
				// 如果获取节点失败，则用户不在线，保存到数据库
				logrus.Errorf("Failed to get node for user %s: %v", receiverId, err)
				h.forwardMsgToMQ(wsMsg, StoreQueueName)
				return
			}
			NodeQueueName := nodeId + ".msg"
			h.forwardMsgToMQ(wsMsg, NodeQueueName)
		}
	}
}

// 将消息转发到指定的 MQ
func (h *MessageHandler) forwardMsgToMQ(msg WsMsg, queueName string) {
	body, err := json.Marshal(msg)
	if err != nil {
		logrus.Errorf("Failed to marshal message: %v", err)
		return
	}
	logrus.Infof("queueName: %s", queueName)

	// 发布消息到指定队列
	err = h.Channel.Publish(
		"", // 默认交换机
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		logrus.Errorf("Failed to publish message to %s: %v", queueName, err)
		// 可选：在这里实现重试逻辑
	}
	logrus.Infof("Message is published to %s", queueName)
}

func createQueue(ch *amqp.Channel, queueName string, durable bool) amqp.Queue {
	q, err := ch.QueueDeclare(
		queueName, // 队列名称
		durable,   // 队列是否持久化
		false,     // 是否自动删除
		false,     // 是否排他
		false,     // 是否等待
		nil,       // 队列参数
	)
	if err != nil {
		logrus.Fatalf("Failed to declare a queue: %v", err)
	}
	return q
}

func InitMq() {
	// 连接到 RabbitMQ
	rabbitConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		logrus.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()

	// 打开通道
	ch, err := rabbitConn.Channel()
	if err != nil {
		logrus.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// 连接到 Redis
	redisClient, err := connectRedis()
	if err != nil {
		logrus.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()
	// 创建队列（如果已经存在则跳过）
	createQueue(ch, InputQueueName, true)
	createQueue(ch, StoreQueueName, true)
	createQueue(ch, Node1QueueName, true)
	createQueue(ch, Node2QueueName, true)
	createQueue(ch, Node3QueueName, true)
	createQueue(ch, Node4QueueName, true)
	createQueue(ch, Node5QueueName, true)
	logrus.Info("Queues are created successfully")
}
