/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package model

import (
	"github.com/pgvector/pgvector-go"
	"time"
)

// MatchProfile 映射 match_profile 表
type MatchProfile struct {
	ID              int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          int64  `gorm:"uniqueIndex;not null" json:"user_id"`
	FortuneRecordID int64  `gorm:"not null" json:"fortune_record_id"`
	WechatID        string `gorm:"type:varchar(100);not null" json:"wechat_id"`

	Gender string `gorm:"type:varchar(10);not null;index:idx_match_profile_filters" json:"gender"`
	City   string `gorm:"type:varchar(100);not null;index:idx_match_profile_filters" json:"city"`

	// 引入 pgvector 的 Vector 类型来映射 PostgreSQL 中的 vector(1536)
	PartnerEmbedding pgvector.Vector `gorm:"type:vector(1536)" json:"partner_embedding"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	IsDeleted int8      `gorm:"column:is_deleted;default:0;index:idx_match_profile_filters;index:idx_match_profile_is_deleted" json:"is_deleted"`
}
