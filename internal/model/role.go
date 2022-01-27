package model

import (
	"context"
	"fmt"
	"time"

	"github.com/galaxy-future/BridgX/internal/clients"

	"gorm.io/gorm"
)

type Role struct {
	Base
	Name     string `gorm:"column:name" json:"name"`           //角色名称
	Code     string `gorm:"column:code" json:"code"`           //角色编码
	Status   *int8  `gorm:"column:status" json:"status"`       //状态  0:禁用 1:启用
	Sort     int    `gorm:"column:sort" json:"sort"`           //排序 值越小越靠前
	CreateBy string `gorm:"column:create_by" json:"create_by"` //创建人
	UpdateBy string `gorm:"column:update_by" json:"update_by"` //更新人
}

func (Role) TableName() string {
	return "role"
}

func (r *Role) BeforeCreate(*gorm.DB) (err error) {
	now := time.Now()
	r.CreateAt = &now
	r.UpdateAt = &now
	return
}

func (r *Role) BeforeSave(*gorm.DB) (err error) {
	now := time.Now()
	r.UpdateAt = &now
	return
}

func (r *Role) BeforeUpdate(*gorm.DB) (err error) {
	now := time.Now()
	r.UpdateAt = &now
	return
}

//GetRoles search roles by condition
func GetRoles(ctx context.Context, roleName string, status *int8, pageNum, pageSize int) ([]*Role, int64, error) {
	res := make([]*Role, 0)
	query := clients.ReadDBCli.WithContext(ctx)
	if roleName != "" {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%v%%", roleName))
	}
	if status != nil {
		query = query.Where("status = ?", status)
	}
	count, err := QueryWhere(query, pageNum, pageSize, &res, "sort ASC", true)
	if err != nil {
		return nil, 0, err
	}
	return res, count, nil
}
