package main

import (
	"CyberMatchmaker/config"
	"CyberMatchmaker/pkg/app"
	"CyberMatchmaker/pkg/infra"
	"CyberMatchmaker/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// 1. 必选的基础初始化
	config.InitConfig()
	logger.InitLogger()
	defer zap.L().Sync()

	// 2. 基础设施初始化 (内部调用你已写好的 InitDB 和 InitRedis)
	infra.InitAll()
	defer infra.CloseAll()

	// 3. 运行应用生命周期 (包含 Web、MQ 和优雅退出)
	app.Run()
}
