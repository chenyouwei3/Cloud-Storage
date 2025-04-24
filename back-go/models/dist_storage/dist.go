package dist_storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type distInfos struct {
	DefaultPath string               `json:"defaultPath"`
	TotalDist   string               `json:"totalDist"`
	Item        map[string]*distInfo `json:"item"`
	Tree        *TreeNode            `json:"tree"`
}

type distInfo struct {
	IsDir bool   `json:"isDir"`
	Size  string `json:"size"`
	Date  string `json:"date"`
}

type TreeNode struct {
	Name     string      `json:"name"`
	Children []*TreeNode `json:"children,omitempty"`
}

// 获取目录信息和树结构
func GetDirInfoWithTree(rootPath string) (*distInfos, error) {
	items := make(map[string]*distInfo)
	var totalSize int64

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
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

	tree := buildTreeFromPath(rootPath)

	return &distInfos{
		DefaultPath: rootPath,
		TotalDist:   formatSize(totalSize),
		Item:        items,
		Tree:        tree,
	}, nil
}

// 构建目录树结构（修复了重复节点问题）
func buildTreeFromPath(root string) *TreeNode {
	tree := &TreeNode{Name: filepath.Base(filepath.Dir(root))}
	rootNode := &TreeNode{Name: filepath.Base(root)}
	tree.Children = []*TreeNode{rootNode}

	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		relPath, _ := filepath.Rel(filepath.Dir(root), path)
		if relPath == "." {
			return nil
		}
		// 使用 ToSlash 统一路径分隔符，保证跨平台一致性
		parts := strings.Split(filepath.ToSlash(relPath), "/")
		insertPath(rootNode, parts)
		return nil
	})

	return tree
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
