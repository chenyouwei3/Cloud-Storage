package routers

import (
	"gin-web/utils/extendController"
	"github.com/gin-gonic/gin"
	"sync"
)

type WaitConn struct {
	Ctx      *gin.Context
	route    string
	result   extendController.Response
	done     chan struct{}
	doneOnce sync.Once
}

func newWaitConn(ctx *gin.Context, route string) *WaitConn {
	return &WaitConn{
		Ctx:   ctx,
		route: route,
		done:  make(chan struct{}),
	}
}

func (this *WaitConn) Done() {
	this.doneOnce.Do(func() {
		close(this.done)
	})
}

func (this *WaitConn) Wait() {
	<-this.done
	//async
}

func (this *WaitConn) SetResult(messageZhCn, messageEnUs string, code int, data interface{}) {
	this.result.Message.ZhCn = messageZhCn
	this.result.Message.EnUs = messageEnUs
	this.result.Code = code
	this.result.Data = data
}

type webTask func()

func (t webTask) Do() {
	t()
}
