/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// SetupRouter 初始化并配置路由
func SetupRouter() *gin.Engine {
	// 使用 Default() 会默认挂载 Logger 和 Recovery(防 panic 崩溃) 中间件
	r := gin.Default()

	// 预留位置：后续在这里挂载全局跨域中间件
	// r.Use(middleware.Cors())

	// 统一定义 /api/v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 1. 用户模块 (免鉴权)
		userGroup := v1.Group("/user")
		{
			userGroup.POST("/register", func(c *gin.Context) {
				// 伪函数：模拟注册成功
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "伪接口: 注册成功", "data": nil})
			})
			userGroup.POST("/login", func(c *gin.Context) {
				// 伪函数：模拟登录成功并返回 Token
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "伪接口: 登录成功", "data": gin.H{"token": "dummy_jwt_token_123"}})
			})
		}
	}

	return r
}
