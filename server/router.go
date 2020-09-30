package server

import (
	"singo/api"
	"singo/middleware"

	"github.com/gin-gonic/gin"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Cors())

	// 路由
	v1 := r.Group("/api/v1")
	{
		v1.POST("ping", api.Ping)

		// 用户登录
		v1.POST("user/register", api.UserCreate)

		// 用户登录
		v1.POST("user/login", middleware.JwtMiddleware.LoginHandler)

		// 需要登录保护的
		auth := v1.Group("")
		auth.Use(middleware.JwtMiddleware.MiddlewareFunc())
		{
			// User Routing
			auth.GET("user/me", api.UserRead)
			auth.PUT("user",api.UserUpdate)
		}
	}
	return r
}
