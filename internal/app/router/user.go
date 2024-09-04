package router

import (
	v1 "crispy-garbanzo/internal/app/api/v1"

	"github.com/gin-gonic/gin"
)

type UserApiRouter struct {
}

func (s *UserApiRouter) InitApiRouter(Router *gin.RouterGroup) {
	var sysUser = v1.ApiGroupSys.SysUserApi
	{
		Router.POST("/changePassword", sysUser.ChangePassword)
		Router.GET("/deposit-history", sysUser.GetUserDepositList)
		Router.GET("/withdraw-history", sysUser.GetUserWithdrawList)
	}
}
