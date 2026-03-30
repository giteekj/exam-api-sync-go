package redis

import (
	"context"
	"fmt"
	"log"

	"exam-api-sync-go/common/setting"

	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
)

// Init 初始化Redis连接
func Init() {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", setting.Redis.Host, setting.Redis.Port),
		Password: setting.Redis.Password,
		DB:       setting.Redis.Db,
	})

	// 测试连接
	ctx := context.Background()
	_, err := Client.Ping(ctx).Result()
	if err != nil {
		log.Printf("Redis连接失败: %v", err)
	} else {
		log.Println("Redis连接成功")
	}
}

// GetClient 获取Redis客户端
func GetClient() *redis.Client {
	return Client
}
