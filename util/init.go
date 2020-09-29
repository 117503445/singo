package util

import (
	"os"
	"path"
)

func init() {
	DirPassword = path.Dir(GetCurrentPath()) + "/data/password"
	FilePasswordAdmin = DirPassword + "/admin.txt"
	if err := os.MkdirAll(DirPassword, os.ModePerm); err != nil {
		Log().Panic("can't create dir data/password", err)
	}
}
