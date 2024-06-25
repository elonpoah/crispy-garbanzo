package initialize

import (
	"crispy-garbanzo/global"
	"crispy-garbanzo/internal/admin/router"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "crispy-garbanzo/docs"

	swaggerfiles "github.com/swaggo/files"
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
	AdminGroup := Router.Group("admin")
	router.RouterGroupSys.InitApiRouter(AdminGroup)
	// swagger；注意：生产环境可以注释掉
	if global.FPG_CONFIG.Application.Mode != "prod" {
		Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	return Router
}
