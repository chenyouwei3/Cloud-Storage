package file

import (
	"fmt"
	"gin-web/initialize/runLog"
	"gin-web/utils"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// 在当前目录下创建一个新的子文件/子目录
// 参数说明   name(要创建的子文件/子目录的名称) osDir(ture为目录/false为文件)
func (this *fileInfo) makeChild(name string, isDir bool) (*fileInfo, error) {
	info := &fileInfo{
		Path:     path.Join(this.Path, this.Name),
		Name:     name,
		AbsPath:  path.Join(this.AbsPath, name),
		IsDir:    isDir,
		ModeTime: utils.NowFormat(),
	}
	if isDir {
		info.FileInfos = map[string]*fileInfo{}
		if err := os.MkdirAll(info.AbsPath, os.ModePerm); err != nil { //创建文件夹
			return nil, err
		}
	}
	return info, nil
}

// 查找目录
// 参数说明	filePath(要查找的目录路径)  mkdir(ture不存在的时候创建,false不存在返回错误)

func (this *fileInfo) FindDir(filePath string, mkdir bool) (*fileInfo, error) {
	//去除多余的/
	paths := strings.Split(path.Clean(filePath), "/")
	//遍历路径，查找目录
	info := this
	for i := 1; i < len(paths); i++ {
		dirName := paths[i]
		cInfo, ok := info.FileInfos[dirName] //在当前目录的子目录 map 中查找目标目录
		//查看是否同名
		if ok {
			if !cInfo.IsDir {
				return nil, fmt.Errorf("已存在同名文件")
			}
		} else {
			if mkdir {
				var err error
				if cInfo, err = info.makeChild(dirName, true); err != nil {
					return nil, err
				}
				info.FileInfos[cInfo.Name] = cInfo
			} else {
				return nil, fmt.Errorf("路径不存在")
			}
		}
		info = cInfo
	}
	return info, nil
}

// 清除文件上传过程中生成的临时分片文件
func (this *fileInfo) clearUpload() {
	if this.FileUpload != nil {
		for part := range this.FileUpload.ExistSlice {
			filename := utils.MakeFilePart(this.AbsPath, part)
			_ = os.RemoveAll(filename)
		}
	}
}

// 分片合并为一个大文件
func (this *fileInfo) mergeUpload() {
	if this.FileUpload != nil || this.FileUpload.Total != len(this.FileUpload.ExistSlice) {
		return
	}
	f, err := os.Create(this.AbsPath)
	if err != nil {
		return
	}
	defer f.Close()
	for i := 0; i < this.FileUpload.Total; i++ {
		// 获取当前分片文件的路径
		partFile := utils.MakeFilePart(this.AbsPath, strconv.Itoa(i))
		pf, err := os.Open(partFile)
		if err != nil {
			return // 如果打开分片文件失败，返回
		}
		// 将分片文件的内容写入目标文件
		written, err := io.Copy(f, pf)
		_ = pf.Close() // 关闭分片文件
		if err != nil {
			return // 如果写入文件失败，返回
		}
		// 记录写入的字节数
		runLog.ZapLog.Info("input " + this.AbsPath + " from " + string(partFile) + " written " + string(written))
	}
	this.clearUpload() // 清除文件上传过程中生成的临时分片文件
	if this.FileMD5 != "" {
		removeMD5File(this.FileMD5, this.AbsPath)
	}
	this.FileMD5 = this.FileUpload.Md5
	this.FileSize = this.FileUpload.Size
	this.ModeTime = utils.NowFormat()
	this.FileUpload = nil

	addMD5File(this.FileMD5, this)
	calUsedDisk()
}

// 回调函数,使所有文件夹都执行某个函数
func walk(info *fileInfo, f func(file *fileInfo) error) (err error) {
	// 如果当前文件信息为空，直接返回
	if info == nil {
		return
	}
	// 对当前文件或目录执行回调函数 f
	if err = f(info); err != nil {
		return // 如果回调函数返回错误，则直接返回
	}
	// 如果当前文件是一个目录，则递归遍历该目录下的所有文件
	for _, cInfo := range info.FileInfos {
		if cInfo.IsDir {
			// 递归遍历子目录
			err = walk(cInfo, f)
		} else {
			// 对文件执行回调函数
			err = f(cInfo)
		}
		// 如果遍历过程中发生错误，则返回
		if err != nil {
			return
		}
	}
	// 返回 nil，表示没有错误
	return
}

// 计算出磁盘的使用情况
func calUsedDisk() {
	used := uint64(0)
	walk(FilePtr.FileInfo, func(file *fileInfo) error {
		if !file.IsDir && file.FileSize != 0 {
			used += file.FileSize
		}
		return nil
	})
	FilePtr.UsedDisk = used
}

// 文件删除，
func remove(parent *fileInfo, name string) error {
	info, ok := parent.FileInfos[name]
	if !ok {
		return fmt.Errorf("%s 文件不存在", name)
	}
	delMd5 := map[string]struct{}{} // 待删除的md5文件，源文件
	// 遍历文件
	if err := walk(info, func(file *fileInfo) error {
		if !file.IsDir && file.FileMD5 != "" { //不是文件夹并且md5不为0
			if !saveFileMultiple {
				if md5File_, ok := FilePtr.MD5Files[file.FileMD5]; ok {
					if md5File_.File == file.AbsPath {
						// 此文件为源文件
						delMd5[file.FileMD5] = struct{}{}
					}
				}
			}
			// 删除md5指向
			removeMD5File(file.FileMD5, file.AbsPath)
			// 清理上传的分片
			file.clearUpload()
		}

		return nil
	}); err != nil {
		return err
	}
	//从该文件当中删除其子文件
	delete(parent.FileInfos, info.Name)
	//防止文件多处访问
	if !saveFileMultiple {
		// 文件夹中包含源文件需要拷贝到他处
		for md5 := range delMd5 {
			md5File_, ok := FilePtr.MD5Files[md5]
			if ok {
				//还存在他处引用
				//把a移动到b
				_ = os.Rename(md5File_.File, md5File_.Ptr[0])
				md5File_.File = md5File_.Ptr[0]
			}
		}
	}

	// 删除文件、文件夹
	if err := os.RemoveAll(info.AbsPath); err != nil {
		return err
	}
	return nil

}

// 拷贝到目标目录下
func copy2(src, destParent *fileInfo, destName string) error {
	srcPath := path.Join(src.Path, src.Name)
	return walk(src, func(file *fileInfo) error {
		var fileName string
		var dirInfo, newInfo *fileInfo
		var err error
		if file.AbsPath == src.AbsPath {
			// 自己
			fileName = destName
			dirInfo = destParent
		} else {
			// 当前分支 目录拷贝
			filePath := path.Join(file.Path, file.Name)
			revPath := strings.TrimPrefix(filePath, srcPath+"/")
			revPath = path.Dir(revPath)

			if destParent.Path == "" {
				// 根目录
				revPath = path.Join("cloud", destName, revPath)
			} else {
				revPath = path.Join(destParent.Path, destName, revPath)
			}

			fileName = file.Name

			//fmt.Println(22222, filePath, srcPath, revPath, fileName)
			if dirInfo, err = destParent.FindDir(revPath, true); err != nil {
				return err
			}

		}

		if newInfo, err = dirInfo.makeChild(fileName, file.IsDir); err != nil {
			return err
		}
		//fmt.Println(srcPath, file.Path, file.Name, newInfo.Path, newInfo.Name)
		if !file.IsDir && file.FileMD5 != "" {
			if saveFileMultiple {
				// 真实保存,拷贝文件
				files, _ := FilePtr.MD5Files[file.FileMD5]
				//logger.Info(files, newInfo)
				if _, err := utils.CopyFile(files.Ptr[0], newInfo.AbsPath); err != nil {
					return err
				}
			}

			newInfo.FileSize = file.FileSize
			newInfo.FileMD5 = file.FileMD5
			addMD5File(newInfo.FileMD5, newInfo)
		}
		dirInfo.FileInfos[newInfo.Name] = newInfo

		return nil
	})
}

// 这个函数 loadFilePath 的作用是加载指定目录的文件结构，并构建文件信息映射，包括文件的路径、大小、MD5、修改时间等信息，同时清理上传的分片文件
func loadFilePath(filePath string) {
	// 规范化文件路径，去除多余的 `/` 和 `..`
	filePath = path.Clean(filePath)

	// 确保 `filePath` 目录存在，不存在则创建
	_ = os.MkdirAll(filePath, os.ModePerm)

	// 获取 `filePath` 目录的前缀部分
	dirPrefix, _ := path.Split(filePath)

	// 初始化文件系统结构 `filePtr`
	FilePtr = &fileInfos{
		FileInfo: &fileInfo{
			Path:      "",                     // 根路径为空
			Name:      "cloud",                // 根目录名称
			AbsPath:   filePath,               // 目录的绝对路径
			IsDir:     true,                   // 标记为目录
			FileInfos: map[string]*fileInfo{}, // 存储子文件和子目录
		},
		MD5Files: map[string]*md5File{}, // 存储 MD5 映射的文件信息
	}

	// 遍历 `filePath` 目录下的所有文件和文件夹
	utils.Must(nil, filepath.Walk(filePath, func(absPath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算当前路径相对于 `dirPrefix` 的相对路径
		relativePath := strings.TrimPrefix(absPath, dirPrefix)

		if !f.IsDir() { // 处理文件
			_, filename := path.Split(absPath) // 获取文件名

			if strings.Contains(filename, ".part") {
				// 如果是 `.part` 文件（上传分片文件），则删除
				_ = os.RemoveAll(absPath)
			} else {
				// 计算文件的 MD5 值
				md5, e := fileMD5(absPath)
				if e != nil {
					return e
				}

				// 解析文件所在目录和文件名
				dir, file := path.Split(relativePath)

				// 查找或创建该文件所在的目录信息
				dirInfo, _ := FilePtr.FileInfo.FindDir(dir, true)

				// 创建并存储文件信息
				if fileInfo, err := dirInfo.makeChild(file, false); err != nil {
					return err
				} else {
					fileInfo.FileSize = uint64(f.Size())                     // 记录文件大小
					fileInfo.FileMD5 = md5                                   // 记录文件 MD5
					fileInfo.ModeTime = f.ModTime().Format(utils.TimeFormat) // 记录文件修改时间
					dirInfo.FileInfos[file] = fileInfo                       // 将文件信息加入目录
					addMD5File(md5, fileInfo)                                // 记录 MD5 对应的文件
				}
			}
		} else { // 处理目录
			// 查找或创建目录信息，并记录修改时间
			dirInfo, _ := FilePtr.FileInfo.FindDir(relativePath, true)
			dirInfo.ModeTime = f.ModTime().Format(utils.TimeFormat)
		}

		return nil
	}))

	// 计算磁盘使用情况
	calUsedDisk()
}
