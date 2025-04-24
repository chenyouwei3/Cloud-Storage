package routers

import (
	"fmt"
	"gin-web/controller"
	"gin-web/initialize/config"
	"gin-web/middleware"
	"gin-web/utils/asyncRoute"

	"github.com/gin-gonic/gin"
)

//注释是没有注释代码的参数
//http.Handle()     //http.Handler()
//http.HandleFunc() //	http.HandlerFunc()
//r:=gin.Default()//自带gin.Logger()和gin.Recovery()两个中间件

func RouterServer() {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery()) //动记录所有 HTTP 请求的详细信息，如请求方法、请求路径、状态码、响应时间等。

	//全局中间件注册
	r.Use(
		middleware.CorsMiddleware(),
		middleware.OperationLog("/Gin/V1.0", nil),
	)
	r.POST("/login", (&controller.UserController{}).Login)
	user := r.Group("/user")
	{
		user.POST("/add", (&controller.UserController{}).Add)           //增加user
		user.DELETE("/deleted", (&controller.UserController{}).Deleted) //删除user
		user.PUT("/update", (&controller.UserController{}).Update)      //更新user
		user.GET("/getAll", (&controller.UserController{}).GetAll)      //查询user

	}
	role := r.Group("/role")
	{
		role.POST("/add", (&controller.RoleController{}).Add)           //增加role
		role.DELETE("/deleted", (&controller.RoleController{}).Deleted) //删除role
		role.PUT("/update", (&controller.RoleController{}).Update)      //更新role
		role.GET("/getAll", (&controller.RoleController{}).GetAll)      //查询role
	}
	api := r.Group("/api")
	{
		api.POST("/add", (&controller.ApiController{}).Add)           //增加api
		api.DELETE("/deleted", (&controller.ApiController{}).Deleted) //删除api
		api.PUT("/update", (&controller.ApiController{}).Update)      //更新api
		api.GET("/getAll", (&controller.ApiController{}).GetAll)      //查询api
	}
	dist := r.Group("/dist") //基于运行的操作系统的文件管理系统
	{
		dist.GET("/list", asyncRoute.WarpHandle((&controller.DistStorageController{}).List))      //
		dist.POST("/mkdir", asyncRoute.WarpHandle((&controller.DistStorageController{}).Mkdir))   //
		dist.POST("/rename", asyncRoute.WarpHandle((&controller.DistStorageController{}).ReName)) //
		dist.POST("/remove", asyncRoute.WarpHandle((&controller.DistStorageController{}).Remove)) //
		dist.POST("/copy", asyncRoute.WarpHandle((&controller.DistStorageController{}).Copy))     //
		dist.POST("/move", asyncRoute.WarpHandle((&controller.DistStorageController{}).Move))     //
		dist.GET("/dropdownMenu", asyncRoute.WarpHandle((&controller.DistStorageController{}).DropdownMenu))
		dist.POST("/download", (&controller.DistStorageController{}).DownLoad)
		dist.POST("/upload", (&controller.DistStorageController{}).Upload)
		dist.GET("/category", (&controller.DistStorageController{}).Category)
	}

	//linuxConnect := r.Group("/linux")
	//{
	//
	//}

	// 捕捉不允许的方法
	r.NoMethod(controller.DefaultController{}.HandleNotFound) //无法匹配的方法
	r.NoRoute(controller.DefaultController{}.HandleNotFound)  //无法匹配的路由
	if err := r.Run(fmt.Sprintf("0.0.0.0:%d", config.Conf.APP.Port)); err != nil {
		panic(err)
	}
}
