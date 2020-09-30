package model

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io/ioutil"
	"singo/util"
	"strings"
)

// DB 数据库链接单例
var DB *gorm.DB

// InitDatabase 在中间件中初始化mysql链接
func InitDatabase() {

	host := viper.Get("mysql.host")
	port := viper.Get("mysql.port")
	username := viper.Get("mysql.username")
	password := viper.Get("mysql.password")
	dbname := viper.GetString("mysql.dbname")

	if strings.Contains(dbname, "test") {
		util.Log().Info("dbname contain \"test\", drop and create new\n")
		_, _ = Exec("drop database " + dbname)
	}

	_, err := Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %v", dbname))
	if err != nil {
		util.Log().Panic("创建数据库不成功", err)
	}

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbname)
	ormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// Error
	if err != nil {
		util.Log().Panic("连接数据库不成功", err)
	}
	//
	//ormDB.LogMode(true)
	//
	////设置连接池
	////空闲
	//ormDB.DB().SetMaxIdleConns(50)
	////打开
	//ormDB.DB().SetMaxOpenConns(100)
	////超时
	//ormDB.DB().SetConnMaxLifetime(time.Second * 30)

	DB = ormDB

	migration()
	createAdminUser()
}

//createAdminUser 创建管理员账号
//用户名 admin,密码随机12位,权限 Admin,User
//输出在 data/password/admin.txt,
func createAdminUser() {
	var err error

	var password string
	bytes, err := ioutil.ReadFile(util.FilePasswordAdmin)
	if err == nil && len(string(bytes)) != 0 {
		password = string(bytes)

	} else {
		password = util.RandStringRunes(12)
		_ = ioutil.WriteFile(util.FilePasswordAdmin, []byte(password), 0777)
	}

	user, err := ReadUserByName("admin")
	if err != nil {
		// 如果没找到 admin 账号,就重新创建
		user = User{
			Username: "admin",
			Roles:    []Role{{Name: "admin"}, {Name: "user"}},
			Avatar:   "",
		}

		if err = user.SetPassword(password); err != nil {
			util.Log().Error("set admin password failed", err)
		}

		DB.Create(&user)
	} else {
		// 如果找到了 admin 账号,就重新设置密码
		if err = user.SetPassword(password); err != nil {
			util.Log().Error("set admin password failed", err)
		}

		DB.Save(&user)
	}

}

// Exec 执行单条 SQL
func Exec(query string) (sql.Result, error) {
	host := viper.Get("mysql.host")
	port := viper.Get("mysql.port")
	username := viper.Get("mysql.username")
	password := viper.Get("mysql.password")

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/?charset=utf8&parseTime=True&loc=Local", username, password, host, port)
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		util.Log().Panic("连接数据库不成功", err)
	}
	result, err := sqlDB.Exec(query)
	return result, err
}
