package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if err := StatFileTypes("../cloud"); err != nil {
		panic(err)
	}

}

// 统计文件夹下所有文件的类型、大小和占比
func StatFileTypes(dirPath string) error {
	typeStats := make(map[string]int64)
	var totalSize int64 = 0

	// 遍历目录
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 跳过目录
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(info.Name()))
		size := info.Size()

		typeStats[ext] += size
		totalSize += size

		return nil
	})

	if err != nil {
		return err
	}

	// 打印统计信息
	fmt.Println("文件类型统计：")
	for ext, size := range typeStats {
		percentage := float64(size) / float64(totalSize) * 100
		fmt.Printf("类型: %-10s 大小: %-10d 字节 占比: %.2f%%\n", ext, size, percentage)
	}
	fmt.Printf("总大小: %d 字节\n", totalSize)

	return nil
}
