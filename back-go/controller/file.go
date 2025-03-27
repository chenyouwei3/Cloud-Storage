package controller

import (
	"gin-web/initialize/runLog"
	"gin-web/models/file"
	"gin-web/routers"
	"gin-web/utils/extendController"
	"net/http"
)

type FileController struct {
	extendController.BaseController
}

type fileRequest struct {
	Path string `json:"path"`
}

func (f *FileController) Mkdir(wait routers.WaitConn) {
	//线程池复用wait
	defer func() {
		wait.Done()
	}()
	var req fileRequest
	if err := wait.Ctx.BindJSON(&req); err != nil {
		runLog.ZapLog.Info("参数错误,file绑定错误" + err.Error())
		wait.SetResult("参数错误,file绑定错误", err.Error(), http.StatusBadRequest, nil)
		return
	}
	if req.Path == "" {
		wait.SetResult("创建路径错误", "create path error", http.StatusBadRequest, nil)
		return
	}
	_, err := file.FilePtr.FileInfo.FindDir(req.Path, true)
	if err != nil {
		runLog.ZapLog.Info("参数错误,user绑定错误" + err.Error())
		wait.SetResult(err.Error(), err.Error(), http.StatusBadRequest, nil)
		return
	}
}
