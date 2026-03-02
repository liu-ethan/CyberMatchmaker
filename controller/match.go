/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package controller

import (
	"CyberMatchmaker/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JoinMatch POST /api/v1/match/join
func JoinMatch(c *gin.Context) {
	// 1. 获取userID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}
	// 2. 调用Service层的JoinMatch函数
	err := service.JoinMatch(c, userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 3. 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "成功加入匹配，请等待5秒后查看匹配结果"})
}

// SearchMatch POST /api/v1/match/search
func SearchMatch(c *gin.Context) {
	// 1. 获取userID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}
	// 2. 调用Service层的SearchMatch函数
	result, err := service.SearchMatch(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 3. 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "匹配成功", "data": result})
}
