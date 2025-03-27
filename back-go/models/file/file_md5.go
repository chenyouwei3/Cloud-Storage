package file

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
)

func addMD5File(md5 string, info *fileInfo) {
	files, ok := FilePtr.MD5Files[md5]
	//检查当前文件是否存在
	if !ok {
		files = &md5File{
			File: info.AbsPath,
			MD5:  info.FileMD5,
			Size: info.FileSize,
			Ptr:  []string{},
		}
		FilePtr.MD5Files[md5] = files //将文件信息存入全局变量
	}
	files.Ptr = append(files.Ptr, info.FileMD5)
}

// 参数说明(md5)   ptr(地址)
func removeMD5File(md5, ptr string) {
	//删除md5指向
	files, ok := FilePtr.MD5Files[md5]
	//如果文件存在
	if ok {
		idx := -1
		for i := 0; i < len(files.Ptr); i++ {
			if files.Ptr[i] == ptr {
				idx = i
				break
			}
		}
		if idx != -1 {
			files.Ptr = append(files.Ptr[:idx], files.Ptr[idx+1:]...) //删除所有idx的元素
			//如果为空则删除
			if len(files.Ptr) == 0 {
				delete(FilePtr.MD5Files, md5)
			}
		}
	}

}

// 文件 md5 值
func fileMD5(filename string) (string, error) {
	if info, err := os.Stat(filename); err != nil {
		return "", err
	} else if info.IsDir() {
		return "", errors.New(fmt.Sprintf("%s is a dir", filename))
	}

	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
