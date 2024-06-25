package router

import (
	v1 "crispy-garbanzo/internal/admin/api/v1"

	"github.com/gin-gonic/gin"
)

type ApiRouter struct {
}

func (s *ApiRouter) InitApiRouter(Router *gin.RouterGroup) {
	var sysUser = v1.ApiGroupSys.SysUserApi
	{
		Router.POST("/register", sysUser.Login)
		Router.POST("/login", sysUser.Login)
	}
}
