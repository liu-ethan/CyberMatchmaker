/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package mapper

import (
	"CyberMatchmaker/model"
	global "CyberMatchmaker/pkg"
	"errors"
	"gorm.io/gorm"
)

// GetUserByUsername 根据账号查询用户
func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := global.DB.Where("username = ? AND is_deleted = 0", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没找到记录在业务上往往不算 err，返回 nil 方便上层做判断
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// CreateUser 插入新用户
func CreateUser(user *model.User) error {
	return global.DB.Create(user).Error
}
