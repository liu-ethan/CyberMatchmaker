/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package route

import (
	"CyberMatchmaker/controller"
	"CyberMatchmaker/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 初始化 Gin 路由配置
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 全局约定：Base URL 为 /api/v1
	apiV1 := r.Group("/api/v1")
	{
		// ============================================================
		// 1. 用户模块 (User Module) - 无需鉴权
		// ============================================================
		userGroup := apiV1.Group("/user")
		{
			// 用户注册: POST /api/v1/user/register
			userGroup.POST("/register", controller.Register)
			// 用户登录: POST /api/v1/user/login
			userGroup.POST("/login", controller.Login)
		}

		// ============================================================
		// 需鉴权接口组 - 使用你提供的 JWTAuth 中间件进行拦截
		// ============================================================
		authRequired := apiV1.Group("/")
		authRequired.Use(middleware.JWTAuth()) // 拦截未登录请求并注入 user_id
		{
			// 2. 算命与异步处理模块 (Fortune Module)
			fortuneGroup := authRequired.Group("/fortune")
			{
				// 提交算命参数: POST /api/v1/fortune/submit
				fortuneGroup.POST("/submit", controller.SubmitFortune)
				// 轮询查询算命结果: GET /api/v1/fortune/result
				fortuneGroup.GET("/result", controller.GetLatestFortuneResult)
			}

			//3. 匹配交友广场模块 (Match Module)
			matchGroup := authRequired.Group("/match")
			{
				// 开启匹配 (加入广场): POST /api/v1/match/join
				matchGroup.POST("/join", controller.JoinMatch)
				// 匹配寻找异性 (核心向量检索): GET /api/v1/match/search
				matchGroup.GET("/search", controller.SearchMatch)
				// 退出匹配广场: POST /api/v1/match/leave
				matchGroup.POST("/leave", controller.LeaveMatch)
			}
		}
	}

	return r
}
