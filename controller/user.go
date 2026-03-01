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

type UserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	// 调用 service 处理逻辑
	err := service.RegisterUser(req.Username, req.Password)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, nil)
}

func Login(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	token, err := service.LoginUser(req.Username, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}
	response.Success(c, gin.H{"token": token})
}
