package service

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"singo/dto"

	"singo/model"
)

// Login 用户登录函数
func Login(c *gin.Context) (interface{}, error) {
	var userLoginDto dto.UserLoginIn
	if err := c.ShouldBindJSON(&userLoginDto); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	username := userLoginDto.UserName
	password := userLoginDto.Password
	queryUser, err := model.QueryByUsername(username)
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}
	if queryUser.CheckPassword(password) {
		return &queryUser, nil
	}
	return nil, jwt.ErrFailedAuthentication
}
