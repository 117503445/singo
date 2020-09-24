package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"singo/cache"
	"singo/model"
	"singo/util"

	"github.com/joho/godotenv"
)

// Init 初始化配置项
func Init() {
	filepathBase := filepath.Dir(util.GetCurrentPath())
	filePathEnv := filepath.Join(filepathBase, ".env")
	// 从本地读取环境变量
	if err := godotenv.Load(filePathEnv); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "load config failed")
		panic(err)
	}

	// 设置日志级别
	util.BuildLogger(os.Getenv("LOG_LEVEL"))

	// 读取翻译文件
	if err := LoadLocales(filepath.Join(util.GetCurrentPath(), "locales", "zh-cn.yaml")); err != nil {
		util.Log().Panic("翻译文件加载失败", err)
	}

	// 连接数据库
	model.Database(os.Getenv("MYSQL_DSN"))
	cache.Redis()
}
