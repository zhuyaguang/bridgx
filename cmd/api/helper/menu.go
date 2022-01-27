package helper

import (
	"github.com/galaxy-future/BridgX/cmd/api/response"
	"github.com/galaxy-future/BridgX/internal/model"
)

func ConvertToMenuList(menus []*model.Menu) []*response.MenuBase {
	res := make([]*response.MenuBase, 0, len(menus))
	for _, menu := range menus {
		res = append(res, BuildMenuBase(menu))
	}
	return res
}

func BuildMenuBase(menu *model.Menu) *response.MenuBase {
	return &response.MenuBase{
		Id:            menu.Id,
		ParentId:      menu.ParentId,
		Name:          menu.Name,
		Icon:          menu.Icon,
		Type:          menu.Type,
		Path:          menu.Path,
		Component:     menu.Component,
		Permission:    menu.Permission,
		Visible:       menu.Visible,
		OuterLinkFlag: menu.OuterLinkFlag,
		Sort:          menu.Sort,
		CreateAt:      menu.CreateAt.Format("2006-01-02 15:04:05"),
		CreateBy:      menu.CreateBy,
		UpdateAt:      menu.UpdateAt.Format("2006-01-02 15:04:05"),
		UpdateBy:      menu.UpdateBy,
	}
}

func ToTree(dbMenus []*model.Menu) []*response.MenuBase {
	menus := ConvertToMenuList(dbMenus)
	mi := make(map[int64]*response.MenuBase)
	for _, item := range menus {
		mi[item.Id] = item
	}

	var list []*response.MenuBase
	for _, item := range menus {
		// root node
		if *item.ParentId == 0 {
			list = append(list, item)
			continue
		}
		if pitem, ok := mi[*item.ParentId]; ok {
			if pitem.Children == nil {
				children := []*response.MenuBase{item}
				pitem.Children = children
				continue
			}
			pitem.Children = append(pitem.Children, item)
		}
	}
	return list
}
