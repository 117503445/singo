package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"io/ioutil"
	"net/http"
	"singo/dto"
	"singo/model"
	"singo/serializer"
	"singo/service"
	"singo/util"
)

// UserCreate 用户注册接口
func UserCreate(c *gin.Context) {
	userCreateIn := &dto.UserCreateIn{}
	var err error
	if err = c.ShouldBindJSON(&userCreateIn); err != nil {
		c.JSON(http.StatusBadRequest, serializer.Err(http.StatusBadRequest, "bad UserCreateIn dto.", err))
		return
	}

	if err = validator.New().Struct(userCreateIn); err != nil {
		c.JSON(http.StatusBadRequest, serializer.Err(serializer.StatusParamNotValid, "StatusParamNotValid", err))
		return
	}

	var user *model.User
	if user, err = userCreateIn.ToUser(); err != nil {
		c.JSON(http.StatusInternalServerError, serializer.Err(serializer.StatusDtoToModelError, "userRegisterInToUser failed", err))
		return
	}

	count := int64(0)
	model.DB.Model(&model.User{}).Where("username = ?", userCreateIn.UserName).Count(&count)
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
	//userUpdateIn := &dto.UserUpdateIn{}
	//var err error

	newUserBytes, _ := ioutil.ReadAll(c.Request.Body)
	var mapNewUser map[string]interface{}
	if err := json.Unmarshal(newUserBytes, &mapNewUser); err != nil {
		fmt.Println(err)
	}

	user := service.CurrentUser(c)

	if username, ok := mapNewUser["username"]; ok {
		if len(username.(string)) < 5 {
			c.JSON(http.StatusBadRequest, serializer.Err(serializer.StatusParamNotValid, "username length should in [5,30]", nil))
			return
		}

		if len(username.(string)) > 30 {
			c.JSON(http.StatusBadRequest, serializer.Err(serializer.StatusParamNotValid, "username length should in [5,30]", nil))
			return
		}
	}
	if password, ok := mapNewUser["password"]; ok {
		if len(password.(string)) < 4 {
			c.JSON(http.StatusBadRequest, serializer.Err(serializer.StatusParamNotValid, "password length should in [4,30]", nil))
			return
		}

		if len(password.(string)) > 30 {
			c.JSON(http.StatusBadRequest, serializer.Err(serializer.StatusParamNotValid, "password length should in [4,30]", nil))
			return
		}
	}

	if newName, ok := mapNewUser["username"]; ok {
		count := int64(0)
		model.DB.Model(&model.User{}).Where("username = ?", newName).Count(&count)
		if count > 0 && newName != user.Username {
			// 修改 Username 后 发生重名
			c.JSON(http.StatusBadRequest, serializer.Err(serializer.StatusUsernameRepeat, "Username has already exists.", nil))
			return
		}
	}

	util.SetStructFieldByMap(user, mapNewUser, []string{"username", "password", "avatar"})
	if _, ok := mapNewUser["password"]; ok {
		if err := user.SetPassword(user.Password); err != nil {
			c.JSON(http.StatusInternalServerError, serializer.Err(serializer.StatusModelToDtoError, "SetPassword failed", err))
			return
		}
	}

	model.DB.Save(user)

	if userOut, err := dto.UserToUserOut(user); err == nil {
		c.JSON(http.StatusOK, userOut)
	} else {
		c.JSON(http.StatusInternalServerError, serializer.Err(serializer.StatusModelToDtoError, "UserToUserOut failed", err))
	}

}
