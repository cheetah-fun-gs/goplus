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
		return "", fmt.Errorf("%s is not found in %s", name, root)
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

// ReplaceOption 替换参数
type ReplaceOption struct {
	Old      string
	New      string
	IsRegexp bool
}

// CopyAndReplace 拷贝文件并替换
func CopyAndReplace(src, dst string, replaces []*ReplaceOption) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	bytes, err := ioutil.ReadAll(srcFile)
	if err != nil {
		return err
	}

	text := string(bytes)
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

	_, err = io.WriteString(dstFile, text)
	return err
}

// Copy 拷贝文件
func Copy(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile) // 把原来文件的内容拷贝到新文件中
	if err != nil {
		return err
	}
	return nil
}
