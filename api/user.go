package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"singo/dto"
	"singo/model"
	"singo/serializer"
	"singo/service"
)

// UserRegister 用户注册接口
func UserRegister(c *gin.Context) {
	userRegisterIn := &dto.UserRegisterIn{}
	var err error
	if err = c.ShouldBindJSON(&userRegisterIn); err != nil {
		c.JSON(http.StatusBadRequest, serializer.Err(http.StatusBadRequest, "bad UserRegisterIn dto.", err))
		return
	}

	var user *model.User
	if user, err = userRegisterIn.ToUser(); err != nil {
		c.JSON(http.StatusInternalServerError, serializer.Err(serializer.StatusDtoToModelError, "userRegisterInToUser failed", err))
		return
	}

	count := int64(0)
	model.DB.Model(&model.User{}).Where("username = ?", userRegisterIn.UserName).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, serializer.Err(serializer.StatusUsernameRepeat, "Username has already exists.", nil))
		return
	}

	if user, err = service.Register(user); err != nil {
		c.JSON(http.StatusInternalServerError, serializer.Err(serializer.StatusRegisterError, "Register failed", err))
		return
	}

	if userOut, err := dto.UserToUserOut(user); err == nil {
		c.JSON(http.StatusOK, userOut)
	} else {
		c.JSON(http.StatusInternalServerError, serializer.Err(serializer.StatusModelToDtoError, "UserToUserOut failed", err))
	}
}

// UserMe 用户详情
func UserMe(c *gin.Context) {
	user := service.CurrentUser(c)
	if userOut, err := dto.UserToUserOut(user); err == nil {
		c.JSON(http.StatusOK, userOut)
	} else {
		c.JSON(http.StatusInternalServerError, serializer.Err(serializer.StatusModelToDtoError, "UserToUserOut failed", err))
	}

}
