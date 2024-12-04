package handler

import (
	"backend/shared"

	"github.com/gin-gonic/gin"
)

func AddFriend(c *gin.Context) {
	// Define a struct to bind the request body
	var request struct {
		UserDetailID int `json:"user_detail_id" binding:"required"`
		FriendID     int `json:"friend_id" binding:"required"`
	}

	// Bind the request data to the struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}

	// Prepare the SQL query to insert friendship into `user_friends` table
	query := `
        INSERT INTO user_friends (user_detail_id, friend_id)
        VALUES (?, ?), (?, ?)
    `

	// Execute the query to create the friendship
	if err := shared.MysqlDb.Exec(query, request.UserDetailID, request.FriendID, request.FriendID, request.UserDetailID).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to add friend"})
		return
	}

	// Return success response
	c.JSON(200, gin.H{"message": "Friend added successfully"})
}

func ListFriends(c *gin.Context) {

}

func DeleteFriend(c *gin.Context) {
	// Delete a friend
}
func FriendDetail(c *gin.Context) {
	friend_id := c.Param("friend_id")
	query := `select user_id as user_detail_id,nickname,avatar_url from user_details where user_id = ?`
	var friend struct {
		UserDetailID int    `json:"user_detail_id"`
		Nickname     string `json:"nickname"`
		AvatarURL    string `json:"avatar_url"`
	}
	if err := shared.MysqlDb.Raw(query, friend_id).Scan(&friend).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve friend details"})
		return
	}
	c.JSON(200, friend)
}
