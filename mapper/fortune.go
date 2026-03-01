/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package mapper

import (
	model "CyberMatchmaker/model"
	global "CyberMatchmaker/pkg"

	"github.com/pgvector/pgvector-go"
)

// JoinMatchSquare 加入匹配广场
func JoinMatchSquare(profile *model.MatchProfile) error {
	return global.DB.Create(profile).Error
}

// LeaveMatchSquare 退出广场（逻辑删除）
func LeaveMatchSquare(userID int64) error {
	return global.DB.Model(&model.MatchProfile{}).
		Where("user_id = ? AND is_deleted = 0", userID).
		Update("is_deleted", 1).Error
}

// SearchPartners 核心黑科技：带标量过滤的向量相似度检索
func SearchPartners(currentUserID int64, embedding pgvector.Vector, targetGender, targetCity string, limit int) ([]model.MatchResult, error) {
	var results []model.MatchResult

	// 1. 构建基础查询：联表 match_profile (mp) 和 fortune_record (fr)
	// 使用 <=> 运算符计算余弦距离，利用 HNSW 索引极速检索
	query := global.DB.Table("match_profile as mp").
		Select("mp.wechat_id, fr.description, mp.partner_embedding <=> ? AS distance", embedding).
		Joins("JOIN fortune_record as fr ON mp.fortune_record_id = fr.id").
		Where("mp.is_deleted = 0 AND fr.is_deleted = 0").
		Where("mp.user_id != ?", currentUserID) // 排除自己，避免搜到自己

	// 2. 动态拼接标量过滤条件（精准命中你建的联合索引）
	if targetGender != "" {
		query = query.Where("mp.gender = ?", targetGender)
	}
	if targetCity != "" {
		query = query.Where("mp.city = ?", targetCity)
	}

	// 3. 执行向量排序（按距离从小到大），并限制条数
	err := query.Order("distance ASC").Limit(limit).Scan(&results).Error

	return results, err
}
