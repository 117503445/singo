package dto

import (
	"github.com/devfeel/mapper"
	"singo/model"
)

type UserRegisterIn struct {
	UserName string `form:"username" json:"username" binding:"required,min=5,max=30"`
	Password string `form:"password" json:"password" binding:"required,min=4,max=40"`
	Avatar   string `json:"avatar" gorm:"size:1000"`
}

func (userRegisterIn UserRegisterIn) ToUser() (*model.User, error) {
	user := &model.User{}
	if err := mapper.AutoMapper(&userRegisterIn, user); err != nil {
		return user, err
	}
	// 加密密码
	if err := user.SetPassword(user.Password); err != nil {
		return user, err
	}

	return user, nil
}
