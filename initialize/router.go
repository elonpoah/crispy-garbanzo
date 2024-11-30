package initialize

import (
	"crispy-garbanzo/global"
	"crispy-garbanzo/internal/app/router"
	"crispy-garbanzo/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "crispy-garbanzo/docs"

	swaggerfiles "github.com/swaggo/files"
)

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
	{
		router.AppRouterGroup.BaseApiRouter.InitApiRouter(PublicGroup)
	}

	PrivateGroup := Router.Group("api")
	PrivateGroup.Use(middleware.JWTAuth())
	{
		router.AppRouterGroup.UserApiRouter.InitApiRouter(PrivateGroup)
		router.AppRouterGroup.SessionApiRouter.InitApiRouter(PrivateGroup)
		router.AppRouterGroup.DrawApiRouter.InitApiRouter(PrivateGroup)
	}
	// swagger；注意：生产环境可以注释掉
	if global.FPG_CONFIG.Application.Mode != "prod" {
		Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
		fmt.Printf("Swagger URL is http://localhost:%s/swagger/index.html\n", global.FPG_CONFIG.Application.Port)

	}

	return Router
}
