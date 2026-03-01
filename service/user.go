/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package service

import (
	config "CyberMatchmaker/config"
	"CyberMatchmaker/mapper"
	"CyberMatchmaker/model"
	global "CyberMatchmaker/pkg"
	"CyberMatchmaker/pkg/jwt"
	"context"
	"errors"
	"fmt"
	"time"
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
func LoginUser(ctx context.Context, username, password string) (string, error) {
	// 1. 查库
	user, err := mapper.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	// 2. 校验账号密码
	if user == nil || user.Password != password {
		return "", errors.New("用户名或密码错误")
	}

	// 3. 生成 JWT Token
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return "", errors.New("生成Token失败")
	}

	// 4. 将 Token 存入 Redis
	key := fmt.Sprintf("%s:%d", config.AppConfig.Jwt.Prefix, user.ID)

	// 安全检查：确保全局对象已初始化
	if global.Redis == nil {
		return "", errors.New("服务器内部错误：Redis未连接")
	}

	// 存入 Redis，过期时间建议与 JWT 过期时间一致
	expire := time.Duration(config.AppConfig.Jwt.Expire) * time.Hour
	err = global.Redis.Set(ctx, key, token, expire).Err()
	if err != nil {
		return "", errors.New("保存登录状态失败")
	}

	return token, nil
}
