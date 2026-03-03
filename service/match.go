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
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/jinzhu/copier"
	"github.com/pgvector/pgvector-go"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type MatchMQ struct {
	Profile     model.MatchProfile
	Description string
}

// JoinMatch 处理用户加入匹配的业务逻辑
func JoinMatch(c *gin.Context, userID int64) error {
	// 1. 检查用户是否已经开启第一次算命
	fortuneRecord, err := mapper.GetLatestFortuneRecordByUserID(userID)

	// 2. 如果没有开启第一次算命，返回错误
	if err != nil {
		// 判断是否是“记录不存在”的错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("您还没有算命记录，请先完成您的第一次算命")
		}
		// 其他数据库层面错误（如连接断开）
		return err
	}

	// 3. 拿到wechat_id
	var matchProfileDTO modelDTO.MatchProfileDTO
	c.ShouldBindJSON(&matchProfileDTO)

	// 4. 如果已经开启第一次算命，将用户加入匹配队列
	var matchProfile model.MatchProfile
	copier.Copy(&matchProfile, fortuneRecord)

	// 5. 补全matchProfile中的其他字段
	matchProfile.FortuneRecordID = fortuneRecord.ID
	matchProfile.WechatID = matchProfileDTO.WechatID
	matchProfile.City = fortuneRecord.CurrentCity

	// 6. 转为MatchMQ
	matchMQ := &MatchMQ{
		Profile:     matchProfile,
		Description: *fortuneRecord.Description,
	}

	// 6. 发送到消息队列，异步处理匹配逻辑
	body, _ := json.Marshal(&matchMQ)
	mq.Publish("", config.AppConfig.RabbitMQ.EmbeddingQName, body)

	// 7. 返回成功
	return nil
}

// JoinMatchConsume 消费者函数，监听JoinMatch的消息队列，处理匹配逻辑
// 只调用一次
func JoinMatchConsume() {
	mq.Consume(config.AppConfig.RabbitMQ.EmbeddingQName, JoinMatchConsumeHandler)
}

// JoinMatchConsumeHandler 处理JoinMatch的消息队列消费逻辑
func JoinMatchConsumeHandler(d amqp.Delivery) {
	// 1. 将消息体反序列化为MatchMQ对象
	var matchMQ MatchMQ
	_ = json.Unmarshal(d.Body, &matchMQ)

	// 2. 根据fortuneRecord.Description计算Embedding向量
	embedding, _ := middleware.LLM.Embedding(context.Background(), matchMQ.Description)
	matchMQ.Profile.PartnerEmbedding = pgvector.NewVector(embedding)

	// 3. 将matchProfile保存到数据库
	err := mapper.CreateMatchProfile(&matchMQ.Profile)
	if err != nil {
		zap.S().Error("数据库保存失败", err)
		return
	}
}

// SearchMatch 处理用户搜索匹配的业务逻辑
func SearchMatch(userID int64) (*modelDTO.SearchMatchDTO, error) {
	// 1. 根据userID查询当前用户的匹配信息
	CurrUsrProfile, err := mapper.GetMatchProfileByUserID(userID)
	if err != nil {
		return nil, errors.New("你还没有加入匹配广场，没有你的微信号记录，无法做匹配哦")
	}

	// 2. 根据当前用户的Embedding向量和性别，查询数据库中异性用户的匹配信息，并计算相似度分数
	MatchUsrProfile, score, err := mapper.FindBestMatch(userID, CurrUsrProfile.Gender, &CurrUsrProfile.PartnerEmbedding)

	// 3. 根据RecordID查询对方的算命记录，获取对方的基本信息（如年龄、职业、兴趣等）
	MatchUsrRecord, err := mapper.GetFortuneRecordByID(MatchUsrProfile.FortuneRecordID)

	// 4. 将对方的基本信息和相似度分数封装成SearchMatchDTO对象，返回给前端
	MatchUsrDTO := &modelDTO.SearchMatchDTO{
		RealName:     MatchUsrRecord.RealName,
		WechatID:     MatchUsrProfile.WechatID,
		Gender:       MatchUsrProfile.Gender,
		BirthDate:    MatchUsrRecord.BirthDate,
		CurrentCity:  MatchUsrProfile.City,
		Bazi:         *MatchUsrRecord.Bazi,
		FiveElements: *MatchUsrRecord.FiveElements,
		Similarity:   float32(score),
	}

	return MatchUsrDTO, nil
}

// LeaveMatch 处理用户退出匹配的业务逻辑
func LeaveMatch(userID int64) error {
	// 1. 从数据库中删除当前用户的匹配信息（逻辑删除，设置is_deleted字段为true）
	err := mapper.DeleteMatchProfileByUserID(userID)
	if err != nil {
		return errors.New("退出匹配广场失败，请稍后再试")
	}
	return nil
}
