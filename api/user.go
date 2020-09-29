package api

import (
	"net/http"
	"singo/dto"
	"singo/model"
	"singo/serializer"
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
			c.JSON(http.StatusBadRequest, serializer.Err(serializer.StatusUsernameRepeat, "Username has already exists.", nil))
			return
		}

		user, _ := service.Register(&userRegisterIn)
		c.JSON(http.StatusOK, user)
	} else {
		c.JSON(http.StatusBadRequest, serializer.Err(http.StatusBadRequest, "bad UserRegisterIn dto.", err))
	}
}

// UserMe 用户详情
func UserMe(c *gin.Context) {
	user := service.CurrentUser(c)
	model.DB.Preload("Roles").Find(&user)
	c.JSON(200, user)
}
