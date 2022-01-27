package model

import (
	"context"
	"fmt"
	"time"

	"github.com/galaxy-future/BridgX/internal/clients"
	"gorm.io/gorm"
)

type Menu struct {
	Base
	ParentId      *int64 `gorm:"column:parent_id" json:"parent_id"`             //父节点ID
	Name          string `gorm:"column:name" json:"name"`                       //菜单名称
	Icon          string `gorm:"column:icon" json:"icon"`                       //图标
	Type          *int8  `gorm:"column:type" json:"type"`                       //菜单类型 0:目录 1:菜单 2:按钮
	Path          string `gorm:"column:path" json:"path"`                       //路径
	Component     string `gorm:"column:component" json:"component"`             //组件
	Permission    string `gorm:"column:permission" json:"permission"`           //权限编码
	Visible       *int8  `gorm:"column:visible" json:"visible"`                 //是否展示  0:否 1:是
	OuterLinkFlag *int8  `gorm:"column:outer_link_flag" json:"outer_link_flag"` //外链标识  0:否 1:是
	Sort          *int   `gorm:"column:sort" json:"sort"`                       //排序 值越小越靠前
	CreateBy      string `gorm:"column:create_by" json:"create_by"`             //创建人
	UpdateBy      string `gorm:"column:update_by" json:"update_by"`             //更新人
}

func (Menu) TableName() string {
	return "menu"
}

func (r *Menu) BeforeCreate(*gorm.DB) (err error) {
	now := time.Now()
	r.CreateAt = &now
	r.UpdateAt = &now
	return
}

func (r *Menu) BeforeSave(*gorm.DB) (err error) {
	now := time.Now()
	r.UpdateAt = &now
	return
}

func (r *Menu) BeforeUpdate(*gorm.DB) (err error) {
	now := time.Now()
	r.UpdateAt = &now
	return
}

//GetMenus search menus by condition
func GetMenus(ctx context.Context, menuName string, visible *int8, menuIds []int64, pageNum, pageSize int) ([]*Menu, int64, error) {
	res := make([]*Menu, 0)
	query := clients.ReadDBCli.WithContext(ctx)
	if menuName != "" {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%v%%", menuName))
	}
	if visible != nil {
		query = query.Where("visible = ?", visible)
	}
	if menuIds != nil && len(menuIds) > 0 {
		query = query.Where("id IN ?", menuIds)
	}
	count, err := QueryWhere(query, pageNum, pageSize, &res, "sort ASC", true)
	if err != nil {
		return nil, 0, err
	}
	return res, count, nil
}

func GetMenusAll(ctx context.Context) ([]*Menu, error) {
	res := make([]*Menu, 0)
	err := clients.ReadDBCli.WithContext(ctx).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func GetMenusByUserId(ctx context.Context, userId int64) ([]*Menu, error) {
	res := make([]*Menu, 0)
	db := clients.ReadDBCli.WithContext(ctx)
	err := db.Raw("SELECT DISTINCT(m.id), m.parent_id FROM role_menu rm, role r, user_role ur, menu m "+
		"WHERE rm.role_id = r.id AND r.id = ur.role_id AND rm.menu_id = m.id AND r.`status` = 1 AND ur.user_id = ?", userId).Scan(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
