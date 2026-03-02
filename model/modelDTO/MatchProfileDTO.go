/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package modelDTO

type MatchProfileDTO struct {
	WechatID string `gorm:"type:varchar(100);not null" json:"wechat_id"`
}
