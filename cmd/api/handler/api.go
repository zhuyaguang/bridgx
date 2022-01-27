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

func CreateApi(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	req := request.CreateApiRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	err := service.CreateApi(ctx, req.Name, req.Path, req.Method, user.Name, req.Status)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
}

func DeleteApi(ctx *gin.Context) {
	param := ctx.Param("ids")
	ids := strings.Split(param, ",")
	err := service.DeleteApi(ctx, utils.ToInt64Slice(ids))
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
}

func UpdateApi(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	req := request.UpdateApiRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	err := service.UpdateApi(ctx, req.Id, req.Name, req.Path, req.Method, user.Name)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
}

func UpdateApiStatus(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	req := request.UpdateRoleStatusRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	err := service.UpdateApiStatus(ctx, req.Id, req.Status, user.Name)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
}

func GetApiDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	api, err := service.GetApiById(ctx, cast.ToInt64(id))
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, helper.BuildApiDetailResponse(api))
}

func GetApiList(ctx *gin.Context) {
	req := request.ApiListRequest{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	apis, total, err := service.GetApiList(ctx, req.Name, req.Path, req.Method, req.Status, req.PageNumber, req.PageSize)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	pager := response.Pager{
		PageNumber: req.PageNumber,
		PageSize:   req.PageSize,
		Total:      int(total),
	}
	data := &response.ApiListResponse{
		ApiList: helper.ConvertToApiList(apis),
		Pager:   pager,
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, data)
}
