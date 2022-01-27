package service

import (
	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/model"
)

func GetMenuIdsByRoleId(roleId int64) ([]int64, error) {
	var roleMenus []model.RoleMenu
	err := model.QueryAll(map[string]interface{}{"role_id": roleId}, &roleMenus, "")
	if err != nil {
		return nil, err
	}
	menuIds := make([]int64, 0, len(roleMenus))
	for _, roleMenu := range roleMenus {
		menuIds = append(menuIds, roleMenu.MenuId)
	}
	return menuIds, nil
}

func GetRoleIdsByMenuIds(menuIds []int64) ([]int64, error) {
	var roleMenus []model.RoleMenu
	err := clients.ReadDBCli.Find(&roleMenus, "menu_id IN ?", menuIds).Error
	if err != nil {
		return nil, err
	}
	roleIds := make([]int64, 0, len(roleMenus))
	for _, roleMenu := range roleMenus {
		roleIds = append(roleIds, roleMenu.RoleId)
	}
	return roleIds, nil
}

func GetRoleMenusByMenuIds(menuIds []int64) ([]model.RoleMenu, error) {
	var roleMenus []model.RoleMenu
	err := clients.ReadDBCli.Find(&roleMenus, "menu_id IN ?", menuIds).Error
	if err != nil {
		return nil, err
	}
	return roleMenus, nil
}

func GetRoleIdsByApiIds(apiIds []int64) ([]int64, error) {
	menuIds, err := GetMenuIdsByApiIds(apiIds)
	if err != nil {
		return nil, err
	}
	roleIds, err := GetRoleIdsByMenuIds(menuIds)
	if err != nil {
		return nil, err
	}
	return roleIds, nil
}

func HasRoleMenuByRoleIdAndMenuIds(roleId int64, menuIds []int64) (bool, error) {
	var cnt int64
	err := clients.ReadDBCli.Model(&model.RoleMenu{}).Where("role_id = ? AND menu_id IN ?", roleId, menuIds).Count(&cnt).Error
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}
