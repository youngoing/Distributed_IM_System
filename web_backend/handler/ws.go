package handler

// func AuthToken(c *gin.Context) {
// 	// 从 Authorization 头部获取 token
// 	tokenString := c.GetHeader("Authorization")
// 	if tokenString == "" {
// 		// 如果没有提供 token，返回 401 错误
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
// 		c.Abort()
// 		return
// 	}

// 	// 移除 Bearer 前缀，保留 token 部分
// 	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
// 		tokenString = tokenString[7:]
// 	}

// 	// 解析 token
// 	claims, err := shared.ParseToken(tokenString)
// 	if err != nil {
// 		// 如果 token 无效，返回 401 错误
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 		c.Abort()
// 		return
// 	}
// 	var response struct {
// 		UserId string `json:"user_id"`
// 	}
// 	response.UserId = claims.UserId
// 	c.JSON(http.StatusOK, response)
// }

//发送ws消息
func sendWsMsg()
