package model

import "singo/util"

//执行数据迁移

func migration() {
	// 自动迁移模式
	if err := DB.AutoMigrate(&User{}); err != nil {
		util.Log().Error("AutoMigrate User Failed", err)
	}

}
