/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package middleware

import (
	"CyberMatchmaker/config"
	"CyberMatchmaker/pkg/jwt"
	"CyberMatchmaker/pkg/response"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	// 假设你的 redis 客户端在这个包下，请根据实际路径修改
	global "CyberMatchmaker/pkg"
)

// var REDIS_KEY_PREFIX = config.AppConfig.Jwt.Prefix

// JWTAuth 鉴权中间件入口
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		// 1. 解析 Header 获取 userID 和 原始 token
		userID, tokenStr, err := parseAuthHeader(authHeader)
		if err != nil {
			abortWithUnauthorized(c, err.Error())
			return
		}

		// 2. Redis 二次校验：确保 Token 与缓存中一致
		// Key 格式: CyberMatchmaker:user:userID
		redisKey := fmt.Sprintf("%s:%d", config.AppConfig.Jwt.Prefix, userID)

		// 从 Redis 中获取存储的 Token
		if global.Redis == nil {
			abortWithUnauthorized(c, "登录已过期，请重新登录")
			return
		}

		storedToken, err := global.Redis.Get(context.Background(), redisKey).Result()

		if err != nil {
			// 如果 Redis 找不到对应的 Key，说明 Session 已过期或用户已登出
			abortWithUnauthorized(c, "登录已过期，请重新登录")
			return
		}

		// 比较传入的 Token 与 Redis 中的 Token 是否一致
		if storedToken != tokenStr {
			abortWithUnauthorized(c, "账号在其他地方登录或 Token 已失效")
			return
		}

		// 3. 将解析出的 user_id 存入 Gin 的 Context 中
		c.Set("user_id", userID)
		c.Next()
	}
}

// 辅助函数：统一处理 401 返回
func abortWithUnauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, response.Result{
		Code: http.StatusUnauthorized,
		Msg:  msg,
	})
	c.Abort()
}

// parseAuthHeader 修改：增加返回 tokenStr 字符串，用于后续 Redis 比对
func parseAuthHeader(header string) (int64, string, error) {
	if header == "" {
		return 0, "", errors.New("请求未携带 Token，无权访问")
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, "", errors.New("Token 格式错误，请使用 Bearer 格式")
	}

	tokenStr := parts[1]
	claims, err := jwt.ParseToken(tokenStr)
	if err != nil {
		return 0, "", errors.New("Token 无效或已过期")
	}

	return claims.UserID, tokenStr, nil
}
