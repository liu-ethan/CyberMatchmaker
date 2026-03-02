/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package redis

import (
	global "CyberMatchmaker/pkg"
	"context"
	"log"
	"sync"
	"time"

	redisClient "github.com/redis/go-redis/v9"
)

var redisOnce sync.Once

// InitRedis 初始化 Redis 客户端
func InitRedis(addr, password string, db int) {
	redisOnce.Do(func() {
		client := redisClient.NewClient(&redisClient.Options{
			Addr:         addr,
			Password:     password,
			DB:           db,
			PoolSize:     20, // 连接池大小
			MinIdleConns: 5,  // 最小空闲连接数
		})
		// 探活测试，验证连接是否真实可用
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := client.Ping(ctx).Result(); err != nil {
			log.Fatalf("Redis 连接失败: %v", err)
		}
		global.Redis = client
		log.Println("Redis 初始化成功")
	})
}
