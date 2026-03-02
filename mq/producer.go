/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package mq

import (
	"CyberMatchmaker/pkg/rabbitmq"
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publish 封装生产者：支持指定交换机、路由键和消息内容
func Publish(exchange, routingKey string, body []byte) error {
	return rabbitmq.MQ.Channel.PublishWithContext(
		context.Background(),
		exchange,   // 交换机名称
		routingKey, // 路由键
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         body,
			DeliveryMode: amqp.Persistent, // 消息持久化
		},
	)
}
