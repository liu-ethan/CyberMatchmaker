/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package mq

import (
	"context"
	"encoding/json"
	"sync"

	"CyberMatchmaker/mapper"
	"CyberMatchmaker/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// 增加一个 llmFunc 参数，类型为接收 string 返回 (string, error) 的函数
func StartFortuneConsumer(ctx context.Context, wg *sync.WaitGroup, llmFunc func(string) (string, error)) {
	q, err := rabbitmq.Ch.QueueDeclare(
		FortuneQueueName, true, false, false, false, nil,
	)
	if err != nil {
		zap.S().Fatalf("Failed to declare a queue: %v", err)
	}

	err = rabbitmq.Ch.Qos(1, 0, false)
	msgs, err := rabbitmq.Ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		zap.S().Fatalf("Failed to register a consumer: %v", err)
	}

	zap.S().Info("算命队列消费者已启动，等待接收任务...")

	for {
		select {
		case <-ctx.Done():
			zap.S().Info("收到退出信号，消费者停止接收新任务。")
			return
		case d, ok := <-msgs:
			if !ok {
				return
			}
			wg.Add(1)
			// 把 llmFunc 一并传给处理具体的 Worker 协程
			go processTask(d, wg, llmFunc)
		}
	}
}

// processTask 接收这个回调函数
func processTask(d amqp.Delivery, wg *sync.WaitGroup, llmFunc func(string) (string, error)) {
	defer wg.Done()

	var task FortuneTask
	if err := json.Unmarshal(d.Body, &task); err != nil {
		zap.S().Errorf("消息反序列化失败: %v", err)
		d.Reject(false)
		return
	}

	zap.S().Infof("开始处理算命任务, record_id: %d", task.RecordID)

	mapper.UpdateFortuneRecord(task.RecordID, map[string]interface{}{"status": "in_process"})

	// 核心改变：不直接依赖 service 包，而是调用外面传进来的回调函数
	resultText, err := llmFunc(task.Prompt)
	if err != nil {
		zap.S().Errorf("LLM API 调用失败, record_id: %d, err: %v", task.RecordID, err)
		mapper.UpdateFortuneRecord(task.RecordID, map[string]interface{}{"status": "failed"})
		d.Ack(false)
		return
	}

	mapper.UpdateFortuneRecord(task.RecordID, map[string]interface{}{
		"status":      "completed",
		"description": resultText,
	})

	zap.S().Infof("算命任务处理完成, record_id: %d", task.RecordID)
	d.Ack(false)
}
