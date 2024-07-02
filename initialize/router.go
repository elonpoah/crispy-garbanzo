package initialize

import (
	"crispy-garbanzo/global"
	"crispy-garbanzo/internal/app/router"
	"fmt"

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
	address := fmt.Sprintf(":%s", global.FPG_CONFIG.Application.Port)
	AdminGroup := Router.Group("api")
	router.RouterGroupSys.InitApiRouter(AdminGroup)
	// swagger；注意：生产环境可以注释掉
	if global.FPG_CONFIG.Application.Mode != "prod" {
		Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
		fmt.Printf("Swagger URL is http://localhost:%s/swagger/index.html\n", address)

	}

	return Router
}
