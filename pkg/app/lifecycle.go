/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package app

import (
	"CyberMatchmaker/config"
	"CyberMatchmaker/route"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// Run 启动整个应用并处理优雅退出
func Run() {
	// 准备上下文和 WaitGroup
	//ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// 1. 启动后台消费者
	//go mq.StartFortuneConsumer(ctx, &wg, service.GenerateFortune)

	// 2. 配置 Web 服务
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.AppConfig.Server.Port),
		Handler: route.SetupRouter(),
	}

	// 3. 异步启动 Web 服务
	go func() {
		zap.S().Infof("服务已启动，监听端口 %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.S().Fatalf("监听异常: %v", err)
		}
	}()

	// 4. 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 5. 优雅关机流程
	zap.S().Info("正在安全退出...")
	//cancel() // 通知消费者停止接收新消息

	// 设置 5 秒超时强行关闭 Web 服务
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)

	wg.Wait() // 等待所有已领取的 MQ 任务处理完成
	zap.S().Info("所有后台任务已结束，进程退出。")
}
