package asyncRoute

import (
	"fmt"
	"gin-web/internal/initialize/runLog"
	"gin-web/pkg/extendController"
	"gin-web/pkg/task"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	GinTaskQueue = task.NewTaskPool(1, 1024) //gin协程池
)

type WaitConn struct {
	Ctx      *gin.Context
	HttpCode int
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

func (wc *WaitConn) Done() {
	wc.doneOnce.Do(func() {
		close(wc.done)
	})
}

func (wc *WaitConn) Wait() {
	<-wc.done
	//async
}

func (wc *WaitConn) Test(content string, isError bool) {
	if isError {
		runLog.ZapLog.Error(content)
	} else {
		runLog.ZapLog.Error(content)
	}
}

func (wc *WaitConn) setResult(httpResponseCode, CustomCode int, ZhCn, EnUs string, data interface{}, err error) {
	if err != nil {
		wc.result.Message.EnUs = fmt.Sprintf("%s : %v", EnUs, err)
		runLog.ZapLog.Error(ZhCn + "/|^_^|/" + wc.result.Message.EnUs) //定义日志输出格式
	} else {
		wc.result.Message.EnUs = EnUs
		runLog.ZapLog.Info(ZhCn + "/|^_^|/" + EnUs) //定义日志输出格式
	}
	//http状态码  自定义状态码
	wc.HttpCode, wc.result.Code = httpResponseCode, CustomCode
	//定义消息
	wc.result.Message.ZhCn = ZhCn
	wc.result.Data = data
}

// 成功
func (wc *WaitConn) SetSuccessResult(data interface{}) {
	wc.setResult(http.StatusOK, 2000, "请求成功", "success", data, nil)
}

// 客户端自定义错误
func (wc *WaitConn) SetCustomResultByFront(ZhCn, EnUs string, err error) {
	wc.setResult(http.StatusBadRequest, 4000, ZhCn, EnUs, nil, err)
}

// 服务端自定义错误
func (wc *WaitConn) SetCustomResultByBacked(ZhCn, EnUs string, err error) {
	wc.setResult(http.StatusOK, 5000, ZhCn, EnUs, nil, err)
}

// 客户端错误
func (wc *WaitConn) SetResultByFront(customCode int, err error) {
	wc.setResult(http.StatusBadRequest, customCode, extendController.ErrorCodeMap[customCode].EnUs, extendController.ErrorCodeMap[customCode].ZhCn, nil, err)
}

// 服务端错误
func (wc *WaitConn) SetResultByBacked(customCode int, err error) {
	wc.setResult(http.StatusOK, customCode, extendController.ErrorCodeMap[customCode].EnUs, extendController.ErrorCodeMap[customCode].ZhCn, nil, err)
}

// 请求过多
func (wc *WaitConn) SetTooManyResponse(err error) {
	wc.setResult(http.StatusTooManyRequests, 4290, extendController.ErrorCodeMap[4290].EnUs, extendController.ErrorCodeMap[4290].ZhCn, nil, err)
}
