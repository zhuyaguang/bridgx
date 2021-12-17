package helper

import (
	"fmt"

	"github.com/galaxy-future/BridgX/internal/constants"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func GetPagerParamFromQuery(c *gin.Context) (pageNumber int, pageSize int) {
	pageNumber = cast.ToInt(c.Query("page_number"))
	pageSize = cast.ToInt(c.Query("page_size"))

	if pageNumber < 1 {
		pageNumber = 1
	}
	if pageSize < 1 || pageSize > constants.DefaultPageSize {
		pageSize = constants.DefaultPageSize
	}
	return pageNumber, pageSize
}

func GetUserToken(c *gin.Context) (string, error) {
	value := c.Value(gf_cluster.HeaderTokenName)
	if value == nil {
		return "", fmt.Errorf("获取登录信息失败")
	}

	return value.(string), nil
}
