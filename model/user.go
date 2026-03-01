/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package model

import "time"

// User 映射 "user" 表 (注意表名在 PGSQL 中是保留字，GORM 会自动处理首字母大写的复数映射，但建议显式指定 TableName)
type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"` // JSON 序列化时忽略密码，防止接口泄露
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	IsDeleted int8      `gorm:"column:is_deleted;default:0;index" json:"is_deleted"` // 0: 正常, 1: 已删除
}

// TableName 显式指定表名，防止 GORM 将表名解析为 users
func (User) TableName() string {
	return `user`
}
