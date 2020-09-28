package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"singo/util"
	"strconv"
)

// RedisClient Redis缓存客户端单例
var RedisClient *redis.Client

// Redis 在中间件中初始化redis链接
func Redis() {
	db, _ := strconv.ParseUint(viper.GetString("redis.dbname"), 10, 64)
	host := viper.GetString("redis.host")
	port := viper.GetString("redis.port")

	addr := fmt.Sprintf("%v:%v", host, port)
	client := redis.NewClient(&redis.Options{
		Addr:       addr,
		Password:   viper.GetString("redis.password"),
		DB:         int(db),
		MaxRetries: 1,
	})

	_, err := client.Ping().Result()

	if err != nil {
		util.Log().Error("连接Redis不成功", err)
	} else {
		RedisClient = client
	}
}
