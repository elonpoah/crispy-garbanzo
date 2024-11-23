package router

import (
	v1 "crispy-garbanzo/internal/app/api/v1"

	"github.com/gin-gonic/gin"
)

type BaseApiRouter struct {
}

func (s *BaseApiRouter) InitApiRouter(Router *gin.RouterGroup) {
	var sysUser = v1.ApiGroupSys.SysUserApi
	var session = v1.ApiGroupSys.SessionApi
	var system = v1.ApiGroupSys.SystemApi
	{
		Router.POST("/register", sysUser.Register)
		Router.POST("/login", sysUser.Login)
		Router.POST("/home/data", session.GetHomeRecommand)
		Router.POST("/session/list", session.GetSessionList)
		Router.POST("/session/detail", session.GetSessionById)
		Router.GET("/platform/setting", system.GetPlatformSetting)
		Router.POST("/draw/detail", sysUser.GetDrawById)
	}
}
