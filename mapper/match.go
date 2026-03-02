/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package mapper

import (
	"CyberMatchmaker/model"
	global "CyberMatchmaker/pkg"
)

func CreateMatchProfile(profile *model.MatchProfile) error {
	return global.DB.Create(profile).Error
}
