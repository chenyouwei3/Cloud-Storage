package dist_storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// 新创建文件夹
func MakeDir(targetPath string) error {
	return os.Mkdir(targetPath, 0755)
}

// 把文件或者文件夹重命名
func RenameFileOrDir(oldPath, newName string) error {
	// 获取原路径的目录部分
	dir := filepath.Dir(oldPath)
	// 拼接新路径
	newPath := filepath.Join(dir, newName)
	// 使用 os.Rename 进行重命名（适用于文件或目录）
	fmt.Println("重命名1111", oldPath, newPath)
	err := os.Rename(oldPath, newPath)
	if err != nil {
		fmt.Println("重命名失败", err)
		return fmt.Errorf("重命名失败: %w", err)
	}
	return nil
}

// 删除文件或者文件夹
func RemoveFileOrDir(targetPath string) error {
	err := os.RemoveAll(targetPath)
	if err != nil {
		return fmt.Errorf("删除失败: %w", err)
	}

	return nil
}

// 根据路径判断是文件还是目录，执行复制
func CopyFileOrDir(srcPath, dstPath string) error {
	info, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return copyDir(srcPath, dstPath)
	}
	return copyFile(srcPath, dstPath)
}

// 将文件从 sourcePath 移动到 targetPath
func MoveFile(sourcePath, targetPath string) error {
	// 确保目标目录存在
	err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm)
	if err != nil {
		return err
	}

	// 移动文件（重命名）
	err = os.Rename(sourcePath, targetPath)
	if err != nil {
		return err
	}

	return nil
}

// 下拉菜单 - 只返回文件夹的路径
func DirMapping(Path string) ([]string, error) {
	var folders []string

	if err := filepath.WalkDir(Path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 只收集目录（排除文件）
		if d.IsDir() {
			folders = append(folders, path)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return folders, nil
}

// 递归复制目录
func copyDir(srcDir, dstDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算相对路径
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dstDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}
		return copyFile(path, targetPath)
	})
}

// 复制文件
func copyFile(srcFile, dstFile string) error {
	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer dst.Close()

	// 拷贝内容
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	// 保持文件权限一致
	srcInfo, err := os.Stat(srcFile)
	if err == nil {
		err = os.Chmod(dstFile, srcInfo.Mode())
	}
	return err
}
