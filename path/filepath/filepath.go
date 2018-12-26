package filepath

import (
	"fmt"
	"os"
	"path/filepath"
)

// NoExt 去掉扩展名
func NoExt(path string) string {
	for i := len(path) - 1; i >= 0 && !os.IsPathSeparator(path[i]); i-- {
		if path[i] == '.' {
			return path[:i]
		}
	}
	return ""
}

// Exists 判断路径是否存在
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Files 遍历文件夹，获取所有文件
func Files(root string) ([]string, error) {
	files := []string{}
	walkFn := func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	}
	err := filepath.Walk(root, walkFn)
	if err != nil {
		return []string{}, err
	}
	return files, nil
}

// Dirs 遍历文件夹，获取所有文件夹
func Dirs(root string) ([]string, error) {
	dirs := []string{}
	walkFn := func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if info.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	}
	err := filepath.Walk(root, walkFn)
	if err != nil {
		return []string{}, err
	}
	return dirs, nil
}

// Find 遍历文件夹，查找指定路径
func Find(root, name string, isDir bool) (string, error) {
	p := ""
	walkFn := func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if info.IsDir() == isDir && info.Name() == name {
			p = path
		}
		return nil
	}
	err := filepath.Walk(root, walkFn)
	if err != nil {
		return "", err
	}
	if p == "" {
		return "", fmt.Errorf("%s is not found in %s", name, root)
	}
	return p, nil
}
