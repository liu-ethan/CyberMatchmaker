/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package modelDTO

type SubmitFortuneRequestDTO struct {
	RealName    string `json:"real_name" binding:"required"`
	Gender      string `json:"gender" binding:"required"`
	BirthDate   string `json:"birth_date" binding:"required"` // 格式: YYYY-MM-DD
	BirthTime   string `json:"birth_time"`
	CurrentCity string `json:"current_city" binding:"required"`
}
