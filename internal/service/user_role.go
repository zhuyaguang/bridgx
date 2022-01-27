package service

import (
	"errors"

	"github.com/galaxy-future/BridgX/internal/model"
	"gorm.io/gorm"
)

func GetRoleIdsByUserId(userId int64) ([]int64, error) {
	var userRoles []model.UserRole
	err := model.QueryAll(map[string]interface{}{"user_id": userId}, &userRoles, "")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []int64{}, nil
		}
		return nil, err
	}
	roleIds := make([]int64, 0, len(userRoles))
	for _, userRole := range userRoles {
		roleIds = append(roleIds, userRole.RoleId)
	}
	return roleIds, nil
}
