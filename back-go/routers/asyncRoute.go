package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

// 二次封装routes树
func warpHandle(fn interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		val := reflect.ValueOf(fn) //获取fn的所有信息
		if val.Kind() != reflect.Func {
			panic("value not func")
		}
		//开始异步处理
		transBegin(c, fn)
	}
}

func transBegin(c *gin.Context, fn interface{}) {
	val := reflect.ValueOf(fn) //获取fn的所有信息
	route := getCurrentRoute(c)
	wait := newWaitConn(c, route)
	if err := taskQueue.SubmitTask(webTask(func() {
		fmt.Println("异步处理ing3")
		val.Call(append([]reflect.Value{reflect.ValueOf(wait)}))
	})); err != nil {
		wait.SetResult("访问人数过多", "too many requests", http.StatusTooManyRequests, nil)
		wait.Done()
		return
	}
	wait.Wait()
	c.JSON(wait.result.Code, wait.result)
}

// ctx.FullPath() 会返回当前 gin.Context 对象中的请求路径，
func getCurrentRoute(ctx *gin.Context) string {
	return ctx.FullPath()
}
