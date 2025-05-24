package asyncRoute

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

type webTask func()

func (t webTask) Do() {
	t()
}

// 二次封装routes树
func WarpHandle(fn interface{}) gin.HandlerFunc {
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
	if err := GinTaskQueue.SubmitTask(webTask(func() {
		val.Call([]reflect.Value{reflect.ValueOf(wait)})
	})); err != nil {
		wait.SetTooManyResponse(err)
		wait.Done()
		return
	}
	wait.Wait()
	c.JSON(wait.HttpCode, wait.result)

}

// ctx.FullPath() 会返回当前 gin.Context 对象中的请求路径，
func getCurrentRoute(ctx *gin.Context) string {
	return ctx.FullPath()
}
