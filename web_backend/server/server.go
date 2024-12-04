package server

import (
	"backend/handler"
	"backend/shared"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func Run() {
	if err := shared.ConfirmEnv(); err != nil {
		log.Fatalf("Environment variable check failed: %v", err)
	}
	// 初始化数据库连接
	shared.InitDB()
	defer shared.CloseDB() // 确保在程序结束时关闭数据库连接
	shared.TestDB()        // 测试数据库连接是否成功
	// 初始化 Redis 连接
	shared.InitRedis()
	defer shared.RedisClient.Close()
	// 创建默认的 Gin 路由器
	router := gin.Default()
	// 配置 CORS 中间件
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:3000", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	//判断是否为int类型,如果不是就设置默认值24
	userExpirationHoursStr := os.Getenv("USER_EXPIRATION_HOURS")

	// 尝试将字符串转换为整数
	userExpirationHours, err := strconv.Atoi(userExpirationHoursStr)
	if err != nil {
		log.Println("Invalid USER_EXPIRATION_HOURS, using default value of 24 hours.")
		userExpirationHours = 24
	}
	// 使用 Cookie 作为会话存储，并设置 SameSite 属性为 None
	cookieStore := cookie.NewStore([]byte("secret"))
	cookieStore.Options(sessions.Options{
		Path:     "/",
		Domain:   "127.0.0.1",                   // 根据实际情况设置域名
		MaxAge:   60 * 60 * userExpirationHours, // s
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode, // 设置 SameSite 属性为 None
	})
	// 使用 Redis 作为会话存储
	redisStore, err := redis.NewStore(10, "tcp", "172.25.59.171:6379", "", []byte("secret")) //用于加密会话数据的密钥。
	if err != nil {
		panic(err)
	}
	router.Use(sessions.Sessions("backend_session", redisStore))

	consul := router.Group("/consul")
	{
		consul.POST("/register", handler.RegisterServiceHandler)
		consul.DELETE("/deregister/:id", handler.DeregisterServiceHandler)
		consul.GET("/services", handler.ListServicesHandler)
		consul.PUT("/update/:id", handler.UpdateServiceHandler)
	}

	user := router.Group("/user")

	{
		user.POST("/register", handler.RegisterUserHandler)
		user.POST("/login", handler.LoginHandler)
		user.GET("/logout", AuthMiddleware(), handler.LogoutHandler)
		user.POST("/:user_detail_id/update", AuthMiddleware(), handler.UpdateUserDetails)
		user.GET("/:user_detail_id/groups", AuthMiddleware(), handler.GetGroupDetails)
		user.GET("/:user_detail_id/friends", AuthMiddleware(), handler.GetUserFriends)
		user.GET("/auth", handler.AuthLogged)
	}

	group := router.Group("/group")

	{
		group.POST("/create", AuthMiddleware(), handler.CreateGroup)
		group.POST("/update/:id", AuthMiddleware(), handler.UpdateGroup)
		group.GET("/:group_id/detail", AuthMiddleware(), handler.GetGroupByID)
	}

	friend := router.Group("/friend")
	{
		friend.POST("create", AuthMiddleware(), handler.AddFriend)
		friend.GET("/:friend_id/detail", AuthMiddleware(), handler.FriendDetail)

	}
	search := router.Group("/search")
	{
		search.GET("/", AuthMiddleware(), handler.SearchUserOrGroup)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 捕获退出信号
	quit := make(chan os.Signal, 1)
	//捕获 SIGINT 和 SIGTERM 信号（即 Ctrl+C 信号）
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		log.Printf("Received signal: %v\n", <-sigChan)
		<-quit
		// 清除所有会话
		log.Println("Clearing all sessions...")
		shared.ClearSession()
		// 退出程序
		os.Exit(0)
	}()

	router.Run(":8080")
}
