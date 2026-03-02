/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package mq

import (
	"CyberMatchmaker/config"
	"CyberMatchmaker/mapper"
	"CyberMatchmaker/middleware"
	"CyberMatchmaker/model"
	"CyberMatchmaker/pkg/utils"
	"context"
	"encoding/json"
	"log"

	"go.uber.org/zap"
)

// StartConsumer 启动监听服务
func StartConsumer(p *FortuneProducer) {
	msgs, err := p.channel.Consume(
		p.queue, "", true, false, false, false, nil,
	)
	if err != nil {
		log.Fatal("无法注册消费者:", err)
	}

	// 开启协程异步处理
	go func() {
		for d := range msgs {
			var task model.FortuneTaskMessage
			if err := json.Unmarshal(d.Body, &task); err != nil {
				log.Println("消息解析失败:", err)
				continue
			}

			// --- 核心逻辑开始 ---
			// TODO: 1. 根据 task.RecordID 从数据库查出完整记录
			record, err := mapper.GetFortuneRecordByID(task.RecordID)
			if err != nil {
				zap.S().Error("查询不到记录 ID %d: %v", task.RecordID, err)
			}
			// TODO: 2. 调用大模型 API (如 OpenAI/DeepSeek)
			zap.S().Info("正在处理用户 %d 的订单 %d...", task.UserID, task.RecordID)
			sysPrompt := config.GetPrompt("fortune_task.system")
			userPrompt := config.GetPrompt("fortune_task.user")
			AIResponse, err := middleware.LLM.CallAI(context.Background(), sysPrompt, userPrompt)
			if err != nil {
				zap.S().Error("调用大模型失败: %v", err)
			}
			// TODO: 3. 将大模型返回的结果封装回 FortuneRecord
			utils.CleanMarkdown(&AIResponse)
			utils.StringtoClass(AIResponse, record)
			// TODO: 4. 更新数据库：status = 'completed', 填充 Bazi, RecentFortune 等
			mapper.UpdateFortuneRecord(record)
			zap.S().Info("订单 %d 处理完成，结果已更新到数据库", task.RecordID)
		}
	}()
}
