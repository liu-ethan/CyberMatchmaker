/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package infra

import (
	"CyberMatchmaker/config"
	"CyberMatchmaker/mq"
	"CyberMatchmaker/pkg/postgres"
	"CyberMatchmaker/pkg/rabbitmq"
	"CyberMatchmaker/pkg/redis"
	"fmt"
)

// InitAll 统一初始化所有基础设施
func InitAll() {
	// 1. 拼接并初始化 Postgres [cite: 8]
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Shanghai",
		config.AppConfig.Database.Host,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.DBName,
		config.AppConfig.Database.SSLMode,
	)
	postgres.InitDB(dsn)

	// 2. 拼接并初始化 Redis [cite: 9]
	redisAddr := fmt.Sprintf("%s:%d", config.AppConfig.Redis.Host, config.AppConfig.Redis.Port)
	redis.InitRedis(redisAddr, config.AppConfig.Redis.Password, config.AppConfig.Redis.DB)

	// 3. 初始化 RabbitMQ [cite: 10]
	rabbitmq.InitRabbitMQ()

	// 4. 初始化 RabbitMQ 生产者
	mq.InitProducer(rabbitmq.Ch, config.AppConfig.RabbitMQ.QName)

	// 5. 启动 RabbitMQ 消息消费者
	mq.StartConsumer(mq.GlobalProducer)
}

// CloseAll 统一释放资源
func CloseAll() {
	rabbitmq.Close()
	// 如果你的 postgres/redis 也有 Close 方法也可以加在这里
}
