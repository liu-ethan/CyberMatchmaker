/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package service

import (
	"CyberMatchmaker/mapper"
	"CyberMatchmaker/middleware"
	"CyberMatchmaker/model"
	"CyberMatchmaker/model/modelDTO"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/pgvector/pgvector-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

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
	// 根据fortuneRecord.Description计算Embedding向量
	embedding, _ := middleware.LLM.Embedding(context.Background(), *fortuneRecord.Description)
	matchProfile.PartnerEmbedding = pgvector.NewVector(embedding)

	// 6. 将matchProfile保存到数据库
	err = mapper.CreateMatchProfile(&matchProfile)
	if err != nil {
		err := errors.New("保存数据库失败，微信号已重复")
		return err
	}

	// 7. 返回成功
	return nil
}
