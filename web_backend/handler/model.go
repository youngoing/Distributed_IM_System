package handler

import (
	"time"
)

// Define the request structure globally to reuse in multiple places
type UpdateGroupRequest struct {
	Name        string `json:"name"`
	AvatarURL   string `json:"avatar_url"`
	Description string `json:"description"`
}

// ConsulService 定义服务信息的结构
type ConsulService struct {
	ID      string        `json:"ID"`
	Name    string        `json:"Name"`
	Tags    []string      `json:"Tags,omitempty"`
	Address string        `json:"Address"`
	Port    int           `json:"Port"`
	Check   *ServiceCheck `json:"Check,omitempty"`
}

// ServiceCheck 定义服务健康检查的结构
type ServiceCheck struct {
	HTTP     string `json:"http"`
	Interval string `json:"interval"`
	Timeout  string `json:"timeout"`
}

// User 定义用户结构体
type User struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"` // 使用 int 类型的自增主键
	Username  string    `gorm:"type:varchar(255);unique;not null" json:"username"`
	Password  string    `gorm:"type:varchar(255);not null" json:"password"`
	Email     string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"` // 自动创建时间戳
}

// UserDetail 定义用户详情结构体
type UserDetail struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`     // 使用 int 类型的自增主键
	UserID    int       `gorm:"type:int;not null;index" json:"user_id"` // 将 UserID 改为 int 类型
	Nickname  string    `gorm:"type:varchar(255);not null" json:"nickname"`
	AvatarURL string    `gorm:"type:varchar(255);default:null" json:"avatar_url,omitempty"` // 可空字段，默认 null
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`                 // 自动创建时间戳
}

// ChatGroup 定义聊天群组结构体
type ChatGroup struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`            // 使用 int 类型的自增主键
	UserDetailID int       `gorm:"type:int;not null;index" json:"user_detail_id"` // 将 UserDetailID 改为 int 类型
	Name         string    `gorm:"type:varchar(255);not null" json:"name"`
	AvatarURL    string    `gorm:"type:varchar(255);not null" json:"avatar_url"`
	Description  string    `gorm:"type:text" json:"description"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// GroupMember 定义群组成员结构体
type GroupMember struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	GroupID      int       `gorm:"type:int;not null;index" json:"group_id"`
	UserDetailID int       `gorm:"type:int;not null;index" json:"user_detail_id"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// UserFriend 定义用户好友结构体
type UserFriend struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserDetailID int       `gorm:"type:int;not null;index" json:"user_detail_id"`
	FriendID     int       `gorm:"type:int;not null;index" json:"friend_id"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}
