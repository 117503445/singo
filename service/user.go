package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"singo/dto"
	"singo/model"
)

// Register 用户注册
func Register(service *dto.UserRegisterIn) (model.User, error) {
	user := model.User{
		Username: service.UserName,
	}

	// 加密密码
	if err := user.SetPassword(service.Password); err != nil {
		return user, err
	}

	roleName := "user"
	role, err := model.QueryRoleByName(roleName)
	if err == gorm.ErrRecordNotFound {
		role = model.Role{
			Name: "user",
		}
	}
	user.Roles = []model.Role{role}

	// 创建用户
	if err := model.DB.Create(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

// CurrentUser 获取当前用户
func CurrentUser(c *gin.Context) *model.User {
	if user, _ := c.Get("user"); user != nil {
		if u, ok := user.(*model.User); ok {
			return u
		}
	}
	return nil
}
