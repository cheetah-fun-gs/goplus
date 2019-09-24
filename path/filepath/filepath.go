// Package filepath 路径, 目录, 文件名方法
package filepath

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
		return "", fmt.Errorf("not found")
	}
	return p, nil
}

// DirNames 获取目录下的目录名列表
func DirNames(dirPath string) ([]string, error) {
	dirNames := []string{}

	fpList, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, fp := range fpList {
		if !fp.IsDir() {
			continue
		}
		dirNames = append(dirNames, fp.Name())
	}
	return dirNames, nil
}

// FileNames 获取目录下的文件名列表
func FileNames(dirPath string) ([]string, error) {
	fileNames := []string{}

	fpList, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, fp := range fpList {
		if fp.IsDir() {
			continue
		}
		fileNames = append(fileNames, fp.Name())
	}
	return fileNames, nil
}

// CopyFile 拷贝文件
func CopyFile(src, dst string) error {
	srcAbs, err := filepath.Abs(src)
	if err != nil {
		return err
	}
	dstAbs, err := filepath.Abs(dst)
	if err != nil {
		return err
	}

	if srcAbs == dstAbs {
		return fmt.Errorf("dst must different from src")
	}

	srcFile, err := os.Open(srcAbs)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstAbs)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// ReplaceOption 替换参数
type ReplaceOption struct {
	Old      string
	New      string
	IsRegexp bool
}

// CopyFileAndReplace 拷贝文件并替换
func CopyFileAndReplace(src, dst string, replaces []*ReplaceOption) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	data, err := ioutil.ReadAll(srcFile)
	if err != nil {
		return err
	}

	text := string(data)
	for _, r := range replaces {
		if !r.IsRegexp {
			text = strings.Replace(text, r.Old, r.New, -1)
		} else {
			ok, err := regexp.Match(r.Old, []byte(text))
			if err != nil {
				return err
			}
			if ok {
				text = regexp.MustCompile(r.Old).ReplaceAllString(text, r.New)
			}
		}
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.WriteString(dstFile, text)
	return err
}
