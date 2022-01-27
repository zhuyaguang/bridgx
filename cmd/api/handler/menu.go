package handler

import (
	"net/http"
	"strings"

	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/model"

	"github.com/spf13/cast"

	"github.com/galaxy-future/BridgX/cmd/api/helper"
	"github.com/galaxy-future/BridgX/cmd/api/request"
	"github.com/galaxy-future/BridgX/cmd/api/response"
	"github.com/galaxy-future/BridgX/internal/service"
	"github.com/galaxy-future/BridgX/pkg/utils"
	"github.com/gin-gonic/gin"
)

func CreateMenu(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	req := request.CreateMenuRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	err := service.CreateMenu(ctx, req.ParentId, req.Name, req.Icon, req.Path, req.Component, req.Permission,
		user.Name, req.Sort, req.Type, req.Visible, req.OuterLinkFlag, req.ApiIds)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
}

func DeleteMenu(ctx *gin.Context) {
	param := ctx.Param("ids")
	ids := strings.Split(param, ",")
	err := service.DeleteMenu(ctx, utils.ToInt64Slice(ids))
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
}

func UpdateMenu(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	req := request.UpdateMenuRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	err := service.UpdateMenu(ctx, req.Id, req.ParentId, req.Name, req.Icon, req.Path, req.Component, req.Permission,
		user.Name, req.Sort, req.Type, req.Visible, req.OuterLinkFlag, req.ApiIds)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
}

func GetMenuDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	menu, err := service.GetMenuById(ctx, cast.ToInt64(id))
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	apiIds, err := service.GetApiIdsByMenuId(menu.Id)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, response.MenuDetailResponse{
		MenuBase: *helper.BuildMenuBase(menu),
		ApiIds:   apiIds,
	})
}

func GetMenuList(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	req := request.MenuListRequest{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	var menus []*model.Menu
	var total int64
	var err error
	if user.Name == constants.UsernameRoot {
		menus, total, err = service.GetMenuList(ctx, req.Name, req.Visible, req.PageNumber, req.PageSize)
	} else {
		menus, total, err = service.GetMenuListByUserId(ctx, user.UserId, req.Name, req.Visible, req.PageNumber, req.PageSize)
	}
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	pager := response.Pager{
		PageNumber: req.PageNumber,
		PageSize:   req.PageSize,
		Total:      int(total),
	}
	data := &response.MenuListResponse{
		MenuList: helper.ToTree(menus),
		Pager:    pager,
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, data)
}
