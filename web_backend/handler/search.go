package handler

import (
	"backend/shared"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SearchUserOrGroup(c *gin.Context) {
	// 获取查询参数
	query := c.Query("query")
	searchType := c.Query("type") // 新增type参数: "user", "group", 或空

	// 初始化结果容器
	var users []UserDetail
	var groups []ChatGroup

	// 根据查询类型和关键词进行搜索
	switch searchType {
	case "user":
		// 只搜索用户
		if query == "" {
			if err := shared.MysqlDb.Find(&users).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "search user failed",
					"details": err.Error(),
				})
				return
			}
		} else {
			if err := shared.MysqlDb.Where("nickname LIKE ?", "%"+query+"%").Find(&users).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "search user failed",
					"details": err.Error(),
				})
				return
			}
		}

	case "group":
		// 只搜索群组
		if query == "" {
			if err := shared.MysqlDb.Find(&groups).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "search group failed",
					"details": err.Error(),
				})
				return
			}
		} else {
			if err := shared.MysqlDb.Where("name LIKE ?", "%"+query+"%").Find(&groups).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "search group failed",
					"details": err.Error(),
				})
				return
			}
		}

	default:
		// 搜索所有（用户和群组）
		if query == "" {
			// 获取所有用户
			if err := shared.MysqlDb.Find(&users).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "search user failed",
					"details": err.Error(),
				})
				return
			}
			// 获取所有群组
			if err := shared.MysqlDb.Find(&groups).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "search group failed",
					"details": err.Error(),
				})
				return
			}
		} else {
			// 搜索用户
			if err := shared.MysqlDb.Where("nickname LIKE ?", "%"+query+"%").Find(&users).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "search user failed",
					"details": err.Error(),
				})
				return
			}
			// 搜索群组
			if err := shared.MysqlDb.Where("name LIKE ?", "%"+query+"%").Find(&groups).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "search group failed",
					"details": err.Error(),
				})
				return
			}
		}
	}

	// 根据搜索类型返回结果
	response := gin.H{}
	if searchType != "group" {
		response["users"] = formatUserDetails(users)
	}
	if searchType != "user" {
		response["groups"] = formatChatGroups(groups)
	}

	c.JSON(http.StatusOK, response)
}

// 格式化用户数据，调整字段名
func formatUserDetails(users []UserDetail) []map[string]interface{} {
	var formattedUsers []map[string]interface{}
	for _, user := range users {
		formattedUsers = append(formattedUsers, map[string]interface{}{
			"user_detail_id": user.ID,
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
