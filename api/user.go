package api

import (
	"singo/model"
	"singo/serializer"
	"singo/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// UserRegister 用户注册接口
func UserRegister(c *gin.Context) {
	var service service.UserRegisterDto


	if err := c.ShouldBindJSON(&service); err == nil {

		count := int64(0)
		model.DB.Model(&model.User{}).Where("username = ?", service.UserName).Count(&count)
		if count > 0 {
			c.JSON(200, gin.H{"message": "Username has already exists."})
			return
		}

		user, _ := service.Register()
		c.JSON(200, user)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

//todo

// UserLogin 用户登录接口
//func UserLogin(c *gin.Context) {
//	var service service.UserLoginDto
//	if err := c.ShouldBindJSON(&service); err == nil {
//		res := service.Login(c)
//		c.JSON(200, res)
//	} else {
//		c.JSON(200, ErrorResponse(err))
//	}
//}

// UserMe 用户详情
func UserMe(c *gin.Context) {
	user := CurrentUser(c)
	//res := serializer.BuildUserResponse(*user)
	model.DB.Preload("Roles").Find(&user)
	c.JSON(200, user)
}

// UserLogout 用户登出
func UserLogout(c *gin.Context) {
	s := sessions.Default(c)
	s.Clear()
	s.Save()
	c.JSON(200, serializer.Response{
		Code: 0,
		Msg:  "登出成功",
	})
}
