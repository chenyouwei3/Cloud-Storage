package file

import (
	"fmt"
	"os"
	"path"
)

func InitFilePath(filepath string) {
	filePath := path.Clean(filepath)       //去除多余的./ ../之类的符号,确保路径格式正确
	_ = os.MkdirAll(filePath, os.ModePerm) //确保 filePath 目录存在，如果不存在则创建
	dirPrefix, _ := path.Split(filePath)   //获取目录前缀
	fmt.Println(dirPrefix)
}
