package wsserver

import (
	"database/sql"
	"net/http"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

var (
	RedisConn *redis.Client
	DbConn    *sql.DB
)

// ResourceManager: control the resources
type ResourceManagerStruct struct {
	MqConn    *amqp.Connection
	MqChannel *amqp.Channel
}

// UserManager: manage the users
type UserManagerStruct struct {
	userChannels map[string]chan []byte
	userList     map[string]string // 存储在线的用户ID
	mu           sync.RWMutex      // 使用 读写锁
}

// WebSocketServer: the server
type WebSocketServerStruct struct {
	host            string
	httpPort        string
	webSocketPort   string
	WebSocketPath   string
	nodeId          string
	resourceManager *ResourceManagerStruct
	userManager     *UserManagerStruct
	upgrader        websocket.Upgrader
}

func NewWebSocketServer(host, webSocketPort, httpPort, path, nodeId string) *WebSocketServerStruct {
	return &WebSocketServerStruct{
		host:          host,
		httpPort:      httpPort,
		webSocketPort: webSocketPort,
		WebSocketPath: path,
		nodeId:        nodeId,
		resourceManager: &ResourceManagerStruct{
			MqConn:    nil,
			MqChannel: nil,
		},
		userManager: &UserManagerStruct{
			userChannels: make(map[string]chan []byte),
			userList:     make(map[string]string), // 初始化哈希表
			mu:           sync.RWMutex{},
		},
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}
