package router

import (
	v1 "crispy-garbanzo/internal/app/api/v1"

	"github.com/gin-gonic/gin"
)

type DrawApiRouter struct {
}

func (s *DrawApiRouter) InitApiRouter(Router *gin.RouterGroup) {
	var drawApi = v1.ApiGroupSys.DrawApi
	{
		Router.POST("/draw/make", drawApi.MakeDraw)
		Router.GET("/draw/history", drawApi.GetUserDrawList)
		Router.POST("/draw/join", drawApi.ClaimDrawById)
	}
}
