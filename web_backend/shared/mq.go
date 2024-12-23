package shared

import (
	"os"

	"github.com/streadway/amqp"
)

var (
	rabbitConn    *amqp.Connection
	rabbitChannel *amqp.Channel
)

func InitRabbitMQ() error {
	var err error
	rabbitUrl := os.Getenv("RABBITMQ_URL")
	rabbitConn, err = amqp.Dial(os.Getenv(rabbitUrl))
	if err != nil {
		return err
	}
	rabbitChannel, err = rabbitConn.Channel()
	if err != nil {
		return err
	}
	return nil
}
func CloseRabbitMQ() {
	if rabbitChannel != nil {
		rabbitChannel.Close()
	}
	if rabbitConn != nil {
		rabbitConn.Close()
	}
}
