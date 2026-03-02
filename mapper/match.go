/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package mapper

import (
	"CyberMatchmaker/model"
	global "CyberMatchmaker/pkg"
	"errors"

	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

// CreateMatchProfile 将用户的匹配信息保存到数据库中
func CreateMatchProfile(profile *model.MatchProfile) error {
	return global.DB.Create(profile).Error
}

// GetMatchProfileByUserID 根据用户ID查询匹配信息
func GetMatchProfileByUserID(userID int64) (*model.MatchProfile, error) {
	var profile model.MatchProfile
	err := global.DB.Where("user_id = ? and is_deleted = ?", userID, 0).First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("没有找到匹配信息，请先加入匹配广场")
		}
		return nil, err
	}
	return &profile, nil
}

// FindBestMatch 根据当前用户的Embedding向量和性别，查询数据库中异性用户的匹配信息，并计算相似度分数
func FindBestMatch(currUserID int64, curUserGender string, currUserEmbedding *pgvector.Vector) (*model.MatchProfile, float64, error) {
	// 定义一个临时结构体，嵌入原始模型并加上 Score 字段
	type ScanResult struct {
		model.MatchProfile
		Similarity float64 `gorm:"column:similarity"`
	}
	var result ScanResult
	// 1. 确定异性的过滤条件
	targetGender := "male"
	if curUserGender == "male" {
		targetGender = "female"
	}
	// 2. 执行查询
	// Select 部分：计算相似度 (1 - distance)
	// Order 部分：按相似度倒序排列 (DESC)
	err := global.DB.Table("match_profile").
		Select("*, (1 - (partner_embedding <=> ?)) AS similarity", currUserEmbedding).
		// 显式排除当前用户 ID
		Where("id <> ? AND gender = ? AND is_deleted = ?", currUserID, targetGender, 0).
		Order("similarity DESC").
		First(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, errors.New("没有找到匹配的用户，但是已经将你的信息加入匹配池，别人匹配时可能会匹配到你哦！（此系统将保障你的数据安全，如果你不想被匹配到，可以随时删除你的信息）")
		}
		return nil, 0, err
	}
	// 返回匹配的模型对象和计算出的分数
	return &result.MatchProfile, result.Similarity, nil
}
