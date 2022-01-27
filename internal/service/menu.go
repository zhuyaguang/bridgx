package service

import (
	"context"
	"errors"

	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/permission"

	"github.com/spf13/cast"

	"github.com/galaxy-future/BridgX/internal/clients"

	"github.com/galaxy-future/BridgX/internal/model"
)

func CreateMenu(ctx context.Context, parentId *int64, name, icon, path, component, permission, operator string, sort *int, menuType, visible, outerLinkFlag *int8, apiIds []int64) error {
	var err error
	tx := clients.WriteDBCli.WithContext(ctx).Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	menu := &model.Menu{CreateBy: operator}
	buildMenu(menu, parentId, name, icon, path, component, permission, operator, sort, menuType, visible, outerLinkFlag)
	err = tx.Create(menu).Error
	if err != nil {
		return err
	}
	menuApis := buildMenuApis(menu.Id, apiIds, operator)
	return tx.CreateInBatches(menuApis, len(apiIds)).Error
}

func UpdateMenu(ctx context.Context, id int64, parentId *int64, name, icon, path, component, permissionCode, operator string, sort *int, menuType, visible, outerLinkFlag *int8, apiIds []int64) error {
	var err error
	menu := &model.Menu{}
	err = model.Get(id, menu)
	if err != nil {
		return err
	}
	addRules, removeRules, err := buildUpdCasbinRules(id, apiIds)
	if err != nil {
		return err
	}
	tx := clients.WriteDBCli.WithContext(ctx).Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			// 4.update casbin policy
			// 4.1remove old policy
			permission.E.RemovePolicies(removeRules)
			// 4.2add new policy
			permission.E.AddPolicies(addRules)
		}
	}()
	// 1.delete menu_api
	err = tx.Delete(model.MenuApi{}, "menu_id = ?", id).Error
	if err != nil {
		return err
	}
	// 2.insert menu_api
	menuApis := buildMenuApis(id, apiIds, operator)
	err = tx.CreateInBatches(menuApis, len(apiIds)).Error
	if err != nil {
		return err
	}
	// 3.update menu
	buildMenu(menu, parentId, name, icon, path, component, permissionCode, operator, sort, menuType, visible, outerLinkFlag)
	return tx.Save(menu).Error
}

func buildUpdCasbinRules(id int64, apiIds []int64) ([][]string, [][]string, error) {
	var addRules [][]string
	var removeRules [][]string
	roleIds, err := GetRoleIdsByMenuIds([]int64{id})
	if err != nil {
		logs.Logger.Errorf("update casbin policy failed, err:[%v], menuId:[%d]", err, id)
		return nil, nil, err
	}
	oldApis, err := GetApisByMenuIds([]int64{id})
	if err != nil {
		logs.Logger.Errorf("update casbin policy failed, err:[%v], menuId:[%d]", err, id)
		return nil, nil, err
	}
	// remove old policy
	for _, api := range oldApis {
		menuIds, err := GetMenuIdsByApiIdAndNotInMenuIds(api.Id, []int64{id})
		if err != nil {
			logs.Logger.Errorf("update casbin policy failed, err:[%v], menuId:[%d]", err, id)
			return nil, nil, err
		}
		for _, roleId := range roleIds {
			// 防止一个api被相同角色的不同菜单引用，如果直接删除，可能会误删当前角色在另外一个菜单的相同api权限
			hasOther, err := HasRoleMenuByRoleIdAndMenuIds(roleId, menuIds)
			if err != nil {
				logs.Logger.Errorf("update casbin policy failed, err:[%v], menuId:[%d]", err, id)
				return nil, nil, err
			}
			if hasOther {
				continue
			}
			removeRules = append(removeRules, []string{cast.ToString(roleId), cast.ToString(api.Path), cast.ToString(api.Method)})
		}
	}
	// add new policy
	apis, err := GetApisByIds(apiIds)
	if err != nil {
		logs.Logger.Errorf("update casbin policy failed, err:[%v], menuId:[%d]", err, id)
		return nil, nil, err
	}
	for _, roleId := range roleIds {
		for _, api := range apis {
			addRules = append(addRules, []string{cast.ToString(roleId), api.Path, api.Method})
		}
	}
	return addRules, removeRules, nil
}

func buildMenuApis(menuId int64, apiIds []int64, operator string) []model.MenuApi {
	var menuApis = make([]model.MenuApi, 0, len(apiIds))
	for _, apiId := range apiIds {
		menuApis = append(menuApis, model.MenuApi{MenuId: menuId, ApiId: apiId, CreateBy: operator, UpdateBy: operator})
	}
	return menuApis
}

func buildMenu(menu *model.Menu, parentId *int64, name, icon, path, component, permission, operator string, sort *int, menuType, visible, outerLinkFlag *int8) {
	menu.ParentId = parentId
	menu.Name = name
	menu.Path = path
	menu.Icon = icon
	menu.Type = menuType
	menu.Component = component
	menu.Permission = permission
	menu.Sort = sort
	menu.Visible = visible
	menu.OuterLinkFlag = outerLinkFlag
	menu.UpdateBy = operator
}

func DeleteMenu(ctx context.Context, ids []int64) error {
	var err error
	removeRules, err := buildRemoveCasbinRules(ids)
	tx := clients.WriteDBCli.WithContext(ctx).Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			// 4.update casbin policy
			permission.E.RemovePolicies(removeRules)
		}
	}()
	// 1.delete role_menu
	err = tx.Delete(model.RoleMenu{}, "menu_id IN ?", ids).Error
	if err != nil {
		return err
	}
	// 2.delete menu_api
	err = tx.Delete(model.MenuApi{}, "menu_id IN ?", ids).Error
	if err != nil {
		return err
	}
	// 3.delete menu
	return tx.Delete(model.Menu{}, ids).Error
}

type roleApi struct {
	RoleIds []int64
	Apis    []*model.Api
}

func buildRemoveCasbinRules(ids []int64) ([][]string, error) {
	var removeRules [][]string
	// key:menuId  value:roleApi
	menuMap := make(map[int64]roleApi)
	roleMenus, err := GetRoleMenusByMenuIds(ids)
	if err != nil {
		logs.Logger.Errorf("update casbin policy failed, err:[%v], menuIds:[%d]", err, ids)
		return nil, err
	}
	if len(roleMenus) == 0 {
		return nil, err
	}
	for _, roleMenu := range roleMenus {
		if temp, ok := menuMap[roleMenu.MenuId]; ok {
			temp.RoleIds = append(temp.RoleIds, roleMenu.RoleId)
		} else {
			menuMap[roleMenu.MenuId] = roleApi{RoleIds: []int64{roleMenu.RoleId}}
		}
	}
	for _, menuId := range ids {
		oldApis, err := GetApisByMenuIds([]int64{menuId})
		if err != nil {
			logs.Logger.Errorf("update casbin policy failed, err:[%v], menuIds:[%d]", err, ids)
			return nil, err
		}
		if temp, ok := menuMap[menuId]; ok {
			temp.Apis = oldApis
		} else {
			menuMap[menuId] = roleApi{Apis: oldApis}
		}
	}
	for _, roleApi := range menuMap {
		for _, roleId := range roleApi.RoleIds {
			for _, api := range roleApi.Apis {
				menuIds, err := GetMenuIdsByApiIdAndNotInMenuIds(api.Id, ids)
				if err != nil {
					logs.Logger.Errorf("update casbin policy failed, err:[%v], menuId:[%v]", err, ids)
					return nil, err
				}
				// 防止一个api被相同角色的不同菜单引用，如果直接删除，可能会误删当前角色在另外一个菜单的相同api权限
				hasOther, err := HasRoleMenuByRoleIdAndMenuIds(roleId, menuIds)
				if err != nil {
					logs.Logger.Errorf("update casbin policy failed, err:[%v], menuId:[%v]", err, ids)
					return nil, err
				}
				if hasOther {
					continue
				}
				removeRules = append(removeRules, []string{cast.ToString(roleId), cast.ToString(api.Path), cast.ToString(api.Method)})
			}
		}
	}
	return removeRules, nil
}

func GetMenuById(ctx context.Context, id int64) (*model.Menu, error) {
	menu := &model.Menu{}
	err := model.Get(id, menu)
	if err != nil {
		return nil, err
	}
	return menu, err
}

func GetMenuList(ctx context.Context, menuName string, visible *int8, pageNum, pageSize int) ([]*model.Menu, int64, error) {
	res, count, err := model.GetMenus(ctx, menuName, visible, nil, pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return res, count, nil
}

func GetMenuListByUserId(ctx context.Context, userId int64, menuName string, visible *int8, pageNum, pageSize int) ([]*model.Menu, int64, error) {
	userMenus, err := model.GetMenusByUserId(ctx, userId)
	if err != nil {
		return nil, 0, err
	}
	allMenus, err := model.GetMenusAll(ctx)
	if err != nil {
		return nil, 0, err
	}
	menus := extractUserMenus(userMenus, allMenus)
	if len(menus) == 0 {
		return nil, 0, errors.New("you don't have any resources")
	}
	var menuIds = make([]int64, 0)
	for _, menu := range menus {
		menuIds = append(menuIds, menu.Id)
	}
	return model.GetMenus(ctx, menuName, visible, menuIds, pageNum, pageSize)
}

func extractUserMenus(userMenus []*model.Menu, allMenus []*model.Menu) []*model.Menu {
	if len(userMenus) == 0 || len(allMenus) == 0 {
		return []*model.Menu{}
	}
	menusMap := make(map[int64]*model.Menu)
	var resultMenus = make([]*model.Menu, 0)
	var tempMap = make(map[int64]interface{}, 0)
	for _, menu := range allMenus {
		menusMap[menu.Id] = menu
	}
	for _, dbMenu := range userMenus {
		if *dbMenu.ParentId == 0 && !existsMenu(tempMap, dbMenu.Id) {
			tempMap[dbMenu.Id] = ""
			resultMenus = append(resultMenus, dbMenu)
			continue
		}
		resultMenus = append(resultMenus, dbMenu)
		getAllParent(tempMap, dbMenu, menusMap, &resultMenus)
	}
	return resultMenus
}

func existsMenu(tempMap map[int64]interface{}, menuId int64) bool {
	_, ok := tempMap[menuId]
	return ok
}

func getAllParent(tempMap map[int64]interface{}, dbMenu *model.Menu, menusMap map[int64]*model.Menu, resultMenus *[]*model.Menu) {
	if *dbMenu.ParentId == 0 {
		return
	} else {
		pMenu := menusMap[*dbMenu.ParentId]
		if !existsMenu(tempMap, pMenu.Id) {
			tempMap[pMenu.Id] = ""
			*resultMenus = append(*resultMenus, pMenu)
		}
		getAllParent(tempMap, pMenu, menusMap, resultMenus)
	}
}
