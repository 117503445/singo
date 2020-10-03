package util

import (
	"path"
	"runtime"
)

var DirPassword string
var FilePasswordAdmin string
var FilePasswordJWT string

// 获取调用者 go 文件 所在文件夹
func GetCurrentPath() string {
	_, filename, _, _ := runtime.Caller(1)

	return path.Dir(filename)
}
