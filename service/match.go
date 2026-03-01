/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package service

import (
	"errors"
	"fmt"

	"CyberMatchmaker/mapper"
	"CyberMatchmaker/model"

	"github.com/pgvector/pgvector-go" // 官方 pgvector 支持库
)

// JoinMatchSquare 加入匹配广场
func JoinMatchSquare(userID, fortuneRecordID int64, wechatID, gender, city string) error {
	// 1. 前置校验：提取该用户的算命报告作为大模型的上下文
	record, err := mapper.GetFortuneRecordByID(fortuneRecordID, userID)
	if err != nil || record == nil {
		return errors.New("未找到关联的算命记录或无权访问")
	}
	if record.Status != "completed" {
		return errors.New("算命任务尚未完成，无法生成伴侣画像")
	}

	// 2. 让 LLM 总结“最适合的伴侣画像”
	// 注：直接调用同 package (service) 下的 GenerateFortune 方法
	prompt := fmt.Sprintf(
		"该用户的命理特征如下：八字[%s]，五行[%s]，运势描述[%s]。请你用一句话精炼总结最适合该用户的伴侣画像特征（例如：适合寻找一个八字带水、性格沉稳、有一定经济基础的伴侣）。",
		record.Bazi, record.FiveElements, record.Description,
	)
	profileText, err := GenerateFortune(prompt)
	if err != nil {
		return errors.New("生成伴侣画像失败: " + err.Error())
	}

	// 3. 将生成的“伴侣画像”文本转成 1536 维向量
	embeddingFloats, err := GenerateEmbedding(profileText)
	if err != nil {
		return errors.New("生成特征向量失败: " + err.Error())
	}

	// 4. 组装数据落库
	profile := &model.MatchProfile{
		UserID:           userID,
		FortuneRecordID:  fortuneRecordID,
		WechatID:         wechatID,
		Gender:           gender,
		City:             city,
		PartnerEmbedding: pgvector.NewVector(embeddingFloats), // 将 []float32 转为 pgvector 类型
	}

	return mapper.JoinMatchSquare(profile)
}

// SearchMatch 根据前置条件和向量距离匹配异性
func SearchMatch(userID int64, targetGender, targetCity string, limit int) ([]model.MatchResult, error) {
	// 1. 先查出当前用户自己留在广场上的伴侣画像向量
	// (在实际业务中，你可以给 mapper 加一个 GetMatchProfileByUserID 方法，这里略写逻辑以示意)
	// myProfile, err := mapper.GetMatchProfileByUserID(userID)
	// if err != nil { return nil, err }

	// 假设拿到了 myProfile.PartnerEmbedding
	// 2. 拿着这个向量去数据库进行高维度的余弦相似度检索
	// return mapper.SearchPartners(userID, myProfile.PartnerEmbedding, targetGender, targetCity, limit)

	return nil, nil // 请根据实际情况补全你的 myProfile 获取逻辑
}

// LeaveMatchSquare 退出广场
func LeaveMatchSquare(userID int64) error {
	return mapper.LeaveMatchSquare(userID)
}
