package shared

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

const (
	MsgStatusSent      = "sent"      // 消息已发送
	MsgStatusDelivered = "delivered" // 消息已送达
	MsgStatusRead      = "read"      // 消息已读
)
const (
	// 用户消息
	MsgTypeUser = "private"
	// // 群组消息
	MsgTypeGroup = "group"
	// Ack 消息
	MsgTypeAck = "ack"
	// 系统消息
	MsgTypeSystem = "system"
	// 邀请消息
	MsgTypeInvition = "invition"
)
const (
	// 消息类型：文本消息
	MsgContentTypeText = "text"
	// 消息类型：图片消息
	MsgContentTypeImage = "image"
)

type WsMsg struct {
	MsgID   string `json:"msg_id"`   // 消息唯一标识符
	MsgType string `json:"msg_type"` // 消息类型(用户消息/群组消息/Ack 消息/系统消息/邀请消息)
	Status  string `json:"status"`   // 消息状态 (sent/delivered/read)
	WsHeader
	WsBody
}

type WsHeader struct {
	SenderId       string   `json:"sender_id"`
	ReceiverId     []string `json:"receiver_id"`
	MsgContentType string   `json:"msg_content_type"`   // 消息内容类型 (text/image)
	Timestamp      int64    `json:"timestamp"`          // 时间戳
	GroupId        string   `json:"group_id",omitempty` // 群组ID（可选）
}
type WsBody struct {
	MsgContent string                 `json:"msg_content"`
	Extra      map[string]interface{} `json:"extra,omitempty"` // 扩展字段，用于自定义场景
}

// 生成随机 ID
func generateRandomID() string {
	bytes := make([]byte, 16) // 生成 16 字节的随机数
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

// 创建一个新的消息
func NewWsUserMessage(senderId string, receiverId []string, msgContent string) WsMsg {
	header := WsHeader{
		SenderId:       senderId,
		ReceiverId:     receiverId,
		MsgContentType: MsgContentTypeText,
		Timestamp:      time.Now().UnixNano(),
	}

	body := WsBody{
		MsgContent: msgContent,
	}

	return WsMsg{
		MsgID:    generateRandomID(),
		MsgType:  MsgTypeUser,
		Status:   MsgStatusSent,
		WsHeader: header,
		WsBody:   body,
	}
}

func NewWsGroupMessage(receiverId []string, senderId, msgContent, GroupId string) WsMsg {
	header := WsHeader{
		SenderId:       senderId,
		ReceiverId:     receiverId,
		MsgContentType: MsgContentTypeText,
		Timestamp:      time.Now().UnixNano(),
		GroupId:        GroupId,
	}

	body := WsBody{
		MsgContent: msgContent,
	}

	return WsMsg{
		MsgID:    generateRandomID(),
		MsgType:  MsgTypeGroup,
		Status:   MsgStatusSent,
		WsHeader: header,
		WsBody:   body,
	}
}

// 消息格式美化
func (msg *WsMsg) PrettyPrint() string {
	formattedJSON, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		return fmt.Sprintf("Failed to format message: %v", err)
	}
	return string(formattedJSON)
}
