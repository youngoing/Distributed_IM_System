package handler

import (
	"backend/shared"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type AddFrienndRequest struct {
	SenderId   int `json:"sender_id" binding:"required"`
	ReceiverId int `json:"receiver_id" binding:"required"`
}

// 请求结构体定义
type ApplicationRequest struct {
	Action     string `json:"action" binding:"required"` // "friend" 或 "group"
	SenderId   int    `json:"sender_id" binding:"required"`
	ReceiverId int    `json:"receiver_id"`
	GroupId    int    `json:"group_id"`
}

type AuthRequest struct {
	Type       string `json:"type" binding:"required"`   // "friend" 或 "group"
	Action     string `json:"action" binding:"required"` // "accept" 或 "reject"
	SenderId   int    `json:"sender_id"`
	ReceiverId int    `json:"receiver_id"`
	GroupId    int    `json:"group_id"`
	MsgId      string `json:"msg_id" binding:"required"`
	Token      string `json:"token" binding:"required"`
}

func CreateApplication(c *gin.Context) {
	var req ApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}
	switch req.Action {
	case "friend":
		// 创建好友申请
		// 查询发送者和接收者信息
		type UserInvitationData struct {
			SenderNickname string         `json:"nickname" gorm:"column:nickname"`
			SenderAvatar   sql.NullString `json:"avatar_url" gorm:"column:avatar_url"`
		}
		var data UserInvitationData
		err := shared.MysqlDb.Table("user_details").
			Select("nickname, avatar_url").
			Where("id = ?", req.SenderId).
			First(&data).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Failed to retrieve invitation data", "details": "record not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve invitation data", "details": err.Error()})
			}
			return
		}
		// 将 SenderAvatar 转换为字符串
		senderAvatar := ""
		if data.SenderAvatar.Valid {
			senderAvatar = data.SenderAvatar.String
		}

		// 创建消息
		receivers := []string{strconv.Itoa(req.ReceiverId)}
		msgCont := data.SenderNickname + " 请求添加你为好友"
		token := uuid.New()
		tokenStr := token.String()
		msg := shared.NewWsUserInvitionMessage(
			receivers,
			strconv.Itoa(req.SenderId),
			msgCont,
			data.SenderNickname,
			senderAvatar,
			tokenStr,
		)

		// 发送消息到 MQ
		if err := shared.ForwardMessageToMQ(msg); err != nil {
			logrus.WithError(err).Error("Failed to send friend invitation")
			c.JSON(500, gin.H{
				"error":   "Failed to send friend invitation",
				"details": err.Error(),
			})
			return
		}

		// 存储到 Redis
		if err := shared.StoreInvitionToken(msg.MsgID, tokenStr); err != nil {
			logrus.WithError(err).Error("Failed to store invitation token")
			c.JSON(500, gin.H{
				"error":   "Failed to store invitation token",
				"details": err.Error(),
			})
			return
		}

		return

	case "group":
		// 创建群组申请
		// 查询发送者和群组信息
		type GroupInfo struct {
			Name         string `json:"name"`
			AvatarURL    string `json:"avatar_url"`
			UserDetailID int    `json:"user_detail_id"`
		}

		var groupInfo GroupInfo
		if err := shared.MysqlDb.Table("chat_groups").
			Select("name, avatar_url, user_detail_id").
			Where("id = ?", req.GroupId).
			First(&groupInfo).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve group data", "details": err.Error()})
			return
		}
		logrus.Info(groupInfo)

		// 然后查询发送者信息
		type SenderInfo struct {
			Nickname  string `json:"nickname"`
			AvatarURL string `json:"avatar_url"`
		}

		var senderInfo SenderInfo
		if err := shared.MysqlDb.Table("user_details").
			Select("nickname, avatar_url").
			Where("id = ?", req.SenderId).
			First(&senderInfo).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve sender data", "details": err.Error()})
			return
		}

		logrus.Info(senderInfo)

		// 创建消息
		receivers := []string{strconv.Itoa(groupInfo.UserDetailID)}
		msgCont := senderInfo.Nickname + " 申请加入群组：" + groupInfo.Name
		token := uuid.New()
		tokenStr := token.String()
		msg := shared.NewWsGroupInvitionMessage(
			receivers,
			strconv.Itoa(req.SenderId),
			msgCont,
			senderInfo.Nickname,
			senderInfo.AvatarURL,
			groupInfo.Name,
			groupInfo.AvatarURL,
			strconv.Itoa(req.GroupId),
			tokenStr,
		)

		// 发送消息到 MQ
		if err := shared.ForwardMessageToMQ(msg); err != nil {
			logrus.WithError(err).Error("Failed to send group invitation")
			c.JSON(500, gin.H{
				"error":   "Failed to send group invitation",
				"details": err.Error(),
			})
			return
		}

		// 存储到 Redis
		if err := shared.StoreInvitionToken(msg.MsgID, tokenStr); err != nil {
			logrus.WithError(err).Error("Failed to store invitation token")
			c.JSON(500, gin.H{
				"error":   "Failed to store invitation token",
				"details": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{"message": "Group invitation sent successfully"})
	default:
		c.JSON(400, gin.H{"error": "Invalid action", "details": "Action must be 'friend' or 'group'"})
	}

}

func AuthInvitation(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	if req.MsgId == "" || req.Token == "" {
		c.JSON(400, gin.H{"error": "Invalid parameters", "details": "msg_id and token cannot be empty"})
		return
	}

	targetToken, err := shared.SearchInvitationToken(req.MsgId)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve invitation token", "details": err.Error()})
		return
	}
	if req.Token != targetToken {
		c.JSON(400, gin.H{"error": "Invalid token", "details": "Token does not match"})
		return
	}

	if err := shared.DeleteInvitionToken(req.MsgId); err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete invitation token", "details": err.Error()})
		return
	}

	if req.Action == "reject" {
		c.JSON(200, gin.H{"message": "Request rejected successfully"})
		return
	}

	switch req.Type {
	case "friend":
		// Add friend relationship
		err := shared.MysqlDb.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec(`INSERT INTO user_friends (user_detail_id, friend_id) VALUES (?, ?)`, req.SenderId, req.ReceiverId).Error; err != nil {
				return fmt.Errorf("failed to add friend: %w", err)
			}
			if err := tx.Exec(`INSERT INTO user_friends (user_detail_id, friend_id) VALUES (?, ?)`, req.ReceiverId, req.SenderId).Error; err != nil {
				return fmt.Errorf("failed to add reverse friend: %w", err)
			}
			return nil
		})
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to add friend", "details": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Friend request accepted successfully"})

	case "group":
		if err := shared.MysqlDb.Exec(`INSERT INTO group_members (group_id, user_detail_id) VALUES (?, ?)`, req.GroupId, req.SenderId).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to add group member", "details": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Group request accepted successfully"})

	default:
		c.JSON(400, gin.H{"error": "Invalid type", "details": "Type must be 'friend' or 'group'"})
	}
}
