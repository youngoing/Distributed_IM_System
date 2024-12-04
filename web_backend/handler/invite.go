package handler

import (
	"backend/shared"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateFriendApplication 创建好友申请
func CreateFriendApplication(c *gin.Context) {

}

// CreateGroupApplication 创建加入群组申请
func CreateGroupInApplication(c *gin.Context) {

}

// AddMember 将成员添加到群组
func AddMember(c *gin.Context) {
	var request AddMemberRequest

	// 解析并验证请求体 JSON 数据
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	// 插入成员到 group_members 表
	query := `
        INSERT INTO group_members (group_id, user_detail_id)
        VALUES (?, ?)
    `
	if err := shared.MysqlDb.Exec(query, request.GroupID, request.UserDetailID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to add member", "error": err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "Member added successfully"})
}
