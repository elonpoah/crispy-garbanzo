package router

import (
	v1 "crispy-garbanzo/internal/app/api/v1"

	"github.com/gin-gonic/gin"
)

type SessionApiRouter struct {
}

func (s *SessionApiRouter) InitApiRouter(Router *gin.RouterGroup) {
	var sysUser = v1.ApiGroupSys.SessionApi
	{
		Router.POST("/session/ticket", sysUser.BuySessionTicket)
		Router.POST("/session/check", sysUser.CheckSession)
		Router.POST("/game/history", sysUser.GetGameHistory)
	}
}
