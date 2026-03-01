/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package service

import (
	"CyberMatchmaker/mapper"
	"CyberMatchmaker/model"
	"CyberMatchmaker/model/modelDTO"
	"CyberMatchmaker/mq"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// SubmitFortune 处理提交算命参数的业务逻辑，提交后返回recordID，结果通过异步处理后存储在数据库中
func SubmitFortune(c *gin.Context, userID int64) (int64, error) {
	var req modelDTO.SubmitFortuneRequestDTO

	// 1. 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		return -1, err
	}

	// 2. 绑定SubmitFortuneRequestDTO到FortuneRecord
	var record model.FortuneRecord
	copier.Copy(&record, req)
	parsedDate, _ := time.Parse("2006-01-02", req.BirthDate) // 处理时间格式
	record.BirthDate = parsedDate
	record.Status = "pending" // 初始状态为 pending
	record.UserID = userID

	// 3. 保存FortuneRecord到数据库
	mapper.CreateFortuneRecord(&record)

	// 4. 发送到消息队列（RabbitMQ）进行异步处理
	task := model.FortuneTaskMessage{
		UserID:   userID,
		RecordID: record.ID,
	}
	mq.GlobalProducer.PublishFortuneTask(task)

	return record.ID, nil
}
