/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package model

// MatchResult 定义一个专用的 DTO (Data Transfer Object) 承接联表和向量查询的返回结果
type MatchResult struct {
	WechatID    string  `gorm:"column:wechat_id"`
	Description string  `gorm:"column:description"`
	Distance    float64 `gorm:"column:distance"` // pgvector 返回的是距离(越小越近)，要在 Service 层用 1-distance 转成匹配度
}
