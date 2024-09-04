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
		Router.POST("/change_password", sysUser.ChangePassword)
		Router.GET("/user_info", sysUser.GetUserInfo)
		Router.GET("/deposit_history", sysUser.GetUserDepositList)
		Router.GET("/withdraw_history", sysUser.GetUserWithdrawList)
	}
}
