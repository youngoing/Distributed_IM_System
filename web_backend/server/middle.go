package server

import (
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// session := sessions.Default(c)
		// userID := session.Get("user_detail_id")

		// // 检查用户是否已登录
		// if userID == nil {
		// 	c.JSON(401, gin.H{"message": "Unauthorized"})
		// 	c.Abort()
		// 	return
		// }

		// 如果已登录，则允许继续请求
		c.Next()
	}
}

// func AuthToken() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// 从 Authorization 头部获取 token
// 		tokenString := c.GetHeader("Authorization")
// 		if tokenString == "" {
// 			// 如果没有提供 token，返回 401 错误
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
// 			c.Abort()
// 			return
// 		}

// 		// 移除 Bearer 前缀，保留 token 部分
// 		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
// 			tokenString = tokenString[7:]
// 		}

// 		// 解析 token
// 		claims, err := shared.ParseToken(tokenString)
// 		if err != nil {
// 			// 如果 token 无效，返回 401 错误
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 			c.Abort()
// 			return
// 		}

// 		// 将 claims 保存到上下文中供后续使用
// 		c.Set("userId", claims.UserId)

// 		// 继续处理请求
// 		c.Next()
// 	}
// }

// func UserVerificationMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		session := sessions.Default(c)
// 		sessionUserID := session.Get("user_detail_id")

// 		// 假设我们期望通过 URL 参数传递 user_id，例如 /user/:user_id
// 		requestUserID := c.Param("user_detail_id")

// 		// 验证用户是否为本人
// 		if sessionUserID != requestUserID {
// 			c.JSON(403, gin.H{"message": "Forbidden"})
// 			c.Abort()
// 			return
// 		}

// 		// 如果验证通过，继续请求
// 		c.Next()
// 	}
// }
