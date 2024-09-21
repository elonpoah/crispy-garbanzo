package router

import (
	v1 "crispy-garbanzo/internal/app/api/v1"

	"github.com/gin-gonic/gin"
)

type SessionApiRouter struct {
}

func (s *SessionApiRouter) InitApiRouter(Router *gin.RouterGroup) {
	var sysSession = v1.ApiGroupSys.SessionApi
	{
		Router.POST("/session/ticket", sysSession.BuySessionTicket)
		Router.POST("/session/check", sysSession.CheckSession)
		Router.POST("/game/history", sysSession.GetGameHistory)
		Router.GET("/session/summary", sysSession.GetUserSummary)
		Router.POST("/free/inviteInfo", sysSession.CheckInviteDuty)
		Router.POST("/free/inviteSpin", sysSession.StartInviteSpin)
	}
}
