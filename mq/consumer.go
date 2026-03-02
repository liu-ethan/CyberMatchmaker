/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package mq

import (
	"CyberMatchmaker/pkg/rabbitmq"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Consume 封装消费者：传入队列名和处理函数
func Consume(queueName string, handler func(d amqp.Delivery)) {
	msgs, err := rabbitmq.MQ.Channel.Consume(
		queueName,
		"",    // 消费者标签
		true,  // 自动应答 (Auto-Ack)
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Printf("消费队列失败: %v", err)
		return
	}
	// 开启协程异步处理
	go func() {
		for d := range msgs {
			handler(d)
		}
	}()
	log.Printf("消费者已启动，监听队列: %s", queueName)
}
