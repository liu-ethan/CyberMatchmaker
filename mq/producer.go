/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package mq

import (
	"CyberMatchmaker/pkg/rabbitmq"
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const FortuneQueueName = "fortune_task_queue"

// FortuneTask 在队列中流转的载荷
type FortuneTask struct {
	RecordID int64  `json:"record_id"`
	UserID   int64  `json:"user_id"`
	Prompt   string `json:"prompt"`
}

// PublishFortuneTask 投递任务
func PublishFortuneTask(task FortuneTask) error {
	body, err := json.Marshal(task)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = rabbitmq.Ch.PublishWithContext(ctx,
		"",               // exchange (使用默认 direct 交换机)
		FortuneQueueName, // routing key (直接对应队列名)
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent, // 消息持久化，防止 MQ 宕机丢失
			Body:         body,
		})

	if err != nil {
		zap.S().Errorf("投递算命任务失败, record_id: %d, err: %v", task.RecordID, err)
		return err
	}

	zap.S().Infof("算命任务已投递入队, record_id: %d", task.RecordID)
	return nil
}
