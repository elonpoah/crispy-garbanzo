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
		Router.POST("/change/password", sysUser.ChangePassword)
		Router.GET("/user/info", sysUser.GetUserInfo)
		Router.POST("/user/withdraw", sysUser.Withdraw)
		Router.POST("/user/deposit", sysUser.Deposit)
		Router.GET("/deposit/history", sysUser.GetUserDepositList)
		Router.GET("/withdraw/history", sysUser.GetUserWithdrawList)
		Router.GET("/freeSpin/history", sysUser.GetUserFreeSpinList)
		Router.POST("/draw/make", sysUser.MakeDraw)
		Router.GET("/draw/history", sysUser.GetUserDrawList)
		Router.POST("/draw/join", sysUser.ClaimDrawById)
	}
}
