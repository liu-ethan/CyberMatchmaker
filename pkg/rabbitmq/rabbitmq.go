/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package rabbitmq

import (
	"CyberMatchmaker/config"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ 封装对象
type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	URL     string
}

var MQ *RabbitMQ

// NewRabbitMQ 初始化连接
func NewRabbitMQ() {
	mq := &RabbitMQ{URL: config.AppConfig.RabbitMQ.URL}
	mq.connect()
	MQ = mq
}

// 内部连接逻辑（含初始化 Channel）
func (r *RabbitMQ) connect() {
	var err error
	r.Conn, err = amqp.Dial(r.URL)
	if err != nil {
		log.Fatalf("无法连接 RabbitMQ: %v", err)
	}
	r.Channel, err = r.Conn.Channel()
	if err != nil {
		log.Fatalf("无法打开 Channel: %v", err)
	}
}

// Close 关闭连接
func (r *RabbitMQ) Close() {
	r.Channel.Close()
	r.Conn.Close()
}
