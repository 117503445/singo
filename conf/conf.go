package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"singo/cache"
	"singo/model"
	"singo/util"

	"github.com/spf13/viper"
)


// Init 初始化配置项
func Init() {
	filepathBase := filepath.Dir(util.GetCurrentPath())
	filePathEnv := filepath.Join(filepathBase, "config.yaml")

	viper.SetConfigFile(filePathEnv) // 配置文件

	if err := viper.ReadInConfig(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "load config failed")
		panic(err)
	}

	// 设置日志级别
	util.BuildLogger(viper.GetString("log.level"))

	// 读取翻译文件
	if err := LoadLocales(filepath.Join(util.GetCurrentPath(), "locales", "zh-cn.yaml")); err != nil {
		util.Log().Panic("翻译文件加载失败", err)
	}

	// 连接数据库
	model.Database()
	cache.Redis()
}
