package dist_storage

import (
	"gin-web/initialize/runLog"
	"gin-web/models/file"
	"gin-web/utils"
	"gin-web/utils/asyncRoute"
	"gin-web/utils/extendController"
	"net/http"
)

type DistStorageController struct {
	extendController.BaseController
}

func (d *DistStorageController) List(wait *asyncRoute.WaitConn) {
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
