package controller

import (
	"fmt"
	dist_storage2 "gin-web/internal/app/models/dist_storage"
	"gin-web/pkg/asyncRoute"
	"gin-web/pkg/extendController"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	var FilterArr []string
	Path := wait.Ctx.Query("path")
	Filter := wait.Ctx.Query("filter")
	if Filter != "" {
		FilterArr = strings.Split(Filter, ",")
	}
	if Path == "" {
		wait.SetResultByFront(4001, nil)
		return
	}
	result, err := dist_storage2.GetDirInfoWithTree(Path, FilterArr)
	if err != nil {
		wait.SetCustomResultByBacked("获取数据失败", "Failed to retrieve data", err)
		return
	}
	wait.SetSuccessResult(result)

}

func (d *DistStorageController) Mkdir(wait *asyncRoute.WaitConn) {
	defer func() { wait.Done() }()
	var req distRequest
	if err := wait.Ctx.BindJSON(&req); err != nil {
		wait.SetResultByFront(4002, err)
		return
	}
	if req.Path == "" {
		wait.SetResultByFront(4001, nil)
		return
	}
	if err := dist_storage2.MakeDir(req.Path); err != nil {
		wait.SetCustomResultByBacked("创建文件夹失败", "Failed to create folder", err)
		return
	}
	wait.SetSuccessResult(nil)
}

func (d *DistStorageController) ReName(wait *asyncRoute.WaitConn) {
	defer func() { wait.Done() }()
	var req struct {
		OldPath string `json:"oldPath"`
		NewPath string `json:"newPath"`
	}
	if err := wait.Ctx.BindJSON(&req); err != nil {
		wait.SetResultByFront(4002, err)
		return
	}
	if req.OldPath == "" || req.NewPath == "" {
		wait.SetResultByFront(4001, nil)
		return
	}

	if req.OldPath == req.NewPath {
		return
	}
	if err := dist_storage2.RenameFileOrDir(req.OldPath, req.NewPath); err != nil {
		wait.SetCustomResultByBacked("重命名失败", "Rename failed", err)
		return
	}
	wait.SetSuccessResult(nil)
}

func (d *DistStorageController) Remove(wait *asyncRoute.WaitConn) {
	defer func() { wait.Done() }()
	var req struct {
		DistsPath []string `json:"distsPath"`
	}
	if err := wait.Ctx.BindJSON(&req); err != nil {
		wait.SetResultByFront(4002, err)
		return
	}
	if len(req.DistsPath) == 0 {
		wait.SetResultByFront(4001, nil)
		return
	}
	for _, filename := range req.DistsPath {
		if err := dist_storage2.RemoveFileOrDir(filename); err != nil {
			wait.SetCustomResultByBacked("删除文件错误", "Delete file/folder error", err)
			return
		}
	}
	wait.SetSuccessResult(nil)
}

func (d *DistStorageController) Copy(wait *asyncRoute.WaitConn) {
	defer func() { wait.Done() }()
	var req distRequest
	if err := wait.Ctx.BindJSON(&req); err != nil {
		wait.SetResultByFront(4002, err)
		return
	}
	if req.Path == "" {
		wait.SetResultByFront(4001, nil)
		return
	}
	//处理路径
	if err := dist_storage2.CopyFileOrDir(req.Path, addCopy(req.Path)); err != nil {
		wait.SetCustomResultByBacked("复制文件/文件夹", "Copy files/folders", err)
		return
	}
	wait.SetSuccessResult(nil)
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
		wait.SetResultByFront(4002, err)
		return
	}
	if req.OldPath == "" || req.NewPath == "" {
		wait.SetResultByFront(4001, nil)
		return
	}
	if err := dist_storage2.MoveFile(req.OldPath, req.NewPath); err != nil {
		wait.SetCustomResultByBacked("创建文件夹失败", "Failed to create folder", err)
		return
	}
	wait.SetSuccessResult(nil)
}

func (d *DistStorageController) DropdownMenu(wait *asyncRoute.WaitConn) {
	defer func() { wait.Done() }()
	Path := wait.Ctx.DefaultQuery("path", "../cloud")
	Query := wait.Ctx.DefaultQuery("query", " ")
	if Path == "" || Query == "" {
		wait.SetResultByFront(4001, nil)
		return
	}
	switch Query {
	case "move":
		res, err := dist_storage2.DirMapping(Path)
		if err != nil {
			wait.SetCustomResultByBacked("获取文件夹列表失败", "Failed to retrieve folder list", err)
			return
		}
		wait.SetSuccessResult(res)
	default:
		wait.SetSuccessResult("query类型不匹配")
	}
}

func (d *DistStorageController) DownLoad(c *gin.Context) {
	var req distRequest
	if err := c.BindJSON(&req); err != nil {
		d.SendParameterErrorResponse(c, 4002, err)
		return
	}

	// 确保文件存在
	if _, err := os.Stat(req.Path); os.IsNotExist(err) {
		d.SendServerErrorResponse(c, 5101, err)
		return
	}

	// 获取文件名和扩展名
	filename := filepath.Base(req.Path)
	ext := strings.ToLower(filepath.Ext(filename))

	// 设置正确的 MIME 类型
	mimeTypes := map[string]string{
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".ppt":  "application/vnd.ms-powerpoint",
		".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".bmp":  "image/bmp",
		".txt":  "text/plain",
		".zip":  "application/zip",
		".rar":  "application/x-rar-compressed",
		".7z":   "application/x-7z-compressed",
	}
	contentType := mimeTypes[ext]
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	// 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", contentType)
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")
	// 直接发送文件
	c.File(req.Path)
	d.SendSuccessResponse(c, "success")
}

func (d *DistStorageController) Upload(c *gin.Context) {
	path := c.PostForm("path")
	if path == "" {
		d.SendParameterErrorResponse(c, 4001, nil)
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		d.SendCustomResponseByFront(c, "获取文件失败", "Failed to retrieve file", err)
		return
	}
	savePath := filepath.Join(path, file.Filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		d.SendCustomResponseByBacked(c, "保存文件失败", "Failed to save file", err)
		return
	}
	d.SendSuccessResponse(c, "success")
}

func (d *DistStorageController) Category(wait *asyncRoute.WaitConn) {
	defer func() { wait.Done() }()
	Path := wait.Ctx.Query("path")
	if Path == "" {
		wait.SetResultByFront(4001, nil)
		return
	}
	var totalSize float64 = 0 // 文件总大小
	// 文件类型的大小/文件类型的数量/ 最近7天的数量和大小
	typeSize, typeNumber, recentlySize, recentlyNumber := make(map[string]float64), make(map[string]int), [7]float64{}, [7]int{}
	now := time.Now()
	startTimes := make([]time.Time, 7)
	for i := 0; i < 7; i++ {
		startTimes[i] = now.AddDate(0, 0, -(6 - i)) // 从6天前到今天
	}
	if err := filepath.Walk(Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() { // 跳过目录
			return nil
		}
		ext := strings.ToLower(filepath.Ext(info.Name()))
		sizeMB := float64(info.Size()) / (1024 * 1024)      // 字节转MB
		sizeMB = math.Round(sizeMB*100) / 100               // 保留两位小数
		typeSize[ext] += sizeMB                             // 统计类型大小
		typeSize[ext] = math.Round(typeSize[ext]*100) / 100 // 再保留两位
		typeNumber[ext]++                                   // 统计种类数量
		totalSize += sizeMB                                 // 总大小
		totalSize = math.Round(totalSize*100) / 100         // 保留两位
		modTime := info.ModTime()
		for i := 0; i < 7; i++ {
			dayStart := time.Date(startTimes[i].Year(), startTimes[i].Month(), startTimes[i].Day(), 0, 0, 0, 0, startTimes[i].Location())
			dayEnd := dayStart.AddDate(0, 0, 1)
			if modTime.After(dayStart) && modTime.Before(dayEnd) {
				recentlyNumber[i]++                                     // 当天的文件数量+1
				recentlySize[i] += sizeMB                               // 当天的文件大小+文件大小
				recentlySize[i] = math.Round(recentlySize[i]*100) / 100 // 保留两位
				break
			}
		}
		return nil
	}); err != nil {
		wait.SetCustomResultByBacked("获取磁盘资源失败", "Failed to retrieve disk resources", err)
		return
	}
	wait.SetSuccessResult(struct {
		Total          float64            `json:"total"`
		Size           map[string]float64 `json:"size"`
		Number         map[string]int     `json:"number"`
		RecentlyNumber [7]int             `json:"recently_number"`
		RecentlySize   [7]float64         `json:"recently_size"`
	}{
		totalSize,
		typeSize,
		typeNumber,
		recentlyNumber,
		recentlySize,
	})
}

func (d *DistStorageController) OnlinePreview(c *gin.Context) {
	path := c.Query("name")
	if path == "" {
		d.SendParameterErrorResponse(c, 4001, nil)
		return
	}

	// 文件路径示例：假设你文件都存在某目录
	fullPath := filepath.Join("/your/file/storage", path)

	// 判断文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.String(http.StatusNotFound, "文件不存在")
		return
	}

	// 获取文件扩展名
	ext := strings.ToLower(filepath.Ext(path))

	// 设置 Content-Type
	mimeType := map[string]string{
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".ppt":  "application/vnd.ms-powerpoint",
		".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".bmp":  "image/bmp",
	}
	contentType := mimeType[ext]
	if contentType == "" {
		contentType = "application/octet-stream" // 如果未知类型，默认用 octet-stream
	}
	c.Header("Content-Type", contentType) // 设置头
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filepath.Base(path)))
	// 读取并返回文件内容
	c.File(fullPath)
}
