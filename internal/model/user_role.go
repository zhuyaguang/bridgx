package model

import (
	"time"

	"gorm.io/gorm"
)

type UserRole struct {
	Base
	UserId   int64  `gorm:"column:user_id" json:"user_id"`     //用户ID
	RoleId   int64  `gorm:"column:role_id" json:"role_id"`     //角色ID
	CreateBy string `gorm:"column:create_by" json:"create_by"` //创建人
	UpdateBy string `gorm:"column:update_by" json:"update_by"` //更新人
}

func (UserRole) TableName() string {
	return "user_role"
}

func (r *UserRole) BeforeCreate(*gorm.DB) (err error) {
	now := time.Now()
	r.CreateAt = &now
	r.UpdateAt = &now
	return
}

func (r *UserRole) BeforeSave(*gorm.DB) (err error) {
	now := time.Now()
	r.UpdateAt = &now
	return
}

func (r *UserRole) BeforeUpdate(*gorm.DB) (err error) {
	now := time.Now()
	r.UpdateAt = &now
	return
}
