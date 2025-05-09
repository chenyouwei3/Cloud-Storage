package routers

import (
	"fmt"
	"gin-web/controller"
	"gin-web/controller/logs"
	"gin-web/initialize/config"
	"gin-web/middleware"
	"gin-web/utils/asyncRoute"
	"os"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

//注释是没有注释代码的参数
//http.Handle()     //http.Handler()
//http.HandleFunc() //	http.HandlerFunc()

func RouterServer() {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery()) //动记录所有 HTTP 请求的详细信息，如请求方法、请求路径、状态码、响应时间等。
	dir, _ := os.Getwd()
	//if config.Conf.APP.StaticFS { // 静态资源浏览
	//	r.StaticFS("/static", gin.Dir(config.Conf.APP.FilePath, true))
	//}
	//if config.Conf.APP.WebIndex != "" {
	//	r.Use(static.Serve("/", static.LocalFile(config.Conf.APP.WebIndex, false)))
	//	r.NoRoute(func(ctx *gin.Context) {
	//		ctx.File(config.Conf.APP.WebIndex + "/index.html")
	//	})
	//}
	/*---------------------- */
	if config.Conf.APP.StaticFS { // 静态资源浏览
		r.StaticFS("/static", gin.Dir(config.Conf.APP.FilePath, true))
	}
	if config.Conf.APP.WebIndex != "" {
		r.Use(static.Serve("/", static.LocalFile(dir+"/initialize/dist", false)))
		r.NoRoute(func(ctx *gin.Context) {
			ctx.File(dir + "/initialize/dist" + "/index.html")
		})
	}
	//全局中间件注册
	r.Use(middleware.CorsMiddleware(), middleware.OperationLog())
	r.POST("/login", (&controller.UserController{}).Login) //登录
	user := r.Group("/user")
	{
		user.POST("/add", (&controller.UserController{}).Insert)       //增加user
		user.DELETE("/deleted", (&controller.UserController{}).Remove) //删除user
		user.PUT("/update", (&controller.UserController{}).Edit)       //更新user
		user.GET("/getAll", (&controller.UserController{}).GetList)    //查询user
	}
	role := r.Group("/role")
	{
		role.POST("/add", (&controller.RoleController{}).Insert)       //增加role
		role.DELETE("/deleted", (&controller.RoleController{}).Remove) //删除role
		role.PUT("/update", (&controller.RoleController{}).Edit)       //更新role
		role.GET("/getAll", (&controller.RoleController{}).GetList)    //查询role
	}
	api := r.Group("/api")
	{
		api.POST("/add", (&controller.ApiController{}).Insert)       //增加api
		api.DELETE("/deleted", (&controller.ApiController{}).Remove) //删除api
		api.PUT("/update", (&controller.ApiController{}).Edit)       //更新api
		api.GET("/getAll", (&controller.ApiController{}).GetList)    //查询api
	}
	log := r.Group("/log")
	{
		operationLog := log.Group("/operation")
		{
			operationLog.GET("/getAll", (&logs.OperationLogController{}).GetList)
		}
	}
	dist := r.Group("/dist") //基于运行的操作系统的文件管理系统
	{
		dist.GET("/list", asyncRoute.WarpHandle((&controller.DistStorageController{}).List))                 //文件列表
		dist.POST("/mkdir", asyncRoute.WarpHandle((&controller.DistStorageController{}).Mkdir))              //创建文件夹
		dist.POST("/rename", asyncRoute.WarpHandle((&controller.DistStorageController{}).ReName))            //重命名(文件/文件夹)
		dist.POST("/remove", asyncRoute.WarpHandle((&controller.DistStorageController{}).Remove))            //删除(文件/文件夹)
		dist.POST("/copy", asyncRoute.WarpHandle((&controller.DistStorageController{}).Copy))                //复制(文件/文件夹)
		dist.POST("/move", asyncRoute.WarpHandle((&controller.DistStorageController{}).Move))                //移动(文件/文件夹)
		dist.GET("/dropdownMenu", asyncRoute.WarpHandle((&controller.DistStorageController{}).DropdownMenu)) //移动(文件/文件夹)下拉框
		dist.POST("/download", (&controller.DistStorageController{}).DownLoad)                               //下载(文件/文件夹)*
		dist.POST("/upload", (&controller.DistStorageController{}).Upload)                                   //上传文件*
		dist.GET("/category", asyncRoute.WarpHandle((&controller.DistStorageController{}).Category))         //可视化数据统计
		dist.GET("/onlicePreview", (&controller.DistStorageController{}).OnlinePreview)                      //在线预览文件
	}
	// 捕捉不允许的方法
	//r.NoMethod(controller.DefaultController{}.HandleNotFound) //无法匹配的方法
	//r.NoRoute(controller.DefaultController{}.HandleNotFound)  //无法匹配的路由
	if err := r.Run(fmt.Sprintf("0.0.0.0:%d", config.Conf.APP.Port)); err != nil {
		panic(err)
	}
}
