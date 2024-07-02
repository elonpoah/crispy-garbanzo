package main

import (
	"crispy-garbanzo/global"
	"crispy-garbanzo/initialize"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	global.FPG_VP = initialize.Viper()    // 初始化配置
	global.FPG_LOG = initialize.Zap()     // 初始化日志
	global.FPG_I18N = initialize.I18n()   // 初始化国际化
	global.FPG_DB = initialize.Gorm()     // 初始化数据库连接
	global.FPG_REIDS = initialize.Redis() // 初始化Redis连接
}

func main() {
	// 程序结束前关闭数据链接
	if global.FPG_DB != nil {
		initialize.RegisterTables() // 初始化表
		db, _ := global.FPG_DB.DB()
		defer db.Close()
	}
	if global.FPG_REIDS != nil {
		defer global.FPG_REIDS.Close()
	}
	Router := initialize.Routers()
	address := fmt.Sprintf(":%s", global.FPG_CONFIG.Application.Port)
	s := &http.Server{
		Addr:           address,
		Handler:        Router,
		ReadTimeout:    time.Duration(global.FPG_CONFIG.Application.Readtimeout) * time.Second,
		WriteTimeout:   time.Duration(global.FPG_CONFIG.Application.Writertimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Printf("Gin mode is %s, server run success on %s\n", gin.Mode(), address)

	s.ListenAndServe()
}
