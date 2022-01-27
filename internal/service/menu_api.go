package service

import (
	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/model"
)

func GetApiIdsByMenuId(menuId int64) ([]int64, error) {
	var menuApis []model.MenuApi
	err := model.QueryAll(map[string]interface{}{"menu_id": menuId}, &menuApis, "")
	if err != nil {
		return nil, err
	}
	apiIds := make([]int64, 0, len(menuApis))
	for _, menuApi := range menuApis {
		apiIds = append(apiIds, menuApi.ApiId)
	}
	return apiIds, nil
}

func GetMenuIdsByApiIds(apiIds []int64) ([]int64, error) {
	var menuApis []model.MenuApi
	err := clients.ReadDBCli.Find(&menuApis, "api_id IN ?", apiIds).Error
	if err != nil {
		return nil, err
	}
	menuIds := make([]int64, 0, len(menuApis))
	for _, menuApi := range menuApis {
		menuIds = append(menuIds, menuApi.MenuId)
	}
	return menuIds, nil
}

func GetMenuIdsByApiIdAndNotInMenuIds(apiId int64, excludeMenuIds []int64) ([]int64, error) {
	var menuApis []model.MenuApi
	err := clients.ReadDBCli.Find(&menuApis, "api_id = ? AND menu_id NOT IN ?", apiId, excludeMenuIds).Error
	if err != nil {
		return nil, err
	}
	menuIds := make([]int64, 0, len(menuApis))
	for _, menuApi := range menuApis {
		menuIds = append(menuIds, menuApi.MenuId)
	}
	return menuIds, nil
}

func GetApiIdsByMenuIds(menuIds []int64) ([]int64, error) {
	var menuApis []model.MenuApi
	err := clients.ReadDBCli.Find(&menuApis, "menu_id IN ?", menuIds).Error
	if err != nil {
		return nil, err
	}
	apiIds := make([]int64, 0, len(menuApis))
	for _, menuApi := range menuApis {
		apiIds = append(apiIds, menuApi.ApiId)
	}
	return apiIds, nil
}

func GetApisByMenuIds(menuIds []int64) ([]*model.Api, error) {
	apiIds, err := GetApiIdsByMenuIds(menuIds)
	if err != nil {
		return nil, err
	}
	apis, err := GetApisByIds(apiIds)
	if err != nil {
		return nil, err
	}
	return apis, nil
}
