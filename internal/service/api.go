package service

import (
	"context"
	"time"

	"github.com/galaxy-future/BridgX/internal/constants"

	"github.com/spf13/cast"

	"github.com/galaxy-future/BridgX/internal/logs"

	"github.com/galaxy-future/BridgX/internal/permission"

	"github.com/galaxy-future/BridgX/internal/clients"

	"github.com/galaxy-future/BridgX/internal/model"
)

func CreateApi(ctx context.Context, name, path, method, operator string, status *int8) error {
	api := &model.Api{CreateBy: operator}
	buildApi(api, name, path, method, operator, status)
	return model.Create(api)
}

func UpdateApi(ctx context.Context, id int64, name, path, method, userName string) error {
	var err error
	api := &model.Api{}
	err = model.Get(id, api)
	if err != nil {
		return err
	}
	oldPath, oldMethod := api.Path, api.Method
	tx := clients.WriteDBCli.WithContext(ctx).Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			//update casbin policy
			updateCasbin(id, path, method, oldPath, oldMethod)
		}
	}()
	buildApi(api, name, path, method, userName, api.Status)
	return tx.Save(api).Error
}

func updateCasbin(id int64, path, method, oldPath, oldMethod string) {
	// remove old policy
	permission.E.RemoveFilteredPolicy(1, oldPath, oldMethod)
	roleIds, err := GetRoleIdsByApiIds([]int64{id})
	if err != nil {
		logs.Logger.Errorf("update casbin policy failed, err:[%v], apiId:[%d]", err, id)
		return
	}
	var rules [][]string
	for _, roleId := range roleIds {
		rules = append(rules, []string{cast.ToString(roleId), path, method})
	}
	//add new policy
	permission.E.AddPolicies(rules)
}

func UpdateApiStatus(ctx context.Context, ids []int64, status *int8, operator string) error {
	err := model.Updates(model.Api{}, ids, map[string]interface{}{"status": status, "update_by": operator, "update_at": time.Now()})
	if err != nil {
		return err
	}
	// update casbin policy
	apis, err := GetApisByIds(ids)
	if err != nil {
		return err
	}
	apiMap := apisToMap(apis)
	for _, id := range ids {
		api := apiMap[id]
		// no --> yes
		if *status == constants.FlagYes {
			roleIds, err := GetRoleIdsByApiIds([]int64{id})
			if err != nil {
				logs.Logger.Errorf("update casbin policy failed, err:[%v], apiId:[%d]", err, id)
				return err
			}
			var rules [][]string
			for _, roleId := range roleIds {
				rules = append(rules, []string{cast.ToString(roleId), api.Path, api.Method})
			}
			permission.E.AddPolicies(rules)
		} else { // yes --> no
			permission.E.RemoveFilteredPolicy(1, api.Path, api.Method)
		}
	}
	return nil
}

func buildApi(api *model.Api, name, path, method, userName string, status *int8) {
	api.Name = name
	api.Path = path
	api.Method = method
	api.Status = status
	api.UpdateBy = userName
}

func DeleteApi(ctx context.Context, ids []int64) error {
	var err error
	apis, err := GetApisByIds(ids)
	if err != nil {
		return err
	}
	tx := clients.WriteDBCli.WithContext(ctx).Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			// 3.update casbin policy
			for _, api := range apis {
				permission.E.RemoveFilteredPolicy(1, api.Path, api.Method)
			}
		}
	}()
	// 1.delete menu_api
	err = tx.Delete(model.MenuApi{}, "api_id IN ?", ids).Error
	if err != nil {
		return err
	}
	// 2.delete api
	return tx.Delete(model.Api{}, ids).Error
}

func GetApiById(ctx context.Context, id int64) (*model.Api, error) {
	api := &model.Api{}
	err := model.Get(id, api)
	if err != nil {
		return nil, err
	}
	return api, err
}

func GetApiList(ctx context.Context, apiName, path, method string, status *int8, pageNum, pageSize int) ([]*model.Api, int64, error) {
	res, count, err := model.GetApis(ctx, apiName, path, method, status, pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return res, count, nil
}

func GetApisByIds(ids []int64) ([]*model.Api, error) {
	var apis []*model.Api
	err := model.Gets(ids, &apis)
	if err != nil {
		return nil, err
	}
	return apis, nil
}

func apisToMap(apis []*model.Api) map[int64]*model.Api {
	var apiMap = make(map[int64]*model.Api, 0)
	for _, api := range apis {
		apiMap[api.Id] = api
	}
	return apiMap
}

func GetApisByRoleId(roleId int64) ([]*model.Api, error) {
	menuIds, err := GetMenuIdsByRoleId(roleId)
	if err != nil {
		return nil, err
	}
	return GetApisByMenuIds(menuIds)
}
