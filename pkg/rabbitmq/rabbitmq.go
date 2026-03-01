/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package rabbitmq

import (
	"CyberMatchmaker/config" // 替换为你的 config 包路径

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// 定义全局的连接和通道变量，首字母大写，方便跨包(mq包)调用
var (
	Conn *amqp.Connection
	Ch   *amqp.Channel
)

// InitRabbitMQ 初始化 MQ 连接池
func InitRabbitMQ() {
	// 直接读取你 yaml 里配置好的完整 url
	// 注意大小写：取决于你 config.go 里定义的结构体字段是 Url 还是 URL
	dsn := config.AppConfig.RabbitMQ.URL

	var err error
	// 1. 直接使用这个 dsn 建立 Connection
	Conn, err = amqp.Dial(dsn)
	if err != nil {
		zap.S().Fatalf("连接 RabbitMQ 失败: %v", err)
	}

	// 2. 建立 Channel
	Ch, err = Conn.Channel()
	if err != nil {
		zap.S().Fatalf("打开 RabbitMQ Channel 失败: %v", err)
	}

	zap.S().Info("RabbitMQ 连接初始化成功！")
}

// Close 在 main.go 退出时优雅关闭连接
func Close() {
	if Ch != nil {
		_ = Ch.Close()
	}
	if Conn != nil {
		_ = Conn.Close()
	}
	zap.S().Info("RabbitMQ 连接已关闭")
}
