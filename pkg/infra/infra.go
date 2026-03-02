/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package infra

import (
	"CyberMatchmaker/config"
	"CyberMatchmaker/middleware"
	"CyberMatchmaker/pkg/postgres"
	"CyberMatchmaker/pkg/rabbitmq"
	"CyberMatchmaker/pkg/redis"
	"fmt"
)

// InitAll 统一初始化所有基础设施
func InitAll() {
	// 1. 拼接并初始化 Postgres
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Shanghai",
		config.AppConfig.Database.Host,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.DBName,
		config.AppConfig.Database.SSLMode,
	)
	postgres.InitDB(dsn)

	// 2. 拼接并初始化 Redis
	redisAddr := fmt.Sprintf("%s:%d", config.AppConfig.Redis.Host, config.AppConfig.Redis.Port)
	redis.InitRedis(redisAddr, config.AppConfig.Redis.Password, config.AppConfig.Redis.DB)

	// 3. 初始化 RabbitMQ
	rabbitmq.NewRabbitMQ()

	// 6. 初始化LLM
	middleware.NewLLMService()
}

func CloseAll() {
	// 关闭 RabbitMQ 连接
	rabbitmq.MQ.Close()
}
