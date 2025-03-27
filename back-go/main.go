package main

import (
	"gin-web/initialize/cacheRedis"
	"gin-web/initialize/config"
	"gin-web/initialize/file"
	mysqlDB "gin-web/initialize/mysql"
	"gin-web/initialize/runLog"
	"gin-web/models"
	"gin-web/models/authcCenter"
	"gin-web/routers"
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
	//设置运行日志
	if err = runLog.InitRunLog(); err != nil {
		panic(err)
	}
	//初始化mysql数据库
	if err = mysqlDB.InitDB(); err != nil {
		panic(err)
	}
	//初始化缓存redis
	if err = cacheRedis.InitRedis(); err != nil {
		panic(err)
	}
	//数据库迁移
	if err = mysqlDB.DB.AutoMigrate(&authcCenter.User{}, &authcCenter.Role{}, &authcCenter.Api{}, &models.OperationLog{}); err != nil {
		panic(err)
	}
	file.InitFilePath(config.Conf.APP.FilePath) //初始化文件存储
	//pprof检测程序性能
	//go func() {
	//	log.Println(http.ListenAndServe("127.0.0.1:6066", nil))
	//}()
}

func main() {
	routers.RouterServer() //http服务

}
