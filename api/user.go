package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"singo/dto"
	"singo/model"
	"singo/serializer"
	"singo/service"
)

// UserCreate 用户注册接口
func UserCreate(c *gin.Context) {
	userRegisterIn := &dto.UserCreateUpdateIn{}
	var err error
	if err = c.ShouldBindJSON(&userRegisterIn); err != nil {
		c.JSON(http.StatusBadRequest, serializer.Err(http.StatusBadRequest, "bad UserCreateUpdateIn dto.", err))
		return
	}

	if err = validator.New().Struct(userRegisterIn); err != nil {
		c.JSON(http.StatusBadRequest, serializer.Err(serializer.StatusParamNotValid, "StatusParamNotValid", err))
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

// UserRead 用户详情
func UserRead(c *gin.Context) {
	user := service.CurrentUser(c)
	if userOut, err := dto.UserToUserOut(user); err == nil {
		c.JSON(http.StatusOK, userOut)
	} else {
		c.JSON(http.StatusInternalServerError, serializer.Err(serializer.StatusModelToDtoError, "UserToUserOut failed", err))
	}

}

// UserUpdate 更新用户信息
func UserUpdate(c *gin.Context) {
	userRegisterIn := &dto.UserCreateUpdateIn{}
	var err error
	if err = c.ShouldBindJSON(&userRegisterIn); err != nil {
		c.JSON(http.StatusBadRequest, serializer.Err(http.StatusBadRequest, "bad UserCreateUpdateIn dto.", err))
		return
	}

	if err = validator.New().Struct(userRegisterIn); err != nil {
		c.JSON(http.StatusBadRequest, serializer.Err(serializer.StatusParamNotValid, "StatusParamNotValid", err))
		return
	}

	count := int64(0)
	model.DB.Model(&model.User{}).Where("username = ?", userRegisterIn.UserName).Count(&count)
	if count > 0 && service.CurrentUser(c).Username != userRegisterIn.UserName {
		// 修改 Username 后 发生重名
		c.JSON(http.StatusBadRequest, serializer.Err(serializer.StatusUsernameRepeat, "Username has already exists.", nil))
		return
	}

	user, _ := userRegisterIn.ToUser()
	user.Model = service.CurrentUser(c).Model
	model.DB.Save(user)

	if userOut, err := dto.UserToUserOut(user); err == nil {
		c.JSON(http.StatusOK, userOut)
	} else {
		c.JSON(http.StatusInternalServerError, serializer.Err(serializer.StatusModelToDtoError, "UserToUserOut failed", err))
	}

}
