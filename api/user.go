package api

import (
	"singo/dto"
	"singo/model"
	"singo/service"

	"github.com/gin-gonic/gin"
)

// UserRegister 用户注册接口
func UserRegister(c *gin.Context) {
	var userRegisterIn dto.UserRegisterIn

	if err := c.ShouldBindJSON(&userRegisterIn); err == nil {

		count := int64(0)
		model.DB.Model(&model.User{}).Where("username = ?", userRegisterIn.UserName).Count(&count)
		if count > 0 {
			c.JSON(200, gin.H{"message": "Username has already exists."})
			return
		}

		user, _ := service.Register(&userRegisterIn)
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
	model.DB.Preload("Roles").Find(&user)
	c.JSON(200, user)
}
