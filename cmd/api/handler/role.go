package handler

import (
	"net/http"
	"strings"

	"github.com/spf13/cast"

	"github.com/galaxy-future/BridgX/cmd/api/helper"
	"github.com/galaxy-future/BridgX/cmd/api/request"
	"github.com/galaxy-future/BridgX/cmd/api/response"
	"github.com/galaxy-future/BridgX/internal/service"
	"github.com/galaxy-future/BridgX/pkg/utils"
	"github.com/gin-gonic/gin"
)

func CreateRole(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	req := request.CreateRoleRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	err := service.CreateRole(ctx, req.Name, req.Code, user.Name, req.Sort, req.Status, req.MenuIds)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
}

func DeleteRole(ctx *gin.Context) {
	param := ctx.Param("ids")
	ids := strings.Split(param, ",")
	err := service.DeleteRole(ctx, utils.ToInt64Slice(ids))
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
}

func UpdateRole(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	req := request.UpdateRoleRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	err := service.UpdateRole(ctx, req.Id, req.Name, req.Code, user.Name, req.Sort, req.MenuIds)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
}

func UpdateRoleStatus(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	req := request.UpdateRoleStatusRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	err := service.UpdateRoleStatus(ctx, req.Id, req.Status, user.Name)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
}

func GetRoleDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	role, err := service.GetRoleById(ctx, cast.ToInt64(id))
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	menuIds, err := service.GetMenuIdsByRoleId(role.Id)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, response.RoleDetailResponse{
		RoleBase: helper.BuildRoleBase(role),
		MenuIds:  menuIds,
	})
}

func GetRoleList(ctx *gin.Context) {
	req := request.RoleListRequest{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	roles, total, err := service.GetRoleList(ctx, req.Name, req.Status, req.PageNumber, req.PageSize)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	pager := response.Pager{
		PageNumber: req.PageNumber,
		PageSize:   req.PageSize,
		Total:      int(total),
	}
	data := &response.RoleListResponse{
		RoleList: helper.ConvertToRoleList(roles),
		Pager:    pager,
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, data)
}
