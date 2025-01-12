package handler

import (
	"backend/shared"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// CreateGroup creates a new chat group and adds the user to the group
func CreateGroup(c *gin.Context) {
	// Define request structure
	var request struct {
		UserDetailID int    `json:"user_detail_id" binding:"required"` // Use int for UserDetailID
		Name         string `json:"name" binding:"required"`           // Group name
		AvatarURL    string `json:"avatar_url"`                        // Group avatar URL
		Description  string `json:"description"`                       // Group description (optional)
	}
	// Parse the incoming JSON request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	// Create the group instance
	group := ChatGroup{
		UserDetailID: request.UserDetailID, // Use the provided UserDetailID
		Name:         request.Name,
		AvatarURL:    request.AvatarURL,
		Description:  request.Description,
	}

	// Insert the group into the database
	if err := shared.MysqlDb.Create(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create group", "error": err.Error()})
		return
	}

	// Now that the group is created, add the user to the group
	groupMember := GroupMember{
		GroupID:      group.ID,             // Use the ID of the newly created group
		UserDetailID: request.UserDetailID, // Add the user to the group
	}

	if err := shared.MysqlDb.Create(&groupMember).Error; err != nil {
		// If adding the user to the group fails, delete the created group
		shared.MysqlDb.Delete(&group) // Rollback the group creation
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to add user to group", "error": err.Error()})
		return
	}

	// Return success response with the group details and user added
	c.JSON(http.StatusOK, gin.H{
		"message": "Group created and user added successfully",
		"group":   group,
	})
}

// UpdateGroup updates an existing chat group
func UpdateGroup(c *gin.Context) {
	// Retrieve the group ID from URL params
	groupIDParam := c.Param("id")
	groupID, err := strconv.Atoi(groupIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid group ID", "error": err.Error()})
		return
	}

	// Define the update request structure
	var request UpdateGroupRequest

	// Parse the incoming JSON request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	// Find the group from the database
	group, err := findGroupByID(groupID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Group not found", "error": err.Error()})
		return
	}

	// Update group fields if provided
	updateGroupFields(&group, request)

	// Save the updated group
	if err := shared.MysqlDb.Save(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update group", "error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Group updated successfully", "group": group})
}

// findGroupByID retrieves a group by its ID
func findGroupByID(groupID int) (ChatGroup, error) {
	var group ChatGroup
	if err := shared.MysqlDb.Where("id = ?", groupID).First(&group).Error; err != nil {
		return group, err
	}
	return group, nil
}

// updateGroupFields updates the group fields based on the request data
func updateGroupFields(group *ChatGroup, request UpdateGroupRequest) {
	if request.Name != "" {
		group.Name = request.Name
	}
	if request.AvatarURL != "" {
		group.AvatarURL = request.AvatarURL
	}
	if request.Description != "" {
		group.Description = request.Description
	}
}

// 定义请求结构体
type AddMemberRequest struct {
	GroupID      int `json:"group_id" binding:"required"`
	UserDetailID int `json:"user_detail_id" binding:"required"`
}

func GetGroupByID(c *gin.Context) {
	// 从 URL 参数获取 group_id
	groupID := c.Param("group_id")

	// 定义一个变量来存储群组的信息
	var group struct {
		GroupID     int    `json:"group_id"`
		Name        string `json:"name"`
		AvatarURL   string `json:"avatar_url"`
		Description string `json:"description"`
	}

	// 修改 SQL 查询来查找指定的群组
	queryGroup := `
		SELECT g.id AS group_id, g.name, g.avatar_url, g.description
		FROM chat_groups g
		WHERE g.id = ?;
	`

	// 执行查询并将结果存储到变量中
	if err := shared.MysqlDb.Raw(queryGroup, groupID).Scan(&group).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve group details", "details": err.Error()})
		return
	}

	// 如果群组不存在，返回错误
	if group.GroupID == 0 {
		c.JSON(404, gin.H{"error": "Group not found"})
		return
	}

	// 查询该群组的成员信息
	var members []struct {
		UserDetailID int    `json:"user_detail_id"`
		Nickname     string `json:"nickname"`
		AvatarURL    string `json:"avatar_url"`
	}

	queryMembers := `
		SELECT u.id AS user_detail_id, u.nickname, u.avatar_url
		FROM user_details u
		WHERE u.id IN (
			SELECT gm.user_detail_id
			FROM group_members gm
			WHERE gm.group_id = ?
		);
	`

	// 执行查询以获取群组成员
	if err := shared.MysqlDb.Raw(queryMembers, group.GroupID).Scan(&members).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve group members", "details": err.Error()})
		return
	}

	// 如果没有成员，返回空数组
	if members == nil {
		members = []struct {
			UserDetailID int    `json:"user_detail_id"`
			Nickname     string `json:"nickname"`
			AvatarURL    string `json:"avatar_url"`
		}{}
	}

	// 构造包含群组详情和成员信息的响应
	groupData := map[string]interface{}{
		"group_id":    group.GroupID,
		"name":        group.Name,
		"avatar_url":  group.AvatarURL,
		"description": group.Description,
		"members":     members, // 成员信息数组
	}

	// 返回包含群组详情和成员的响应
	c.JSON(200, groupData)
}

func DeleteGroup(c *gin.Context) {
	// 从 URL 参数获取 group_id
	groupID := c.Param("group_id")
	session := sessions.Default(c)
	userDetailID := session.Get("user_detail_id").(int)

	// 判断用户是否是群组的创建者
	var group ChatGroup
	if err := shared.MysqlDb.Where("id = ? AND user_detail_id = ?", groupID, userDetailID).First(&group).Error; err != nil {
		c.JSON(400, gin.H{"error": "Failed to delete group", "details": err.Error()})
		return
	}
	// 删除群组成员
	if err := shared.MysqlDb.Where("group_id = ?", groupID).Delete(&GroupMember{}).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete group members", "details": err.Error()})
		return
	}

	// 删除群组
	if err := shared.MysqlDb.Where("id = ?", groupID).Delete(&ChatGroup{}).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete group", "details": err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(200, gin.H{"message": "Group deleted successfully"})

}

func QuitGroup(c *gin.Context) {
	// 从 URL 参数获取 group_id
	groupID := c.Param("group_id")
	session := sessions.Default(c)
	userDetailID := session.Get("user_detail_id").(int)

	//判断用户是否是群组的创建者
	var group ChatGroup
	if err := shared.MysqlDb.Where("id = ? AND user_detail_id = ?", groupID, userDetailID).First(&group).Error; err != nil {
		c.JSON(400, gin.H{"error": "Failed to quit group", "details": err.Error()})
		return
	}

	// 删除群组成员
	if err := shared.MysqlDb.Where("group_id = ? AND user_detail_id = ?", groupID, userDetailID).Delete(&GroupMember{}).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to quit group", "details": err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(200, gin.H{"message": "Quit group successfully"})
}
