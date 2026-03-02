/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package mq

import (
	"CyberMatchmaker/model"

	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// FortuneProducer 负责连接 RabbitMQ 并发送消息到队列
type FortuneProducer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

var GlobalProducer *FortuneProducer

// InitProducer 初始化 RabbitMQ 生产者，连接 MQ 服务器并声明队列
func InitProducer(ch *amqp.Channel, qName string) {
	// 1. 声明队列 (这一步必不可少，确保队列存在)
	_, err := ch.QueueDeclare(
		qName, true, false, false, false, nil,
	)
	if err != nil {
		zap.S().Fatalf("无法在信道上声明队列: %v", err)
	}
	// 2. 封装进全局变量
	GlobalProducer = &FortuneProducer{
		channel: ch,    // 使用传入的信道
		queue:   qName, // 记录队列名
	}
}

// PublishFortuneTask 将算命任务消息发布到 RabbitMQ 队列中
func (p *FortuneProducer) PublishFortuneTask(msg model.FortuneTaskMessage) error {
	// 1. 将结构体转为 JSON 字符串
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	// 2. 发布消息
	return p.channel.Publish(
		"",      // exchange (空字符串表示使用默认交换机)
		p.queue, // routing key (队列名)
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // 消息持久化，防止 MQ 宕机丢失
		},
	)
}
