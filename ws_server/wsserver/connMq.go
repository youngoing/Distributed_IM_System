package wsserver

import (
	"encoding/json"
	"ws_server/shared"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func (server *WebSocketServerStruct) connectRabbitMQ() error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		logrus.Error("Failed to connect to RabbitMQ:", err)
		return err
	}
	server.resourceManager.MqConn = conn
	ch, err := conn.Channel()
	if err != nil {
		logrus.Error("Failed to open RabbitMQ channel:", err)
		return err
	}
	server.resourceManager.MqChannel = ch
	logrus.Info("Connected to RabbitMQ successfully")
	return nil
}

// Forward message to RabbitMQ
func (server *WebSocketServerStruct) forwardMessageToMQ(message shared.WsMsg, dequeueName string) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		logrus.Error("Failed to serialize WsMsg to []byte:", err)
		return err
	}
	err = server.resourceManager.MqChannel.Publish(
		"",
		dequeueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageBytes,
		},
	)
	if err != nil {
		logrus.Error("Failed to forward message to MQ:", err)
		return err
	}
	logrus.Infof("Message forwarded to MQ successfully: %s", dequeueName)
	return nil
}

// Consume messages from RabbitMQ
func (server *WebSocketServerStruct) consumeNodeMessage() error {
	dequeueName := server.nodeId + ".msg"
	msgs, err := server.resourceManager.MqChannel.Consume(
		dequeueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logrus.Error("Failed to consume MQ messages:", err)
		return err
	}
	for msg := range msgs {
		var message shared.WsMsg
		err := json.Unmarshal(msg.Body, &message)
		if err != nil {
			logrus.Error("Failed to parse MQ message:", err)
			continue
		}
		switch message.MsgType {
		case shared.MsgTypeUser:
			server.handleUserMessage(message)
		case shared.MsgTypeGroup:
			server.handleGroupMessage(message)
		}
	}
	return nil
}
