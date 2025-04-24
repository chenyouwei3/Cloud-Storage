package utils

import (
	"fmt"
	"io"
	"math"
	"os"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

var (
	KB = uint64(math.Pow(2, 10))
	MB = uint64(math.Pow(2, 20))
	GB = uint64(math.Pow(2, 30))
	TB = uint64(math.Pow(2, 40))
)

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
		fmt.Println(err.Error(), "xixix")
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

func ConvertBytesString(b uint64) string {
	cf, s := ConvertBytes(b)
	return fmt.Sprintf("%.1f%s", cf, s)
}

func CelsiusToFahrenheit(c int) int {
	return c*9/5 + 32
}

func BytesToKB(b uint64) float64 {
	return float64(b) / float64(KB)
}

func BytesToMB(b uint64) float64 {
	return float64(b) / float64(MB)
}

func BytesToGB(b uint64) float64 {
	return float64(b) / float64(GB)
}

func BytesToTB(b uint64) float64 {
	return float64(b) / float64(TB)
}

func ConvertBytes(b uint64) (float64, string) {
	switch {
	case b < KB:
		return float64(b), "B"
	case b < MB:
		return BytesToKB(b), "KB"
	case b < GB:
		return BytesToMB(b), "MB"
	case b < TB:
		return BytesToGB(b), "GB"
	default:
		return BytesToTB(b), "TB"
	}
}
