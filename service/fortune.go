/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package service

import (
	"CyberMatchmaker/config"
	"CyberMatchmaker/mapper"
	"CyberMatchmaker/middleware"
	"CyberMatchmaker/model"
	"CyberMatchmaker/model/modelDTO"
	"CyberMatchmaker/mq"
	"CyberMatchmaker/pkg/utils"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/jinzhu/copier"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type FortuneMQ struct {
	RecordID int64
	UserID   int64
}

// SubmitFortune 处理提交算命参数的业务逻辑，提交后返回recordID，结果通过异步处理后存储在数据库中
func SubmitFortune(c *gin.Context, userID int64) (int64, error) {
	var req modelDTO.SubmitFortuneRequestDTO

	// 1. 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		return -1, err
	}

	// 2. 绑定SubmitFortuneRequestDTO到FortuneRecord
	var record model.FortuneRecord
	_ = copier.Copy(&record, req)
	parsedDate, _ := time.Parse("2006-01-02", req.BirthDate) // 处理时间格式
	record.BirthDate = parsedDate
	record.Status = "pending" // 初始状态为 pending
	record.UserID = userID

	// 3. 保存FortuneRecord到数据库
	_ = mapper.CreateFortuneRecord(&record)

	// 4. 发送到消息队列（RabbitMQ）进行异步处理
	msg := FortuneMQ{
		RecordID: record.ID,
		UserID:   userID,
	}
	data, _ := json.Marshal(msg)
	err := mq.Publish("", config.AppConfig.RabbitMQ.FortuneQName, data)
	if err != nil {
		zap.S().Info("生产者发送消息失败: ", err)
		return -1, err
	}

	return record.ID, nil
}

// ConsumeFortune 处理从消息队列中接收的算命请求，进行算命逻辑处理，并将结果存储回数据库
// 只能调用一次, 在infra.go里面调用
func ConsumeFortune() {
	mq.Consume(config.AppConfig.RabbitMQ.FortuneQName, ConsumeHandleFortune)
}

// ConsumeHandleFortune 业务代码 处理从消息队列中接收的算命请求，进行算命逻辑处理，并将结果存储回数据库
// ConsumeFortuned的Handler函数，包含核心业务逻辑：
func ConsumeHandleFortune(d amqp.Delivery) {
	zap.S().Info("收到订单消息: ", string(d.Body))

	var data FortuneMQ
	if err := json.Unmarshal(d.Body, &data); err != nil {
		zap.S().Info("消息解析失败: ", err)
		return
	}

	// 1. 根据 data.RecordID 从数据库查出完整记录
	record, err := mapper.GetFortuneRecordByID(data.RecordID)
	if err != nil {
		zap.S().Error("查询不到记录 ID %d: %v", data.RecordID, err)
	}

	// 2. 调用大模型 API (如 OpenAI/DeepSeek)
	zap.S().Info("正在处理用户 %d 的订单 %d...", data.UserID, data.RecordID)
	sysPrompt := config.GetPrompt("fortune_task.system")
	userPrompt := config.GetPrompt("fortune_task.user")
	AIResponse, err := middleware.LLM.CallAI(context.Background(), sysPrompt, userPrompt)
	if err != nil {
		zap.S().Error("调用大模型失败: %v", err)
	}

	// 3. 将大模型返回的结果封装回 FortuneRecord
	utils.CleanMarkdown(&AIResponse)
	utils.StringtoClass(AIResponse, record)
	// 4. 更新数据库：status = 'completed', 填充 Bazi, RecentFortune 等
	record.Status = "completed"
	mapper.UpdateFortuneRecord(record)

	zap.S().Info("订单 %d 处理完成，结果已更新到数据库", data.RecordID)
}

// GetFortuneResult 根据用户ID查询最新的算命记录
func GetFortuneResult(c *gin.Context, userID int64) (modelDTO.FortuneResponseDTO, error) {
	record, err := mapper.GetLatestFortuneRecordByUserID(userID)
	if err != nil {
		return modelDTO.FortuneResponseDTO{}, fmt.Errorf("查询不到算命记录: %v", err)
	}
	var responseData modelDTO.FortuneResponseDTO
	_ = copier.Copy(&responseData, record)
	return responseData, nil
}
