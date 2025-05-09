package dist_storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// distInfos 用于封装目录信息和树结构
type distInfos struct {
	DefaultPath string               `json:"defaultPath"` // 根目录路径
	TotalDist   string               `json:"totalDist"`   // 总大小，格式化后的字符串
	Item        map[string]*distInfo `json:"item"`        // 文件/目录的详细信息映射
	Tree        *TreeNode            `json:"tree"`        // 目录树结构
}

// distInfo 表示单个文件或目录的信息
type distInfo struct {
	IsDir bool   `json:"isDir"` // 是否为目录
	Size  string `json:"size"`  // 大小，格式化后的字符串
	Date  string `json:"date"`  // 最后修改时间，格式化后的字符串
}

// TreeNode 用于表示目录树节点
type TreeNode struct {
	Name     string      `json:"name"`               // 节点名称，对应文件或目录名
	Children []*TreeNode `json:"children,omitempty"` // 子节点列表，使用omitempty避免空 slice 输出
}

// 获取目录信息和树结构
func GetDirInfoWithTree(rootPath string, fileTypes []string) (*distInfos, error) {
	// 用于存放每个相对路径对应distinfo
	items := make(map[string]*distInfo)
	// totalSize 用于累加所有文件大小，目录通过递归计算
	var totalSize int64

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		//查找文件类型
		if err != nil || (len(fileTypes) > 0 && !isValidFileType(path, fileTypes)) {
			return nil
		}
		relPath, _ := filepath.Rel(filepath.Dir(rootPath), path)
		if relPath == "." {
			relPath = info.Name()
		}
		var size int64
		if info.IsDir() {
			size = calcDirSize(path)
		} else {
			size = info.Size()
		}
		totalSize += size

		items[filepath.ToSlash(relPath)] = &distInfo{
			IsDir: info.IsDir(),
			Size:  formatSize(size),
			Date:  info.ModTime().Format("2006-01-02 15:04:05"),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	tree := buildTreeFromPath(rootPath, fileTypes)

	return &distInfos{
		DefaultPath: rootPath,
		TotalDist:   formatSize(totalSize),
		Item:        items,
		Tree:        tree,
	}, nil
}

// 构建目录树结构（修复了重复节点问题）
func buildTreeFromPath(root string, fileTypes []string) *TreeNode {
	rootNode := &TreeNode{Name: filepath.Base(root)}
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || (len(fileTypes) > 0 && !isValidFileType(path, fileTypes)) {
			return nil
		}
		/* relPath, _ := filepath.Rel(filepath.Dir(root), path) */
		relPath, _ := filepath.Rel(root, path)
		if relPath == "." {
			return nil
		}
		// 使用 ToSlash 统一路径分隔符，保证跨平台一致性
		parts := strings.Split(filepath.ToSlash(relPath), "/")
		insertPath(rootNode, parts)
		return nil
	})
	return rootNode
}

// 将路径按层级插入树节点中，避免重复
func insertPath(current *TreeNode, parts []string) {
	if len(parts) == 0 {
		return
	}
	part := parts[0]

	// 查找是否已经存在该子节点
	var child *TreeNode
	for _, c := range current.Children {
		if c.Name == part {
			child = c
			break
		}
	}
	if child == nil {
		child = &TreeNode{Name: part}
		current.Children = append(current.Children, child)
	}

	insertPath(child, parts[1:])
}

// 判断文件后缀是否在 fileTypes 中
func isValidFileType(path string, fileTypes []string) bool {
	for _, fileType := range fileTypes {
		if strings.HasSuffix(path, "."+fileType) {
			return true
		}
	}
	return false
}

// 递归计算目录大小
func calcDirSize(path string) int64 {
	var size int64
	_ = filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size
}

// 格式化文件大小为易读单位
func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
