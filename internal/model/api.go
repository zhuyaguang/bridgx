package model

import (
	"context"
	"fmt"
	"time"

	"github.com/galaxy-future/BridgX/internal/clients"
	"gorm.io/gorm"
)

type Api struct {
	Base
	Name     string `gorm:"column:name" json:"name"`           //接口名称
	Path     string `gorm:"column:path" json:"path"`           //地址
	Method   string `gorm:"column:method" json:"method"`       //请求方法
	Status   *int8  `gorm:"column:status" json:"status"`       //状态  0:禁用 1:启用
	CreateBy string `gorm:"column:create_by" json:"create_by"` //创建人
	UpdateBy string `gorm:"column:update_by" json:"update_by"` //更新人
}

func (Api) TableName() string {
	return "api"
}

func (r *Api) BeforeCreate(*gorm.DB) (err error) {
	now := time.Now()
	r.CreateAt = &now
	r.UpdateAt = &now
	return
}

func (r *Api) BeforeSave(*gorm.DB) (err error) {
	now := time.Now()
	r.UpdateAt = &now
	return
}

func (r *Api) BeforeUpdate(*gorm.DB) (err error) {
	now := time.Now()
	r.UpdateAt = &now
	return
}

//GetApis search apis by condition
func GetApis(ctx context.Context, apiName, path, method string, status *int8, pageNum, pageSize int) ([]*Api, int64, error) {
	res := make([]*Api, 0)
	query := clients.ReadDBCli.WithContext(ctx)
	if apiName != "" {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%v%%", apiName))
	}
	if path != "" {
		query = query.Where("path LIKE ?", fmt.Sprintf("%%%v%%", path))
	}
	if path != "" {
		query = query.Where("method = ?", fmt.Sprintf("%%%v%%", method))
	}
	if status != nil {
		query = query.Where("status = ?", status)
	}
	count, err := QueryWhere(query, pageNum, pageSize, &res, "id DESC", true)
	if err != nil {
		return nil, 0, err
	}
	return res, count, nil
}
