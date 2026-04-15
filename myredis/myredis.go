package myredis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"goweb/config"
	"log"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

// redis初始化
func InitRedis() bool {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.RAddr,     // Redis 地址
		Password: config.RPassword, // 密码（默认无密码）
		DB:       config.RDB,       // 数据库编号
		PoolSize: 25,               // 连接池大小
	})

	// 测试连接
	ctx := context.Background()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		log.Fatal("Redis 连接失败:", err)
		return false
	}
	log.Println("Redis 连接成功")
	return true
}
