package service

import (
	"context"
	"time"

	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/spf13/cast"

	"github.com/galaxy-future/BridgX/internal/permission"

	"github.com/galaxy-future/BridgX/internal/clients"

	"github.com/galaxy-future/BridgX/internal/model"
)

func CreateRole(ctx context.Context, name, code, operator string, sort int, status *int8, menuIds []int64) error {
	var err error
	role := &model.Role{CreateBy: operator}
	tx := clients.WriteDBCli.WithContext(ctx).Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			// update casbin policy
			apis, err := GetApisByMenuIds(menuIds)
			if err != nil {
				logs.Logger.Errorf("update casbin policy failed, err:[%v], roleId:[%d]", err, role.Id)
				return
			}
			var rules [][]string
			for _, api := range apis {
				rules = append(rules, []string{cast.ToString(role.Id), api.Path, api.Method})
			}
			// add new policy
			permission.E.AddPolicies(rules)
		}
	}()
	buildRole(role, name, code, operator, sort, status)
	err = tx.Create(role).Error
	if err != nil {
		return err
	}
	roleMenus := buildRoleMenus(role.Id, menuIds, operator)
	return tx.CreateInBatches(roleMenus, len(menuIds)).Error
}

func UpdateRole(ctx context.Context, id int64, name, code, operator string, sort int, menuIds []int64) error {
	var err error
	role := &model.Role{}
	err = model.Get(id, role)
	if err != nil {
		return err
	}
	tx := clients.WriteDBCli.WithContext(ctx).Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			// update casbin policy
			// remove old policy
			permission.E.RemoveFilteredPolicy(0, cast.ToString(id))
			apis, err := GetApisByMenuIds(menuIds)
			if err != nil {
				logs.Logger.Errorf("update casbin policy failed, err:[%v], roleId:[%d]", err, id)
				return
			}
			var rules [][]string
			for _, api := range apis {
				rules = append(rules, []string{cast.ToString(id), api.Path, api.Method})
			}
			// add new policy
			permission.E.AddPolicies(rules)
		}
	}()
	// 1.delete role_menu
	err = tx.Delete(model.RoleMenu{}, "role_id = ?", id).Error
	if err != nil {
		return err
	}
	// 2.insert role_menu
	roleMenus := buildRoleMenus(id, menuIds, operator)
	err = tx.CreateInBatches(roleMenus, len(menuIds)).Error
	if err != nil {
		return err
	}
	// 3.update role
	buildRole(role, name, code, operator, sort, role.Status)
	return tx.Save(role).Error
}

func UpdateRoleStatus(ctx context.Context, ids []int64, status *int8, operator string) error {
	err := model.Updates(model.Role{}, ids, map[string]interface{}{"status": status, "update_by": operator, "update_at": time.Now()})
	if err != nil {
		return err
	}
	// update casbin policy
	for _, id := range ids {
		// no --> yes
		if *status == constants.FlagYes {
			apis, err := GetApisByRoleId(id)
			if err != nil {
				logs.Logger.Errorf("update casbin policy failed, err:[%v], roleIds:[%v]", err, ids)
				return err
			}
			var rules [][]string
			for _, api := range apis {
				rules = append(rules, []string{cast.ToString(id), api.Path, api.Method})
			}
			permission.E.AddPolicies(rules)
		} else { // yes --> no
			permission.E.RemoveFilteredPolicy(0, cast.ToString(id))
		}
	}
	return nil
}

func buildRoleMenus(roleId int64, menuIds []int64, operator string) []model.RoleMenu {
	var roleMenus = make([]model.RoleMenu, 0, len(menuIds))
	for _, menuId := range menuIds {
		roleMenus = append(roleMenus, model.RoleMenu{RoleId: roleId, MenuId: menuId, CreateBy: operator, UpdateBy: operator})
	}
	return roleMenus
}

func buildRole(role *model.Role, name, code, operator string, sort int, status *int8) {
	role.Name = name
	role.Code = code
	role.Sort = sort
	role.Status = status
	role.UpdateBy = operator
}

func DeleteRole(ctx context.Context, ids []int64) error {
	var err error
	tx := clients.WriteDBCli.WithContext(ctx).Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			for _, id := range ids {
				permission.E.DeleteRole(cast.ToString(id))
			}
		}
	}()
	// 1.delete user_role
	err = tx.Delete(model.UserRole{}, "role_id IN ?", ids).Error
	if err != nil {
		return err
	}
	// 2.delete role_menu
	err = tx.Delete(model.RoleMenu{}, "role_id IN ?", ids).Error
	if err != nil {
		return err
	}
	// 3.delete role
	return tx.Delete(model.Role{}, ids).Error
}

func GetRoleById(ctx context.Context, id int64) (*model.Role, error) {
	role := &model.Role{}
	err := model.Get(id, role)
	if err != nil {
		return nil, err
	}
	return role, err
}

func GetRoleList(ctx context.Context, roleName string, status *int8, pageNum, pageSize int) ([]*model.Role, int64, error) {
	res, count, err := model.GetRoles(ctx, roleName, status, pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return res, count, nil
}
