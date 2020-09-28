package server

import (
	"singo/api"
	"singo/middleware"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()

	// 中间件, 顺序不能改
	r.Use(middleware.Session(viper.GetString("session.secret")))
	r.Use(middleware.Cors())
	r.Use(middleware.CurrentUser())

	// 路由
	v1 := r.Group("/api/v1")
	{
		v1.POST("ping", api.Ping)

		// 用户登录
		v1.POST("user/register", api.UserRegister)

		// 用户登录
		v1.POST("user/login", middleware.JwtMiddleware.LoginHandler)

		// 需要登录保护的
		auth := v1.Group("")
		auth.Use(middleware.JwtMiddleware.MiddlewareFunc())
		{
			// User Routing
			auth.GET("user/me", api.UserMe)
			auth.DELETE("user/logout", api.UserLogout)
		}
	}
	return r
}
