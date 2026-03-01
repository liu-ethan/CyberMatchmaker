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
	"net/http"
	"strconv"

	"CyberMatchmaker/pkg/response"
	"CyberMatchmaker/service"

	"github.com/gin-gonic/gin"
)

// SubmitFortuneRequest 定义提交算命的请求结构
type SubmitFortuneRequest struct {
	RealName    string `json:"real_name" binding:"required"`
	Gender      string `json:"gender" binding:"required"`
	BirthDate   string `json:"birth_date" binding:"required"` // 格式: YYYY-MM-DD
	BirthTime   string `json:"birth_time" binding:"required"` // 格式: HH:mm
	CurrentCity string `json:"current_city" binding:"required"`
}

// SubmitFortune 处理提交算命参数的接口 (POST /fortune/submit)
func SubmitFortune(c *gin.Context) {
	var req SubmitFortuneRequest
	// 1. 参数校验：使用 Gin 的 Binding 功能校验必填项
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数不完整或格式有误")
		return
	}

	// 2. 获取用户信息：从 JWT 中间件注入的上下文中提取 user_id
	// 注意：根据你的 service 签名，这里需要 int64
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未登录")
		return
	}
	userID := userIDVal.(int64)

	// 3. 调用 Service 层处理业务逻辑
	recordID, err := service.SubmitFortune(
		userID,
		req.RealName,
		req.Gender,
		req.BirthDate,
		req.BirthTime,
		req.CurrentCity,
	)

	if err != nil {
		// 如果是业务逻辑错误（如 MQ 投递失败），返回 500
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 4. 按照接口文档返回统一格式
	response.Success(c, gin.H{
		"record_id": recordID,
		"status":    "pending",
	})
}

// GetFortuneResult 处理轮询查询算命结果的接口 (GET /fortune/result)
func GetFortuneResult(c *gin.Context) {
	// 1. 获取 Query 参数中的 record_id
	recordIDStr := c.Query("record_id")
	if recordIDStr == "" {
		response.Error(c, http.StatusBadRequest, "缺少 record_id 参数")
		return
	}

	recordID, err := strconv.ParseInt(recordIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "record_id 格式错误")
		return
	}

	// 2. 获取当前登录用户的 user_id
	// 从 Auth 中间件注入的上下文获取，并断言为 int64 类型
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未登录")
		return
	}
	userID := userIDVal.(int64)

	// 3. 调用 Service 获取结果，传入 recordID 和 userID 进行越权校验
	result, err := service.GetFortuneResult(recordID, userID)
	if err != nil {
		// 如果 Service 返回错误（如记录不存在或 userID 不匹配），返回相应提示
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	// 4. 返回结果
	// 业务逻辑：如果 status 不为 completed，由 service 层控制只返回 status 字段
	response.Success(c, result)
}
