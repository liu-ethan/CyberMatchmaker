/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package model

import (
	"time"
)

// FortuneRecord 映射 fortune_record 表
type FortuneRecord struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int64     `gorm:"not null;index" json:"user_id"` // 外键关联 User
	RealName    string    `gorm:"type:varchar(50);not null" json:"real_name"`
	Gender      string    `gorm:"type:varchar(10);not null" json:"gender"`
	BirthDate   time.Time `gorm:"type:date;not null" json:"birth_date"`
	BirthTime   string    `gorm:"type:varchar(20)" json:"birth_time"`
	CurrentCity string    `gorm:"type:varchar(100)" json:"current_city"`

	// 大模型计算结果 (使用指针处理 NULL 值，因为 pending 状态下这些字段为空)
	Bazi          *string `gorm:"type:varchar(50)" json:"bazi,omitempty"`
	FiveElements  *string `gorm:"type:varchar(50)" json:"five_elements,omitempty"`
	ZodiacSign    *string `gorm:"type:varchar(10)" json:"zodiac_sign,omitempty"`
	BestCity      *string `gorm:"type:varchar(100)" json:"best_city,omitempty"`
	RecentFortune *string `gorm:"type:text" json:"recent_fortune,omitempty"`
	Description   *string `gorm:"type:text" json:"description,omitempty"`

	Status    string    `gorm:"type:varchar(20);default:'pending';index" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	IsDeleted int8      `gorm:"column:is_deleted;default:0" json:"is_deleted"`
}

// TableName 显式指定表名，防止 GORM 将表名解析为 fortune_record
func (FortuneRecord) TableName() string {
	return `fortune_record`
}

// FortuneTaskMessage 定义了发送到队列的消息格式
type FortuneTaskMessage struct {
	RecordID int64 `json:"id"`
	UserID   int64 `json:"user_id"`
}
