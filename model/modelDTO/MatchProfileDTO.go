/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package modelDTO

import "time"

type MatchProfileDTO struct {
	WechatID string `gorm:"type:varchar(100);not null" json:"wechat_id"`
}

type SearchMatchDTO struct {
	RealName     string    `form:"real_name" json:"real_name"`
	WechatID     string    `gorm:"type:varchar(100);not null" json:"wechat_id"`
	Gender       string    `gorm:"type:varchar(100);not null" json:"gender"`
	BirthDate    time.Time `json:"birth_date"`
	CurrentCity  string    `gorm:"type:varchar(100);not null" json:"current_city"`
	Bazi         string    `gorm:"type:varchar(100);not null" json:"bazi"`
	FiveElements string    `gorm:"type:varchar(100);not null" json:"five_elements"`
	Similarity   float32   `json:"similarity"`
}
