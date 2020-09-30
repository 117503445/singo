package model

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Roles    []Role `json:"roles" gorm:"many2many:user_role"`
	Avatar   string `json:"avatar" gorm:"size:1000"`
}

const (
	// PassWordCost 密码加密难度
	PassWordCost = 12
)

// ReadUserById 用ID获取用户
func ReadUserById(id int) (*User, error) {
	var user User
	result := DB.First(&user, id)
	return &user, result.Error
}

// SetPassword 设置密码
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

// CheckPassword 校验密码
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

// ReadUserByName 用 Username 获取用户
func ReadUserByName(username string) (User, error) {
	var user User
	result := DB.Where("username = ?", username).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, result.Error
	}

	return user, nil
}
