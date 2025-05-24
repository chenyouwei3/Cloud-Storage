package routers

import (
	"fmt"
	"gin-web/internal/app/controller"
	"gin-web/internal/app/controller/logs"
	"gin-web/internal/initialize/config"
	"gin-web/internal/middleware"
	"gin-web/pkg/asyncRoute"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"runtime"
)

//注释是没有注释代码的参数
//http.Handle()     //http.Handler()
//http.HandleFunc() //	http.HandlerFunc()

func RouterServer() {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery()) //动记录所有 HTTP 请求的详细信息，如请求方法、请求路径、状态码、响应时间等。
	/*---------------------- */
	if config.Conf.APP.StaticFS && runtime.GOOS == "linux" {
		fmt.Println("testing1")
		r.Use(static.Serve("/", static.LocalFile("/app/cloud-storage/internal/dist", false)))
		r.NoRoute(func(ctx *gin.Context) {
			ctx.File("/app/cloud-storage/internal/dist" + "/index.html")
		})
	} else {
		r.Use(static.Serve("/", static.LocalFile("../internal/dist", false)))
		r.NoRoute(func(ctx *gin.Context) {
			ctx.File("../internal/dist" + "/index.html")
		})
	}
	//全局中间件注册
	r.Use(middleware.CorsMiddleware(), middleware.OperationLog())
	r.POST("/login", (&controller.UserController{}).Login) //登录
	user := r.Group("/user")
	{
		user.POST("/insert", (&controller.UserController{}).Insert)                  //增加user
		user.POST("/remove", (&controller.UserController{}).Remove)                  //删除user
		user.POST("/edit", (&controller.UserController{}).Edit)                      //更新user
		user.GET("/getList", (&controller.UserController{}).GetList)                 //查询user
		user.GET("/getUserByRoles", (&controller.UserController{}).GetRolesByUserID) //根据user_id查role

	}
	role := r.Group("/role")
	{
		role.POST("/insert", (&controller.RoleController{}).Insert)                //增加role
		role.POST("/remove", (&controller.RoleController{}).Remove)                //删除role
		role.POST("/edit", (&controller.RoleController{}).Edit)                    //更新role
		role.GET("/getList", (&controller.RoleController{}).GetList)               //查询role
		role.GET("/getRoleByApis", (&controller.RoleController{}).GetApisByRoleID) //根据role_id查api
	}
	api := r.Group("/api")
	{
		api.POST("/insert", (&controller.ApiController{}).Insert)  //增加api
		api.POST("/remove", (&controller.ApiController{}).Remove)  //删除api
		api.POST("/edit", (&controller.ApiController{}).Edit)      //更新api
		api.GET("/getList", (&controller.ApiController{}).GetList) //查询api
	}
	log := r.Group("/log")
	{
		operationLog := log.Group("/operation")
		{
			operationLog.GET("/getList", (&logs.OperationLogController{}).GetList) //查询操作日志
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
		dist.GET("/category", asyncRoute.WarpHandle((&controller.DistStorageController{}).Category))         //文件可视化数据统计
		dist.GET("/onlicePreview", (&controller.DistStorageController{}).OnlinePreview)                      //文件在线预览文件
	}
	// 捕捉不允许的方法
	//r.NoMethod(controller.DefaultController{}.HandleNotFound) //无法匹配的方法
	//r.NoRoute(controller.DefaultController{}.HandleNotFound)  //无法匹配的路由
	if err := r.Run(fmt.Sprintf("0.0.0.0:%d", config.Conf.APP.Port)); err != nil {
		panic(err)
	}
}
