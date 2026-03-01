/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package mapper

import (
	"CyberMatchmaker/model"
	global "CyberMatchmaker/pkg"
)

// CreateFortuneRecord 插入算命记录（初始状态为 pending）
func CreateFortuneRecord(record *model.FortuneRecord) error {
	return global.DB.Create(record).Error
}

// UpdateFortuneRecord 消费者 Worker 计算完后，通过这个方法把 LLM 的结果批量写回
func UpdateFortuneRecord(id int64, updates map[string]interface{}) error {
	// updates map 可以包含 status, bazi, description, partner_embedding 等字段
	return global.DB.Model(&model.FortuneRecord{}).
		Where("id = ? AND is_deleted = 0", id).
		Updates(updates).Error
}

// GetFortuneRecordByID 前端轮询查询结果
func GetFortuneRecordByID(id int64, userID int64) (*model.FortuneRecord, error) {
	var record model.FortuneRecord
	// 必须加上 user_id 校验，防止越权查询别人的算命结果
	err := global.DB.Where("id = ? AND user_id = ? AND is_deleted = 0", id, userID).First(&record).Error
	return &record, err
}
