package main

import (
	"gin-web/initialize/cacheRedis"
	"gin-web/initialize/config"
	mysqlDB "gin-web/initialize/mysql"
	"gin-web/initialize/runLog"
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
	//if err = mysqlDB.DB.AutoMigrate(&authcCenter.User{}, &authcCenter.Role{}, &authcCenter.Api{}, &models.OperationLog{}); err != nil {
	//	panic(err)
	//}
	//file.LoadFilePath(config.Conf.APP.FilePath)
}

func main() {
	routers.RouterServer() //http服务
	//err := dist_storage.RenameFileOrDir("../cloud/外文翻译.pdf", "../cloud/翻译.pdf")
	//if err != nil {
	//	fmt.Println("错误:", err)
	//}
}
