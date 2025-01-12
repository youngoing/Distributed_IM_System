package handler

import (
	"backend/shared"
	"fmt"
	"net/http"

	"github.com/DanPlayer/randomname"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// // RegisterUser handles user registration with transaction and error handling
func RegisterUserHandler(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required"`
	}
	// 解析 JSON 请求
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}
	// 开始事务
	tx := shared.MysqlDb.Begin()
	// 创建用户对象
	user := User{
		Username: request.Username,
		Password: request.Password,
		Email:    request.Email,
	}
	// 哈希密码
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password"})
		return
	}
	user.Password = hashedPassword
	// 保存用户到数据库
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
		return
	}
	// 创建用户详情
	userDetail := UserDetail{
		UserID:   user.ID, // 使用 user.ID 作为关联
		Nickname: randomname.GenerateName(),
	}
	// 保存用户详情到数据库
	if err := tx.Create(&userDetail).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user details"})
		return
	}
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to commit transaction"})
		return
	}
	var response struct {
		Message string `json:"message"`
		WsToken string `json:"ws_token"`
	}
	// 生成 JWT
	token, err := shared.GenerateWsToken(fmt.Sprint(userDetail.ID), userDetail.Nickname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token"})
		return
	}
	response.Message = "User created successfully"
	response.WsToken = token
	// 返回成功响应
	c.JSON(http.StatusOK, response)
}

// LoginHandler handles user login with session management
func LoginHandler(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Parse JSON request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Find user by username
	var user User
	if err := shared.MysqlDb.Where("username = ?", request.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		return
	}

	// Validate password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		return
	}

	//获取user_detail
	var userDetail UserDetail
	if err := shared.MysqlDb.Where("user_id = ?", user.ID).First(&userDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	type UserLoginDetail struct {
		UserDetailID int    `json:"user_detail_id"`
		Nickname     string `json:"nickname"`
		AvatarURL    string `json:"avatar_url"`
	}
	userLoginDetail := UserLoginDetail{
		UserDetailID: userDetail.ID,
		Nickname:     userDetail.Nickname,
		AvatarURL:    userDetail.AvatarURL,
	}
	// Set session for user
	session := sessions.Default(c)
	session.Set("user_detail_id", userDetail.ID) // 将用户 ID 存储在 session 中

	session.Save()

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully", "user": userLoginDetail})

}

// LogoutHandler handles user logout by clearing the session
func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// UpdateUserDetails handles updating user details such as nickname and avatar
func UpdateUserDetails(c *gin.Context) {
	var request struct {
		Nickname  string `json:"nickname"`
		AvatarURL string `json:"avatar_url"`
	}

	// Parse JSON request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Retrieve user ID from session
	session := sessions.Default(c)
	userIDStr := session.Get("user_detail_id")
	if userIDStr == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not logged in"})
		return
	}

	// Find the user in the database
	var user User
	if err := shared.MysqlDb.Where("id = ?", userIDStr).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	// Update user details
	userDetail := UserDetail{
		UserID:    user.ID,
		Nickname:  request.Nickname,
		AvatarURL: request.AvatarURL,
	}

	// If the user detail does not exist, create a new one
	if err := shared.MysqlDb.Save(&userDetail).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User details updated successfully"})
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func GetGroupDetails(c *gin.Context) {
	// 从 URL 参数获取 user_detail_id
	userDetailId := c.Param("user_detail_id")

	// 定义一个数组来存储多个群组的信息
	var groups []struct {
		GroupID     int    `json:"group_id"`
		Name        string `json:"name"`
		AvatarURL   string `json:"avatar_url"`
		Description string `json:"description"`
	}

	// 修改 SQL 查询来查找用户参与的所有群组
	queryGroup := `
        SELECT g.id AS group_id, g.name, g.avatar_url, g.description
        FROM chat_groups g
        WHERE g.id IN (
            SELECT gm.group_id 
            FROM group_members gm 
            WHERE gm.user_detail_id = ?
        );
    `

	// 执行查询并将结果存储到数组中
	if err := shared.MysqlDb.Raw(queryGroup, userDetailId).Scan(&groups).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve group details", "details": err.Error()})
		return
	}

	// 如果没有群组，返回空数组
	if groups == nil {
		groups = []struct {
			GroupID     int    `json:"group_id"`
			Name        string `json:"name"`
			AvatarURL   string `json:"avatar_url"`
			Description string `json:"description"`
		}{}
	}

	// 遍历每个群组，获取该群组的成员信息
	var groupDetails []map[string]interface{} // 用于存储每个群组及其成员的信息

	for _, group := range groups {
		// 查询该群组的成员（排除请求该信息的用户）
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
            ) AND u.id != ?;  -- 排除请求该群组详情的用户
        `

		// 执行查询以获取群组成员
		if err := shared.MysqlDb.Raw(queryMembers, group.GroupID, userDetailId).Scan(&members).Error; err != nil {
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

		// 将每个群组的数据添加到结果数组
		groupDetails = append(groupDetails, groupData)
	}

	// 如果没有群组详情，返回空数组
	if groupDetails == nil {
		groupDetails = []map[string]interface{}{}
	}

	// 返回包含多个群组及其成员的响应
	c.JSON(200, groupDetails)
}

func GetUserFriends(c *gin.Context) {
	// Get the user_id from the URL parameter
	userDetailID := c.Param("user_detail_id")

	// Query to find the friends of the user
	query := `
        SELECT ud.id as user_detail_id, ud.nickname, ud.avatar_url
        FROM user_friends uf
        JOIN user_details ud ON uf.friend_id = ud.id
        WHERE uf.user_detail_id = ?
    `

	var friends []struct {
		UserDetailID int    `json:"user_detail_id"`
		Nickname     string `json:"nickname"`
		AvatarURL    string `json:"avatar_url"`
	}

	// Execute the query
	if err := shared.MysqlDb.Raw(query, userDetailID).Scan(&friends).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve friends"})
		return
	}

	// Return the friends' details
	c.JSON(200, friends)
}

func AuthLogged(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_detail_id")

	if userID == nil {
		c.JSON(401, gin.H{"message": "Unauthorized"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"message": "Authorized"})
}
