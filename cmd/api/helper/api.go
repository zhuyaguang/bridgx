package helper

import (
	"github.com/galaxy-future/BridgX/cmd/api/response"
	"github.com/galaxy-future/BridgX/internal/model"
)

func ConvertToApiList(apis []*model.Api) []response.ApiDetailResponse {
	res := make([]response.ApiDetailResponse, 0, len(apis))
	for _, api := range apis {
		res = append(res, BuildApiDetailResponse(api))
	}
	return res
}

func BuildApiDetailResponse(api *model.Api) response.ApiDetailResponse {
	return response.ApiDetailResponse{
		Id:       api.Id,
		Name:     api.Name,
		Path:     api.Path,
		Method:   api.Method,
		Status:   api.Status,
		CreateAt: api.CreateAt.Format("2006-01-02 15:04:05"),
		CreateBy: api.CreateBy,
		UpdateAt: api.UpdateAt.Format("2006-01-02 15:04:05"),
		UpdateBy: api.UpdateBy,
	}
}
