/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package modelDTO

import "time"

type SubmitFortuneRequestDTO struct {
	RealName    string `json:"real_name" binding:"required"`
	Gender      string `json:"gender" binding:"required"`
	BirthDate   string `json:"birth_date" binding:"required"` // 格式: YYYY-MM-DD
	BirthTime   string `json:"birth_time"`
	CurrentCity string `json:"current_city" binding:"required"`
}

type FortuneResponseDTO struct {
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

	Status string `gorm:"type:varchar(20);default:'pending';index" json:"status"`
}
