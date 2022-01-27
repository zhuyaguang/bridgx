package model

import (
	"time"

	"gorm.io/gorm"
)

type MenuApi struct {
	Base
	MenuId   int64  `gorm:"column:menu_id" json:"menu_id"`     //菜单ID
	ApiId    int64  `gorm:"column:api_id" json:"api_id"`       //apiID
	CreateBy string `gorm:"column:create_by" json:"create_by"` //创建人
	UpdateBy string `gorm:"column:update_by" json:"update_by"` //更新人
}

func (MenuApi) TableName() string {
	return "menu_api"
}

func (r *MenuApi) BeforeCreate(*gorm.DB) (err error) {
	now := time.Now()
	r.UpdateAt = &now
	return
}

func (r *MenuApi) BeforeUpdate(*gorm.DB) (err error) {
	now := time.Now()
	r.UpdateAt = &now
	return
}
