package handler

import (
	"backend/shared"

	"github.com/gin-gonic/gin"
)

func ListFriends(c *gin.Context) {

}

func DeleteFriend(c *gin.Context) {
	// Delete a friend
	friend_id := c.Query("friend_id")
	user_detail_id := c.Query("user_id")
	query := `delete from user_friends where user_detail_id = ? and friend_id = ?`
	if err := shared.MysqlDb.Exec(query, user_detail_id, friend_id).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete friend"})
		return
	}
	if err := shared.MysqlDb.Exec(query, friend_id, user_detail_id).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete friend"})
		return
	}
	c.JSON(200, gin.H{"message": "Friend deleted successfully"})
}
func FriendDetail(c *gin.Context) {
	friendID := c.Param("friend_id")
	var friend struct {
		UserDetailID int    `json:"user_detail_id"`
		Nickname     string `json:"nickname"`
		AvatarURL    string `json:"avatar_url"`
	}

	// 使用 GORM 的 First 方法查找记录，假设 user_details 表对应的模型为 UserDetail
	if err := shared.MysqlDb.Table("user_details").
		Select("user_id as user_detail_id, nickname, avatar_url").
		Where("id = ?", friendID).
		First(&friend).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve friend details"})
		return
	}

	c.JSON(200, friend)
}
