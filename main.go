package main

import (
	"CyberMatchmaker/config"
	"CyberMatchmaker/pkg/logger" // 引入刚才写的 logger 包
	"CyberMatchmaker/pkg/postgres"
	"CyberMatchmaker/pkg/rabbitmq"
	"CyberMatchmaker/pkg/redis"
	"CyberMatchmaker/route"
	"fmt"

	"go.uber.org/zap"
)

func main() {
	// 1. 加载配置
	config.InitConfig()
	// 2. 初始化全局日志 (务必紧跟在配置加载之后，让后续的组件都能用上日志)
	logger.InitLogger()
	// 延迟同步，确保程序退出前把缓存区里的日志全部刷入磁盘
	defer zap.L().Sync()

	zap.S().Info("日志模块初始化成功，开始加载基础设施...")

	// 3. 组装连接字符串并初始化基础设施连接池
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Shanghai",
		config.AppConfig.Database.Host,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.DBName,
		config.AppConfig.Database.SSLMode,
	)
	postgres.InitDB(dsn)

	redisAddr := fmt.Sprintf("%s:%d", config.AppConfig.Redis.Host, config.AppConfig.Redis.Port)
	redis.InitRedis(redisAddr, config.AppConfig.Redis.Password, config.AppConfig.Redis.DB)

	rabbitmq.InitRabbitMQ(config.AppConfig.RabbitMQ.URL)

	// 4. 注册路由 (Gin)
	r := route.SetupRouter()

	// 5. 启动 Web 服务
	serverPort := config.AppConfig.Server.Port
	addr := fmt.Sprintf(":%d", serverPort)

	zap.S().Infof("CyberMatchmaker 服务启动成功，监听端口 %s", addr)

	// 启动并监听 HTTP 请求
	if err := r.Run(addr); err != nil {
		zap.S().Fatalf("Web 服务启动失败: %v", err)
	}
}
