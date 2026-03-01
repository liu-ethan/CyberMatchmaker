/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package middleware

import (
	"CyberMatchmaker/pkg/jwt"
	"CyberMatchmaker/pkg/response"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth 鉴权中间件入口
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		// 调用拆分出去的解析逻辑
		userID, err := parseAuthHeader(authHeader)
		if err != nil {
			// 统一使用你的 response.Result 结构体返回
			c.JSON(http.StatusUnauthorized, response.Result{
				Code: http.StatusUnauthorized,
				Msg:  err.Error(),
			})
			c.Abort() // 拦截请求
			return
		}

		// 将解析出的 user_id 存入 Gin 的 Context 中
		c.Set("user_id", userID)

		// 放行请求
		c.Next()
	}
}

// parseAuthHeader 拆分出来的核心业务逻辑：校验 Header 格式并解析提取 UserID
func parseAuthHeader(header string) (int64, error) {
	if header == "" {
		return 0, errors.New("请求未携带 Token，无权访问")
	}

	// 按空格分割，验证 Bearer 格式
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, errors.New("Token 格式错误，请使用 Bearer 格式")
	}

	// 解析具体的 Token 字符串
	claims, err := jwt.ParseToken(parts[1])
	if err != nil {
		return 0, errors.New("Token 无效或已过期")
	}

	return claims.UserID, nil
}
