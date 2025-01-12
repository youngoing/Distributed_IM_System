package shared

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

var (
	RabbitConn *amqp091.Connection
	queueName  = "input.msg"
)

// 初始化 RabbitMQ 连接
func InitRabbitMQ() error {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	queueName = "input.msg"

	if rabbitURL == "" || queueName == "" {
		return fmt.Errorf("missing required environment variables RABBITMQ_URL or input.msg")
	}

	conn, err := amqp091.Dial(rabbitURL)
	if err != nil {
		logrus.WithError(err).Error("Failed to connect to RabbitMQ")
		return err
	}

	RabbitConn = conn
	logrus.Info("Connected to RabbitMQ successfully")
	return nil
}

// 关闭 RabbitMQ 连接
func CloseRabbitMQ() {
	if RabbitConn != nil {
		if err := RabbitConn.Close(); err != nil {
			logrus.WithError(err).Warn("Failed to close RabbitMQ connection")
		} else {
			logrus.Info("RabbitMQ connection closed successfully")
		}
	}
}

// 声明队列（可复用）
func declareQueue(channel *amqp091.Channel) error {
	_, err := channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // autoDelete
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		logrus.WithError(err).Error("Failed to declare RabbitMQ queue")
	}
	return err
}

// 发送消息到队列
func ForwardMessageToMQ(message WsMsg) error {
	if RabbitConn == nil {
		return fmt.Errorf("RabbitMQ connection is not initialized")
	}

	// 创建通道
	channel, err := RabbitConn.Channel()
	if err != nil {
		logrus.WithError(err).Error("Failed to create RabbitMQ channel")
		return err
	}
	defer func() {
		if err := channel.Close(); err != nil {
			logrus.WithError(err).Warn("Failed to close RabbitMQ channel")
		}
	}()

	// 声明队列
	if err := declareQueue(channel); err != nil {
		return err
	}

	// 序列化消息
	messageBytes, err := json.Marshal(message)
	if err != nil {
		logrus.WithError(err).Error("Failed to serialize WsMsg")
		return err
	}

	// 发布消息
	err = channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        messageBytes,
		},
	)
	if err != nil {
		logrus.WithError(err).Error("Failed to forward message to MQ")
		return err
	}

	logrus.Infof("Message forwarded to MQ successfully: %s", queueName)
	return nil
}
