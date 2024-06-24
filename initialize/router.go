package initialize

import (
	"crispy-garbanzo/internal/admin/router"

	"github.com/gin-gonic/gin"
)

// 初始化总路由
func Routers() *gin.Engine {
	Router := gin.Default()
	Router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    0,
			"data":    "pong",
			"message": "Success",
		})
	})
	PublicGroup := Router.Group("api")
	router.RouterGroupSys.InitApiRouter(PublicGroup)
	return Router
}
