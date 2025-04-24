package controller

import (
	"errors"
	"gin-web/initialize/runLog"
	"gin-web/models/dist_storage"
	"gin-web/utils/asyncRoute"
	"gin-web/utils/extendController"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type DistStorageController struct {
	extendController.BaseController
}

type distRequest struct {
	Path string `json:"path"`
}

func (d *DistStorageController) List(wait *asyncRoute.WaitConn) {
	//线程池复用wait
	defer func() { wait.Done() }()
	Path := wait.Ctx.DefaultQuery("path", "../cloud")
	if Path == "" {
		runLog.ZapLog.Info("参数错误-绑定")
		wait.SetResult("参数错误-绑定", "Parameter Error - Binding", http.StatusBadRequest, nil)
		return
	}
	result, err := dist_storage.GetDirInfoWithTree(Path)
	if err != nil {
		runLog.ZapLog.Info("获取数据失败" + err.Error())
		wait.SetResult("获取数据失败", err.Error(), http.StatusBadRequest, nil)
		return
	}
	wait.SetResult("请求成功", "success", http.StatusOK, result)
}

func (d *DistStorageController) Mkdir(wait *asyncRoute.WaitConn) {
	defer func() { wait.Done() }()
	var req distRequest
	if err := wait.Ctx.BindJSON(&req); err != nil {
		runLog.ZapLog.Info("参数绑定错误" + err.Error())
		wait.SetResult("参数绑定错误", err.Error(), http.StatusBadRequest, nil)
		return
	}
	if req.Path == "" {
		runLog.ZapLog.Info("参数错误-参数绑定错误")
		wait.SetResult("参数错误,参数绑定错误", "query is nil", http.StatusBadRequest, nil)
		return
	}
	if err := dist_storage.MakeDir(req.Path); err != nil {
		runLog.ZapLog.Info("创建文件夹失败" + err.Error())
		wait.SetResult("创建文件夹失败", err.Error(), http.StatusBadRequest, nil)
		return
	}
	wait.SetResult("请求成功", "success", http.StatusOK, nil)
}

func (d *DistStorageController) ReName(wait *asyncRoute.WaitConn) {
	defer func() { wait.Done() }()
	var req struct {
		OldPath string `json:"oldPath"`
		NewPath string `json:"newPath"`
	}
	if err := wait.Ctx.BindJSON(&req); err != nil {
		runLog.ZapLog.Info("参数错误,file绑定错误" + err.Error())
		wait.SetResult("参数错误,file绑定错误", err.Error(), http.StatusBadRequest, nil)
		return
	}
	if req.OldPath == "" || req.NewPath == "" {
		runLog.ZapLog.Info("参数错误,file参数为空")
		wait.SetResult("参数错误,file参数为空", "query is nil", http.StatusBadRequest, nil)
		return
	}

	if req.OldPath == req.NewPath {
		return
	}
	if err := dist_storage.RenameFileOrDir(req.OldPath, req.NewPath); err != nil {
		runLog.ZapLog.Info("重命名失败" + err.Error())
		wait.SetResult("重命名失败", "Failed to retrieve data"+err.Error(), http.StatusBadRequest, nil)
		return
	}
	wait.SetResult("请求成功", "success", http.StatusOK, nil)
}

func (d *DistStorageController) Remove(wait *asyncRoute.WaitConn) {
	defer func() { wait.Done() }()
	var req struct {
		DistsPath []string `json:"distsPath"`
	}
	if err := wait.Ctx.BindJSON(&req); err != nil {
		runLog.ZapLog.Info("参数错误,file绑定错误" + err.Error())
		wait.SetResult("参数错误,file绑定错误", err.Error(), http.StatusBadRequest, nil)
		return
	}
	if len(req.DistsPath) == 0 {
		runLog.ZapLog.Info("参数错误,file参数为空")
		wait.SetResult("参数错误,file参数为空", "query is nil", http.StatusBadRequest, nil)
		return
	}
	for _, filename := range req.DistsPath {
		if err := dist_storage.RemoveFileOrDir(filename); err != nil {
			runLog.ZapLog.Info("删除文件错误" + err.Error())
			wait.SetResult("删除文件错误", err.Error(), http.StatusBadRequest, nil)
			return
		}
	}
	wait.SetResult("请求成功", "success", http.StatusOK, nil)
}

func (d *DistStorageController) Copy(wait *asyncRoute.WaitConn) {
	defer func() { wait.Done() }()
	var req distRequest
	if err := wait.Ctx.BindJSON(&req); err != nil {
		runLog.ZapLog.Info("参数错误,file绑定错误" + err.Error())
		wait.SetResult("参数错误,file绑定错误", err.Error(), http.StatusBadRequest, nil)
		return
	}
	if req.Path == "" {
		runLog.ZapLog.Info("参数错误,file参数为空")
		wait.SetResult("参数错误,file参数为空", "query is nil", http.StatusBadRequest, nil)
		return
	}
	//处理路径
	if err := dist_storage.CopyFileOrDir(req.Path, addCopy(req.Path)); err != nil {
		runLog.ZapLog.Info("复制文件/文件夹" + err.Error())
		wait.SetResult("复制文件/文件夹", err.Error(), http.StatusBadRequest, nil)
		return
	}
	wait.SetResult("请求成功", "success", http.StatusOK, nil)
}

func addCopy(path string) string {
	idx := strings.LastIndex(path, "/")
	if idx == -1 {
		return ""
	}
	name := path[idx+1:]
	copyName := name + "(copy)"
	return path[:idx+1] + copyName
}

func (d *DistStorageController) Move(wait *asyncRoute.WaitConn) {
	defer func() { wait.Done() }()
	var req struct {
		OldPath string `json:"oldPath"`
		NewPath string `json:"newPath"`
	}
	if err := wait.Ctx.BindJSON(&req); err != nil {
		runLog.ZapLog.Info("参数错误,file绑定错误" + err.Error())
		wait.SetResult("参数错误,file绑定错误", err.Error(), http.StatusBadRequest, nil)
		return
	}
	if req.OldPath == "" || req.NewPath == "" {
		runLog.ZapLog.Info("参数错误,file参数为空")
		wait.SetResult("参数错误,file参数为空", "query is nil", http.StatusBadRequest, nil)
		return
	}
	if err := dist_storage.MoveFile(req.OldPath, req.NewPath); err != nil {
		runLog.ZapLog.Info("创建文件夹失败" + err.Error())
		wait.SetResult("创建文件夹失败", err.Error(), http.StatusBadRequest, nil)
		return
	}
	wait.SetResult("请求成功", "success", http.StatusOK, nil)
}

func (d *DistStorageController) DropdownMenu(wait *asyncRoute.WaitConn) {
	defer func() { wait.Done() }()
	Path := wait.Ctx.DefaultQuery("path", "../cloud")
	Query := wait.Ctx.DefaultQuery("query", " ")
	if Path == "" || Query == "" {
		runLog.ZapLog.Info("参数错误,menu绑定错误")
		wait.SetResult("参数错误,menu绑定错误", "query is nil", http.StatusBadRequest, nil)
		return
	}
	switch Query {
	case "move":
		res, err := dist_storage.DirMapping(Path)
		if err != nil {
			runLog.ZapLog.Info("获取文件夹列表失败" + err.Error())
			wait.SetResult("获取文件夹列表失败", err.Error(), http.StatusBadRequest, nil)
			return
		}
		wait.SetResult("请求成功", "success", http.StatusOK, res)
	default:
		runLog.ZapLog.Info("参数错误,参数为空")
		wait.SetResult("参数错误,参数为空", "query is nil", http.StatusBadRequest, nil)
		return
	}
}

func (d *DistStorageController) DownLoad(c *gin.Context) {
	var req distRequest
	if err := c.BindJSON(&req); err != nil {
		d.SendParameterErrorResponse(c, err)
		return
	}
	// 确保文件存在
	if _, err := os.Stat(req.Path); os.IsNotExist(err) {
		d.SendNotFoundResponse(c, err)
		return
	}
	// 获取文件名
	filename := filepath.Base(req.Path)
	// 设置响应头
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")
	// 返回文件内容
	d.SendSuccessResponse(c, "success")
}

func (d *DistStorageController) Upload(c *gin.Context) {
	path := c.PostForm("path")
	if path == "" {
		d.SendParameterErrorResponse(c, errors.New("path 参数为空"))
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		d.SendCustomResponse(c, "获取文件失败", "Failed to retrieve file", err)
		return
	}
	savePath := filepath.Join(path, file.Filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		d.SendCustomResponse(c, "保存文件失败", "Failed to save file", err)
		return
	}
	d.SendSuccessResponse(c, "success")
}

func (d *DistStorageController) Category(c *gin.Context) {
	var req distRequest
	if err := c.BindJSON(&req); err != nil {
		d.SendParameterErrorResponse(c, err)
		return
	}
	// 大小,数量,时间
	typeSize, typeNumber, typeTime := make(map[string]float64), make(map[string]int), make(map[string]int) // 大小
	var totalSize float64 = 0
	if err := filepath.Walk(req.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 跳过目录
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(info.Name()))

		sizeMB := float64(info.Size()) / (1024 * 1024) // 字节转 MB，保留小数
		typeSize[ext] += sizeMB                        //统计大小
		typeNumber[ext]++                              //统计种类
		typeTime[info.ModTime().Format("2006-01-02")]++
		totalSize += sizeMB //总大小

		return nil
	}); err != nil {
		d.SendCustomResponse(c, "获取磁盘资源失败", "Failed to retrieve disk resources", err)
		return
	}
	d.SendSuccessResponse(c, struct {
		Total  float64            `json:"total"`
		Size   map[string]float64 `json:"size"`
		Number map[string]int     `json:"number"`
		Time   map[string]int     `json:"time"`
	}{
		totalSize,
		typeSize,
		typeNumber,
		typeTime,
	})
}
