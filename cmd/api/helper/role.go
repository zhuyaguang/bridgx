package helper

import (
	"github.com/galaxy-future/BridgX/cmd/api/response"
	"github.com/galaxy-future/BridgX/internal/model"
)

func ConvertToRoleList(roles []*model.Role) []response.RoleBase {
	res := make([]response.RoleBase, 0, len(roles))
	for _, role := range roles {
		res = append(res, BuildRoleBase(role))
	}
	return res
}

func BuildRoleBase(role *model.Role) response.RoleBase {
	return response.RoleBase{
		Id:       role.Id,
		Name:     role.Name,
		Code:     role.Code,
		Sort:     role.Sort,
		Status:   role.Status,
		CreateAt: role.CreateAt.Format("2006-01-02 15:04:05"),
		CreateBy: role.CreateBy,
		UpdateAt: role.UpdateAt.Format("2006-01-02 15:04:05"),
		UpdateBy: role.UpdateBy,
	}
}
