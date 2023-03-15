package util

import (
	"path"
	"strings"
)

// BaseName
// 获取路径中的文件名，不包含扩展名
func NameWithoutExt(p string) string {
	name := path.Base(p)
	ext := path.Ext(name)
	return strings.TrimSuffix(name, ext)
}
