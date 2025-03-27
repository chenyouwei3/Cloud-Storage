package utils

import (
	"fmt"
	"io"
	"os"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

func NowFormat() string {
	return time.Now().Format(TimeFormat)
}

func MakeFilePart(name, part string) string {
	return fmt.Sprintf("%s.part%s", name, part)
}

func CopyFile(src, dest string) (written int64, err error) {
	//将源文件 src 的内容复制到目标文件 dest。
	//src源文件夹   dest目标地址
	srcF, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcF.Close()

	return WriteFile(dest, srcF)
}

func WriteFile(filename string, reader io.Reader) (written int64, err error) {
	//将 reader 中的数据写入到目标文件 filename 中。
	newFile, err := os.Create(filename)
	if err != nil {
		return 0, err
	}
	defer newFile.Close()

	return io.Copy(newFile, reader)
}

func Must(i interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return i
}
