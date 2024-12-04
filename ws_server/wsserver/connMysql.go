package wsserver

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"
	"ws_server/shared"

	_ "github.com/go-sql-driver/mysql" // 导入 MySQL 驱动程序
	"github.com/sirupsen/logrus"
)

type WsPrivateDBMsg struct {
	Id         int          `json:"id"` // 消息唯一标识符
	ReceiverId string       `json:"receiver_id"`
	Msg        shared.WsMsg `json:"msg"`
	CreateTime time.Time    `json:"create_time"`
}
type WsGroupDBMsg struct {
	Id         int          `json:"id"` // 消息唯一标识符
	ReceiverId string       `json:"receiver_id"`
	GroupId    string       `json:"group_id"`
	Msg        shared.WsMsg `json:"msg"`
	// CreateTime time.Time `json:"create_time"`
}

// 创建数据库连接池
func (server *WebSocketServerStruct) createDBPool(dsn string) error {
	// 打开数据库连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logrus.Errorf("failed to open database: %v", err)
		return err
	}
	// 设置连接池参数
	db.SetMaxOpenConns(5)                  // 最大打开连接数
	db.SetMaxIdleConns(5)                  // 最大空闲连接数
	db.SetConnMaxLifetime(5 * time.Minute) // 连接最大生命周期

	// 测试数据库连接
	err = db.Ping()
	if err != nil {
		logrus.Errorf("failed to connect to database: %v", err)
		return err
	}
	DbConn = db
	return nil
}

// 从数据库中获取私有消息
func takePrivateMsgsFromDB(userID string) ([]WsPrivateDBMsg, error) {
	// 执行查询语句
	rows, err := DbConn.Query("SELECT id, receiver_id, msg FROM offline_private_message WHERE receiver_id = ?", userID)
	if err != nil {
		log.Printf("Failed to query messages: %v", err)
		return nil, err
	}
	defer rows.Close()

	// 用于存储结果的切片
	var msgs []WsPrivateDBMsg

	// 解析查询结果
	for rows.Next() {
		var msg WsPrivateDBMsg
		var msgData []byte // 临时变量，用于存储从数据库读取的msg数据
		// 读取一行数据
		err := rows.Scan(&msg.Id, &msg.ReceiverId, &msgData)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
			return nil, err
		}
		// 解析字符串到 WsMsg 结构
		err = json.Unmarshal(msgData, &msg.Msg)
		if err != nil {
			log.Printf("Failed to unmarshal message data: %v", err)
			return nil, err
		}
		// 将解析的消息添加到结果切片中
		msgs = append(msgs, msg)

		// 删除已读取的消息
		_, err = DbConn.Exec("DELETE FROM offline_private_message WHERE id = ?", msg.Id)
		if err != nil {
			log.Printf("Failed to delete message: %v", err)
			return nil, err
		}
	}

	// 检查是否发生了任何扫描错误
	if err = rows.Err(); err != nil {
		log.Printf("Error occurred during rows iteration: %v", err)
		return nil, err
	}

	return msgs, nil
}

// 从数据库中获取群组消息
func takeGroupMsgsFromDB(userID string) ([]WsGroupDBMsg, error) {
	// 执行查询语句
	rows, err := DbConn.Query("SELECT id, receiver_id, group_id, msg FROM offline_group_message WHERE receiver_id = ?", userID)
	if err != nil {
		log.Printf("Failed to query messages: %v", err)
		return nil, err
	}
	defer rows.Close()

	// 用于存储结果的切片
	var msgs []WsGroupDBMsg

	// 解析查询结果
	for rows.Next() {
		var msg WsGroupDBMsg
		var msgData []byte // 临时变量，用于存储从数据库读取的msg数据
		// 读取一行数据
		err := rows.Scan(&msg.Id, &msg.ReceiverId, &msg.GroupId, &msgData)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
			return nil, err
		}
		// 解析字符串到 WsMsg 结构
		err = json.Unmarshal(msgData, &msg.Msg)
		if err != nil {
			log.Printf("Failed to unmarshal message data: %v", err)
			return nil, err
		}

		// 将解析的消息添加到结果切片中
		msgs = append(msgs, msg)

		// 删除已读取的消息
		_, err = DbConn.Exec("DELETE FROM offline_group_message WHERE id = ?", msg.Id)
		if err != nil {
			log.Printf("Failed to delete message: %v", err)
			return nil, err
		}
	}

	// 检查是否发生了任何扫描错误
	if err = rows.Err(); err != nil {
		log.Printf("Error occurred during rows iteration: %v", err)
		return nil, err
	}

	return msgs, nil
}
func (server *WebSocketServerStruct) handleDBmsg(userId string) {
	log.Println("处理数据库存储的消息")
	// 从数据库中获取用户的离线消息
	privateMsgs, err := takePrivateMsgsFromDB(userId)
	if err != nil {
		log.Printf("Failed to take private messages from database: %v", err)
		return
	}
	groupMsgs, err := takeGroupMsgsFromDB(userId)
	if err != nil {
		log.Printf("Failed to take group messages from database: %v", err)
		return
	}
	if len(privateMsgs) == 0 && len(groupMsgs) == 0 {
		log.Printf("用户 %s 没有离线消息", userId)
		return
	}

	// 将消息发送给用户
	for _, msg := range privateMsgs {
		server.sendMessageToUser(userId, msg.Msg)
	}
	for _, msg := range groupMsgs {
		server.sendMessageToUser(userId, msg.Msg)
	}
}
