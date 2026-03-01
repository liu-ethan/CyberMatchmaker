/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package service

import (
	"CyberMatchmaker/mapper"
	"CyberMatchmaker/model"
	"CyberMatchmaker/pkg/jwt"
	"errors"
)

// RegisterUser 处理用户注册
func RegisterUser(username, password string) error {
	// 1. 校验用户名是否已存在
	existingUser, err := mapper.GetUserByUsername(username)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("用户名已存在")
	}

	// 2. 构造数据落库（小项目直接存明文）
	newUser := &model.User{
		Username: username,
		Password: password,
	}
	return mapper.CreateUser(newUser)
}

// LoginUser 处理登录并签发 Token
func LoginUser(username, password string) (string, error) {
	// 1. 查库
	user, err := mapper.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	// 2. 校验账号密码
	if user == nil || user.Password != password {
		return "", errors.New("用户名或密码错误")
	}

	// 3. 生成并返回 JWT Token
	return jwt.GenerateToken(user.ID)
}
