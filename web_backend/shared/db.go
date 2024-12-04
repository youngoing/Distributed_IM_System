package shared

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // 导入 MySQL 驱动程序
)

var ctx = context.Background()
var MysqlDb *gorm.DB

// 初始化 MySQL 连接池
func InitDB() {
	// 数据库连接字符串
	dsn := "root:qweasdzxc@tcp(172.25.59.171:3306)/backend?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	// 使用 gorm 连接数据库
	MysqlDb, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// 配置连接池
	MysqlDb.DB().SetMaxIdleConns(10)    // 设置最大空闲连接数
	MysqlDb.DB().SetMaxOpenConns(100)   // 设置最大打开连接数
	MysqlDb.DB().SetConnMaxLifetime(10) // 设置连接的最大可复用时间（单位秒）

}

// 关闭数据库连接
func CloseDB() {
	if err := MysqlDb.Close(); err != nil {
		log.Fatalf("Failed to close database: %v", err)
	}
}

// 测试数据库连接
func TestDB() {
	if err := MysqlDb.DB().Ping(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	} else {
		log.Println("Successfully connected to the database!")
	}
}

var RedisClient *redis.Client

func InitRedis() {
	// redis连接
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "172.25.59.171:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
		PoolSize: 10, // 连接池大小
	})
}

func ClearSession() {
	// 删除带有 session_ 前缀的所有键
	iter := RedisClient.Scan(ctx, 0, "session_*", 0).Iterator()
	for iter.Next(ctx) {
		err := RedisClient.Del(ctx, iter.Val()).Err()
		if err != nil {
			fmt.Printf("Failed to delete key %s: %v\n", iter.Val(), err)
		}
	}
	if err := iter.Err(); err != nil {
		fmt.Printf("Failed to iterate keys: %v\n", err)
	}
}
