/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package controller

import (
	"CyberMatchmaker/pkg/response"
	"CyberMatchmaker/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SubmitFortune 处理提交算命参数的接口 (POST /fortune/submit)
func SubmitFortune(c *gin.Context) {
	// 1. 从上下文中获取 user_id（JWTAuth 中间件已注入）
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未登录")
		return
	}
	// 2. 传入到Service层
	recordID, err := service.SubmitFortune(c, userID.(int64))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "提交参数失败: "+err.Error())
		return
	}
	// 3. 返回算命记录ID给前端
	response.Success(c, gin.H{"record_id": recordID})
}

// GetFortuneResult 处理轮询查询算命结果的接口 (GET /fortune/result)
func GetFortuneResult(c *gin.Context) {

}
