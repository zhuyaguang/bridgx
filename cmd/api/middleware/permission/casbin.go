package permission

import (
	"net/http"

	"github.com/spf13/cast"

	"github.com/galaxy-future/BridgX/cmd/api/response"

	"github.com/galaxy-future/BridgX/cmd/api/helper"
	"github.com/galaxy-future/BridgX/internal/permission"
	"github.com/gin-gonic/gin"
)

func CheckPermission() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := helper.GetUserClaims(ctx)
		path := ctx.Request.URL.Path
		method := ctx.Request.Method
		if pass, err := permission.E.Enforce(cast.ToString(user.UserId), path, method); err != nil {
			response.MkResponse(ctx, http.StatusForbidden, "permission denied", nil)
			ctx.Abort()
			return
		} else if !pass {
			response.MkResponse(ctx, http.StatusForbidden, "permission denied", nil)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
