package dto

import (
	"github.com/devfeel/mapper"
	"singo/model"
)

type UserCreateIn struct {
	UserName string `json:"username" validate:"required,min=5,max=30"`
	Password string `json:"password" validate:"required,min=4,max=40"`
	Avatar   string `json:"avatar" gorm:"size:1000"`
}

func (userDto UserCreateIn) ToUser() (*model.User, error) {
	user := &model.User{}
	if err := mapper.AutoMapper(&userDto, user); err != nil {
		return user, err
	}
	// 加密密码
	if err := user.SetPassword(user.Password); err != nil {
		return user, err
	}

	return user, nil
}
