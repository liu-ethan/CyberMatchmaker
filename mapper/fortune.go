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

// GetFortuneRecordByID 根据 ID 查询算命记录
func GetFortuneRecordByID(id int64) (*model.FortuneRecord, error) {
	var record model.FortuneRecord
	err := global.DB.First(&record, id).Error
	return &record, err
}

// UpdateFortuneRecord 封装 GORM 的 Save 操作
// 它会根据 record.ID 自动匹配记录并更新所有非零值字段
func UpdateFortuneRecord(record *model.FortuneRecord) error {
	// 使用 pkg.DB 或你定义的全局 DB 变量
	// Save 是一个“全量更新”操作，如果 ID 存在则更新，不存在则插入
	return global.DB.Save(record).Error
}

// GetLatestFortuneRecordByUserID 根据 user_id 查询最新的算命记录（status = completed）
func GetLatestFortuneRecordByUserID(userID int64) (*model.FortuneRecord, error) {
	var record model.FortuneRecord
	err := global.DB.
		Where("user_id = ? and status = ?", userID, "completed").
		Order("created_at desc").
		First(&record).Error
	return &record, err
}
