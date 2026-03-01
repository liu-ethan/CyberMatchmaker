/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package rabbitmq

import (
	global "CyberMatchmaker/pkg"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"sync"
)

var mqOnce sync.Once

// InitRabbitMQ 初始化 RabbitMQ 连接
func InitRabbitMQ(url string) {
	mqOnce.Do(func() {
		conn, err := amqp.Dial(url)
		if err != nil {
			log.Fatalf("RabbitMQ 连接失败: %v", err)
		}

		global.MQ = conn
		log.Println("RabbitMQ 初始化成功")
	})
}
