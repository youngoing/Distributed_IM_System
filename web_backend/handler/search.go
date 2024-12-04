package handler

import (
	"backend/shared"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SearchUserOrGroup(c *gin.Context) {
	// 获取查询参数
	query := c.Query("query")

	// 初始化结果容器
	var users []UserDetail
	var groups []ChatGroup

	// 如果没有查询参数，或者查询为空，返回所有用户和群组
	if query == "" {
		// 获取所有用户
		if err := shared.MysqlDb.Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users", "details": err.Error()})
			return
		}
		// 获取所有群组
		if err := shared.MysqlDb.Find(&groups).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve groups", "details": err.Error()})
			return
		}
	} else {
		// 搜索用户
		if err := shared.MysqlDb.Where("nickname LIKE ?", "%"+query+"%").Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users", "details": err.Error()})
			return
		}

		// 搜索群组
		if err := shared.MysqlDb.Where("name LIKE ?", "%"+query+"%").Find(&groups).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search groups", "details": err.Error()})
			return
		}
	}

	// 返回搜索结果
	c.JSON(http.StatusOK, gin.H{
		"users":  formatUserDetails(users), // 格式化用户数据
		"groups": formatChatGroups(groups), // 格式化群组数据
	})
}

// 格式化用户数据，调整字段名
func formatUserDetails(users []UserDetail) []map[string]interface{} {
	var formattedUsers []map[string]interface{}
	for _, user := range users {
		formattedUsers = append(formattedUsers, map[string]interface{}{
			"user_detail_id": user.UserID,
			"nickname":       user.Nickname,
			"avatar_url":     user.AvatarURL,
			"created_at":     user.CreatedAt,
		})
	}
	return formattedUsers
}

// 格式化群组数据，调整字段名
func formatChatGroups(groups []ChatGroup) []map[string]interface{} {
	var formattedGroups []map[string]interface{}
	for _, group := range groups {
		formattedGroups = append(formattedGroups, map[string]interface{}{
			"group_id":       group.ID,
			"user_detail_id": group.UserDetailID,
			"name":           group.Name,
			"avatar_url":     group.AvatarURL,
			"description":    group.Description,
			"created_at":     group.CreatedAt,
		})
	}
	return formattedGroups
}
