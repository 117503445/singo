package model

import (
	"database/sql"
	"fmt"
	"singo/util"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"

	//
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// DB 数据库链接单例
var DB *gorm.DB

// Database 在中间件中初始化mysql链接
func Database() {

	host := viper.Get("mysql.host")
	port := viper.Get("mysql.port")
	username := viper.Get("mysql.username")
	password := viper.Get("mysql.password")
	dbname := viper.Get("mysql.dbname")

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/?charset=utf8&parseTime=True&loc=Local", username, password, host, port)
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		util.Log().Panic("连接数据库不成功", err)
	}
	_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %v", dbname))
	if err != nil {
		util.Log().Panic("创建数据库不成功", err)
	}

	dsn = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbname)
	ormDB, err := gorm.Open("mysql", dsn)
	// Error
	if err != nil {
		util.Log().Panic("连接数据库不成功", err)
	}

	ormDB.LogMode(true)

	//设置连接池
	//空闲
	ormDB.DB().SetMaxIdleConns(50)
	//打开
	ormDB.DB().SetMaxOpenConns(100)
	//超时
	ormDB.DB().SetConnMaxLifetime(time.Second * 30)

	DB = ormDB

	migration()
}
