package model

import (
	"time"

	"gorm.io/gorm"
)

type RoleMenu struct {
	Base
	RoleId   int64  `gorm:"column:role_id" json:"role_id"`     //角色ID
	MenuId   int64  `gorm:"column:menu_id" json:"menu_id"`     //菜单ID
	CreateBy string `gorm:"column:create_by" json:"create_by"` //创建人
	UpdateBy string `gorm:"column:update_by" json:"update_by"` //更新人
}

func (RoleMenu) TableName() string {
	return "role_menu"
}

func (r *RoleMenu) BeforeCreate(*gorm.DB) (err error) {
	now := time.Now()
	r.CreateAt = &now
	r.UpdateAt = &now
	return
}

func (r *RoleMenu) BeforeSave(*gorm.DB) (err error) {
	now := time.Now()
	r.UpdateAt = &now
	return
}

func (r *RoleMenu) BeforeUpdate(*gorm.DB) (err error) {
	now := time.Now()
	r.UpdateAt = &now
	return
}
