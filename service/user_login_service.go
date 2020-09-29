package service

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"singo/model"
)

// UserLoginDto 管理用户登录的服务
type UserLoginDto struct {
	UserName string `form:"username" json:"username" binding:"required,min=5,max=30"`
	Password string `form:"password" json:"password" binding:"required,min=4,max=40"`
}

// Login 用户登录函数
func Login(c *gin.Context) (interface{}, error) {
	var userLoginDto UserLoginDto
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
