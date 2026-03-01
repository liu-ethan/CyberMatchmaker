/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Result 定义了统一的 JSON 返回结构
type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Success 请求成功时的标准返回
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

// Error 请求失败时的标准返回
func Error(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Result{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
