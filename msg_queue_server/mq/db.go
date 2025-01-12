package mq

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type handler struct {
	Channel    *amqp.Channel
	MysqlConn  *sql.DB
	rabbitConn *amqp.Connection // 保存 RabbitMQ 连接，用于后续关闭
}

func (h *handler) RunDbHandler() {
	h.ConsumeStoreMsg(StoreQueueName)
}

// 创建 Handler 并初始化 RabbitMQ 和 MySQL
func NewHandler() (*handler, error) {
	// 初始化 RabbitMQ
	rabbitConn, ch, err := createRabbitMQ()
	if err != nil {
		return nil, fmt.Errorf("failed to create RabbitMQ channel: %v", err)
	}

	// 初始化 MySQL
	db, err := createDBPool("root:qweasdzxc@unix(/var/run/mysqld/mysqld.sock)/websocket")
	if err != nil {
		// 如果数据库连接失败，需要关闭 RabbitMQ 连接
		rabbitConn.Close()
		return nil, fmt.Errorf("failed to create database pool: %v", err)
	}

	return &handler{
		Channel:    ch,
		MysqlConn:  db,
		rabbitConn: rabbitConn,
	}, nil
}

func (h *handler) ConsumeStoreMsg(queueName string) {
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
		h.handleMsg(msg.Body)
	}
}

// 创建 RabbitMQ 连接和通道
func createRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	// 连接到 RabbitMQ
	rabbitConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	// 打开通道
	ch, err := rabbitConn.Channel()
	if err != nil {
		rabbitConn.Close() // 如果通道创建失败，关闭连接
		return nil, nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	return rabbitConn, ch, nil
}

// 创建数据库连接池
func createDBPool(dsn string) (*sql.DB, error) {
	// 打开数据库连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	// 设置连接池参数
	db.SetMaxOpenConns(5)                  // 最大打开连接数
	db.SetMaxIdleConns(5)                  // 最大空闲连接数
	db.SetConnMaxLifetime(5 * time.Minute) // 连接最大生命周期

	// 测试数据库连接
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return db, nil
}

// 关闭所有连接
func (h *handler) Close() {
	if h.Channel != nil {
		h.Channel.Close()
	}
	if h.rabbitConn != nil {
		h.rabbitConn.Close()
	}
	if h.MysqlConn != nil {
		h.MysqlConn.Close()
	}
}

// 处理队列消息
func (h *handler) handleMsg(msg []byte) {
	// 解析消息
	var wsMsg WsMsg

	// 解析消息
	err := json.Unmarshal(msg, &wsMsg)
	if err != nil {
		logrus.Errorf("Failed to parse message: %v", err)
		return
	}

	// 解析成功，可以使用结构体中的数据
	logrus.Infof("Received message: %+v", wsMsg)
	switch wsMsg.MsgType {
	case MsgTypeUser:
		logrus.Infof("Received user message: %+v", wsMsg)
		h.handlePrivateMsg(wsMsg)
	case MsgTypeGroup:
		logrus.Infof("Received group message: %+v", wsMsg)
		h.handleGroupMsg(wsMsg)
	case MsgTypeInvition:
		logrus.Infof("Received invitation message: %+v", wsMsg)
		h.handlePrivateMsg(wsMsg)
	}

}

// 处理私信
func (h *handler) handlePrivateMsg(wsMsg WsMsg) {
	// 获取在线用户列表
	if wsMsg.ReceiverId == nil {
		logrus.Warn("ReceiverId is nil")
		return
	}
	if len(wsMsg.ReceiverId) == 0 {
		logrus.Warn("ReceiverId is empty")
		return
	}
	if len(wsMsg.ReceiverId) == 1 {
		// 单聊消息
		receiverId := wsMsg.ReceiverId[0]
		err := h.saveUserMsgToDB(wsMsg, receiverId)
		if err != nil {
			logrus.Errorf("Failed to save message to database: %v", err)
			return
		}

	} else {
		// 私聊多发消息
		for _, receiverId := range wsMsg.ReceiverId {
			err := h.saveUserMsgToDB(wsMsg, receiverId)
			if err != nil {
				logrus.Errorf("Failed to save message to database: %v", err)
				continue
			}
		}
	}

}

// 处理群聊
func (h *handler) handleGroupMsg(wsMsg WsMsg) {
	if wsMsg.ReceiverId == nil {
		logrus.Warn("ReceiverId is nil")
		return
	}
	if len(wsMsg.ReceiverId) == 0 {
		logrus.Warn("ReceiverId is empty")
		return
	}
	if len(wsMsg.ReceiverId) == 1 {
		// 单聊消息
		receiverId := wsMsg.ReceiverId[0]
		err := h.saveGroupMsgToDB(wsMsg, receiverId, wsMsg.GroupId)
		if err != nil {
			logrus.Errorf("Failed to save message to database: %v", err)
			return
		}
	}
	// 私聊多发消息
	for _, receiverId := range wsMsg.ReceiverId {
		err := h.saveGroupMsgToDB(wsMsg, receiverId, wsMsg.GroupId)
		if err != nil {
			logrus.Errorf("Failed to save message to database: %v", err)
			continue
		}
	}

}

func (h *handler) saveGroupMsgToDB(wsMsg WsMsg, receiveId, groupId string) (err error) {
	// 将消息序列化为 JSON 格式
	msgData, err := json.Marshal(wsMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal message to JSON: %v", err)
	}

	// 准备 SQL 语句
	stmt, err := h.MysqlConn.Prepare("INSERT INTO offline_group_message (receiver_id,group_id, msg) VALUES (?,?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	// 执行 SQL 语句
	_, err = stmt.Exec(receiveId, groupId, string(msgData))
	if err != nil {
		return fmt.Errorf("failed to execute SQL statement: %v", err)
	}

	logrus.Infof("Message saved to database: %+v", wsMsg)
	return nil
}

func (h *handler) saveUserMsgToDB(wsMsg WsMsg, receiveId string) (err error) {
	// 将消息序列化为 JSON 格式
	msgData, err := json.Marshal(wsMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal message to JSON: %v", err)
	}

	// 准备 SQL 语句
	stmt, err := h.MysqlConn.Prepare("INSERT INTO offline_private_message (receiver_id, msg) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	// 执行 SQL 语句
	_, err = stmt.Exec(receiveId, string(msgData))
	if err != nil {
		return fmt.Errorf("failed to execute SQL statement: %v", err)
	}

	logrus.Infof("Message saved to database: %+v", wsMsg)
	return nil
}
