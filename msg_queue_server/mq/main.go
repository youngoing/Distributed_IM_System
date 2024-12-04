package mq

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// func main() {
// 	InitMq()
// }

func RunMq() {
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

	// 初始化 MessageHandler
	handler := &MessageHandler{
		Channel:     ch,
		RedisClient: redisClient,
	}

	// 启动消息消费
	go handler.ConsumeInputMsg(InputQueueName)
	select {}
}
