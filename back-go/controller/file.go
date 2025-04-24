package controller

import (
	"gin-web/initialize/runLog"
	"gin-web/models/file"
	"gin-web/utils"
	"gin-web/utils/asyncRoute"
	"gin-web/utils/extendController"
	"net/http"
	"path"
	"fmt"
	"strings"
)

type FileController struct {
	extendController.BaseController
}

type item struct {
	Filename string `json:"filename"`
	IsDir    bool   `json:"isDir"`
	Size     string `json:"size"`
	Date     string `json:"date"`
}

type fileListData struct {
	DiskUsed     uint64    `json:"diskUsed"`
	DiskUsedStr  string    `json:"diskUsedStr"`
	DiskTotal    uint64    `json:"diskTotal"`
	DiskTotalStr string    `json:"diskTotalStr"`
	Total        int       `json:"total"`
	Items        []*item   `json:"items"`
	Tree         *treeNode `json:"tree"`
}

func (f *FileController) Mkdir(wait *asyncRoute.WaitConn) {
	//线程池复用wait
	defer func() {
		wait.Done()
	}()
	var req struct {
		Path string `json:"path"`
	}
	if err := wait.Ctx.BindJSON(&req); err != nil {
		runLog.ZapLog.Info("参数错误,file绑定错误" + err.Error())
		wait.SetResult("参数错误,file绑定错误", err.Error(), http.StatusBadRequest, nil)
		return
	}
	if req.Path == "" {
		runLog.ZapLog.Info("创建路径错误")
		wait.SetResult("创建路径错误", "create path error", http.StatusBadRequest, nil)
		return
	}

	_, err := file.FilePtr.FileInfo.FindDir(req.Path, true)
	if err != nil {
		runLog.ZapLog.Info("查询文件夹错误" + err.Error())
		wait.SetResult(err.Error(), err.Error(), http.StatusBadRequest, nil)
		return
	}
	wait.SetResult("请求成功", "success", http.StatusOK, nil)
}

func (f *FileController) List(wait *asyncRoute.WaitConn) {
	//线程池复用wait
	defer func() {
		wait.Done()
	}()
	var req struct {
		Path string `json:"path"`
	}
	if err := wait.Ctx.BindJSON(&req); err != nil {
		runLog.ZapLog.Info("参数错误,file绑定错误" + err.Error())
		wait.SetResult("参数错误,file绑定错误", err.Error(), http.StatusBadRequest, nil)
		return
	}

	info, err := file.FilePtr.FileInfo.FindDir(req.Path, false)
	if err != nil {
		runLog.ZapLog.Info("查询文件夹错误" + err.Error())
		wait.SetResult(err.Error(), err.Error(), http.StatusBadRequest, nil)
		return
	}
	items := make([]*item, 0, len(info.FileInfos))
	for _, info := range info.FileInfos {
		if info.IsDir || info.FileMD5 != "" {
			_item := &item{
				Filename: info.Name,
				IsDir:    info.IsDir,
				Date:     info.ModeTime,
			}
			if info.IsDir {
				_item.Size = "-"
			} else {
				_item.Size = utils.ConvertBytesString(info.FileSize)
			}
			items = append(items, _item)
		}
	}
	wait.SetResult("请求成功", "success", http.StatusOK, &fileListData{
		DiskTotal:    file.FileDiskTotal,
		DiskTotalStr: utils.ConvertBytesString(file.FileDiskTotal),
		DiskUsed:     file.FilePtr.UsedDisk,
		DiskUsedStr:  utils.ConvertBytesString(file.FilePtr.UsedDisk),
		Total:        len(items),
		Items:        items,
		Tree:         buildTree(items),
	})
}

func (f *FileController) Remove(wait *asyncRoute.WaitConn) {
	defer func() { wait.Done() }()
	var req struct {
		Path     string   `json:"path"`
		Filename []string `json:"filename"`
	}
	if err := wait.Ctx.BindJSON(&req); err != nil {
		runLog.ZapLog.Info("参数错误,file绑定错误" + err.Error())
		wait.SetResult("参数错误,file绑定错误", err.Error(), http.StatusBadRequest, nil)
		return
	}
	if req.Path == "" || len(req.Filename) == 0 {
		runLog.ZapLog.Info("参数错误,file参数为空")
		wait.SetResult("参数错误,file参数为空", "query is nil", http.StatusBadRequest, nil)
		return
	}

	info, err := file.FilePtr.FileInfo.FindDir(req.Path, false)
	if err != nil {
		runLog.ZapLog.Info("查询文件夹错误" + err.Error())
		wait.SetResult(err.Error(), err.Error(), http.StatusBadRequest, nil)
		return
	}

	for _, filename := range req.Filename {
		if err = file.Remove(info, filename); err != nil {
			runLog.ZapLog.Info("删除文件错误" + err.Error())
			wait.SetResult(err.Error(), err.Error(), http.StatusBadRequest, nil)
			return
		}
	}
	wait.SetResult("请求成功", "success", http.StatusOK, nil)

	file.CalUsedDisk()
}

func (f *FileController) Rename(wait *asyncRoute.WaitConn) {
	defer func() { wait.Done() }()
	var req struct {
		Path    string `json:"path"`
		OldName string `json:"oldName"`
		NewName string `json:"newName"`
	}
	if err := wait.Ctx.BindJSON(&req); err != nil {
		runLog.ZapLog.Info("参数错误,file绑定错误" + err.Error())
		wait.SetResult("参数错误,file绑定错误", err.Error(), http.StatusBadRequest, nil)
		return
	}
	if req.Path == "" || req.OldName == "" || req.NewName == "" {
		runLog.ZapLog.Info("参数错误,file参数为空")
		wait.SetResult("参数错误,file参数为空", "query is nil", http.StatusBadRequest, nil)
		return
	}

	if req.OldName == req.NewName {
		return
	}

	if strings.Contains(req.NewName, "/") {
		runLog.ZapLog.Info("文件名不能含有'/'")
		wait.SetResult("文件名不能含有'/'", "query is not  '/'", http.StatusBadRequest, nil)
		return
	}

	dirInfo, err := file.FilePtr.FileInfo.FindDir(req.Path, false)
	if err != nil {
		runLog.ZapLog.Info("查询文件夹错误" + err.Error())
		wait.SetResult(err.Error(), err.Error(), http.StatusBadRequest, nil)
		return
	}

	srcInfo, ok := dirInfo.FileInfos[req.OldName]
	if !ok {
		runLog.ZapLog.Info("文件不存在")
		wait.SetResult("文件不存在", "file does not exist", http.StatusBadRequest, nil)
		return
	}

	if err = file.Copy2(srcInfo, dirInfo, req.NewName); err != nil {
		runLog.ZapLog.Info(err.Error())
		fmt.Println("rename", err)
		wait.SetResult(err.Error(), err.Error(), http.StatusBadRequest, nil)
		return
	}

	// 移除原文件
	_ = file.Remove(dirInfo, req.OldName)

	file.CalUsedDisk()
}

// 移动、复制 文件或文件夹
func (f *FileController) Mvcp(wait *asyncRoute.WaitConn) {
	defer func() { wait.Done() }()
	var req struct {
		Source []string `json:"source"`
		Target string   `json:"target"`
		Move   bool     `json:"move"`
	}
	if err := wait.Ctx.BindJSON(&req); err != nil {
		runLog.ZapLog.Info("参数错误,file绑定错误" + err.Error())
		wait.SetResult("参数错误,file绑定错误", err.Error(), http.StatusBadRequest, nil)
		return
	}
	if len(req.Source) == 0 || req.Target == "" {
		runLog.ZapLog.Info("参数错误,file参数为空")
		wait.SetResult("参数错误,file参数为空", "query is nil", http.StatusBadRequest, nil)
		return
	}

	tarDir, err := file.FilePtr.FileInfo.FindDir(req.Target, false)
	if err != nil {
		runLog.ZapLog.Info("查询文件夹错误" + err.Error())
		wait.SetResult(err.Error(), err.Error(), http.StatusBadRequest, nil)
		return
	}

	for _, source := range req.Source {
		srcPath, srcName := path.Split(source)
		srcDir, err := file.FilePtr.FileInfo.FindDir(srcPath, false)
		if err != nil {
			runLog.ZapLog.Info("查询文件夹错误" + err.Error())
			wait.SetResult(err.Error(), err.Error(), http.StatusBadRequest, nil)
			return
		}

		srcInfo, ok := srcDir.FileInfos[srcName]
		if !ok {
			runLog.ZapLog.Info("文件不存在")
			wait.SetResult("文件不存在", "file does not exist", http.StatusBadRequest, nil)
			return
		}

		// 不能移动到自身或子目录下
		if tarDir.AbsPath == srcDir.AbsPath ||
			strings.Contains(tarDir.AbsPath, srcInfo.AbsPath) {
			runLog.ZapLog.Info("不能拷贝、移动文件夹到自身目录或子目录")
			wait.SetResult("不能拷贝、移动文件夹到自身目录或子目录", "Cannot copy or move folders to their own directory or subdirectories", http.StatusBadRequest, nil)
			return
		}

		if _, ok = tarDir.FileInfos[srcName]; ok {
			runLog.ZapLog.Info("目标目录下已存在同名文件")
			wait.SetResult("目标目录下已存在同名文件", "A file with the same name already exists in the target directory", http.StatusBadRequest, nil)
			return
		}

		if err = file.Copy2(srcInfo, tarDir, srcInfo.Name); err != nil {

			runLog.ZapLog.Info(err.Error())
			wait.SetResult(err.Error(), err.Error(), http.StatusBadRequest, nil)
			return
		}

		if req.Move {
			// 移除原文件
			_ = file.Remove(srcDir, srcName)
		}
	}

	file.CalUsedDisk()
}

// 查询子节点
type treeNode struct {
	Name     string      `json:"name"`
	Children []*treeNode `json:"children,omitempty"`
}

func buildTree(items []*item) *treeNode {
	root := &treeNode{
		Name: "cloud",
	}

	for _, item := range items {
		parts := strings.Split(item.Filename, "\\")
		if len(parts) > 1 {
			insertNode(root, parts[1:], item)
		} else {
			insertNode(root, parts, item)
		}
	}

	return root
}

func insertNode(root *treeNode, parts []string, item *item) {
	if len(parts) == 0 {
		return
	}

	name := parts[0]
	var child *treeNode
	for _, c := range root.Children {
		if c.Name == name {
			child = c
			break
		}
	}

	if child == nil {
		child = &treeNode{
			Name: name,
		}
		root.Children = append(root.Children, child)
	}

	insertNode(child, parts[1:], item)

}
