package router

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

		user := v1.Group("user")
		{
			// 用户登录
			user.POST("register", api.UserCreate)

			// 用户登录
			user.POST("login", middleware.JwtMiddleware.LoginHandler)

			auth := user.Group("")
			// 需要登录保护
			auth.Use(middleware.JwtMiddleware.MiddlewareFunc())
			{
				// User Routing
				auth.GET("me", api.UserRead)
				auth.PUT("", api.UserUpdate)
			}
		}
		admin := v1.Group("admin")
		// 需要登录保护
		admin.Use(middleware.JwtMiddleware.MiddlewareFunc())
		admin.Use(middleware.HasRole("admin"))
		{
			user := admin.Group("user")
			{
				user.GET(":id", api.AdminUserRead)
			}
		}
	}
	return r
}
