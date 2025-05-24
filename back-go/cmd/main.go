package main

import (
	"gin-web/internal/app/routers"
	"gin-web/internal/initialize/config"
	mysqlDB "gin-web/internal/initialize/mysql"
	"gin-web/internal/initialize/runLog"
	"gin-web/internal/middleware"
	"github.com/gin-gonic/gin"
	_ "net/http/pprof"
)

func init() {
	//初始化配置文件
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	//设置运行模式
	if config.Conf.APP.Mode == "debug" {
		gin.SetMode(gin.DebugMode)
	}
	//初始化mysql数据库
	if err = mysqlDB.InitDB(); err != nil {
		panic(err)
	}
	//设置运行日志
	if err = runLog.InitRunLog(); err != nil {
		panic(err)
	}
}

func main() {
	middleware.InitOperationLogWorker() //操作日志协程
	routers.RouterServer()              //http服务
	defer runLog.ZapLog.Sync()          //运行日志退出
}
