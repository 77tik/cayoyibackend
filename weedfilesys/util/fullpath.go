package util

import (
	"path/filepath"
	"strings"
)

type FullPath string

// dir="a/b" name="c.txt"
// return: a/b/c.txt
func NewFullPath(dir, name string) FullPath {
	name = strings.TrimSuffix(name, "/")
	return FullPath(dir + "/" + name)
}

// 拆分fullpath为dir 和 name
func (fp FullPath) DirAndName() (string, string) {
	// /a/b/c.txt => dir="/a/b/",name="c.txt"
	// /hello.txt =? dir="/",name-"hello.txt" 这种情况直接返回了
	dir, name := filepath.Split(string(fp))

	// 如果name中存在非法的UTF8字符，就用？代替
	name = strings.ToValidUTF8(name, "?")
	if dir == "/" {
		return dir, name
	}
	if len(dir) < 1 {
		return "/", ""
	}

	// 因为filepath的目录总是返回 "/" 结尾的，所以要剔除最后边的"/"
	return dir[:len(dir)-1], name
}
func (fp FullPath) Name() string {
	_, name := filepath.Split(string(fp))
	name = strings.ToValidUTF8(name, "?")
	return name
}
