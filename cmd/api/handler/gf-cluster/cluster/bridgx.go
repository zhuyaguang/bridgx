package cluster

import (
	"net/http"

	"github.com/galaxy-future/BridgX/cmd/api/helper"
	"github.com/galaxy-future/BridgX/internal/service"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
	"github.com/gin-gonic/gin"
)

//HandleListUnusedBridgxCluster lie列出没有被使用的集群
//TODO 使用Brdigx方法
func HandleListUnusedBridgxCluster(c *gin.Context) {
	claims := helper.GetUserClaims(c)
	if claims == nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse("校验身份出错"))
		return
	}
	token, err := helper.GetUserToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gf_cluster.NewFailedResponse(err.Error()))
		return
	}
	pageNumber, pageSize := helper.GetPagerParamFromQuery(c)

	clusters, total, err := service.GetBridgxUnusedCluster(c, claims, pageSize, pageNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gf_cluster.NewFailedResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gf_cluster.NewListUnusedBridgxClusterResponse(clusters, gf_cluster.Pager{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		Total:      total,
	}))
}
