package service

import (
	"gorm.io/gorm"
	"singo/model"
	"singo/serializer"
)

// UserRegisterDto 管理用户注册服务
type UserRegisterDto struct {
	UserName string `form:"username" json:"username" binding:"required,min=5,max=30"`
	Password string `form:"password" json:"password" binding:"required,min=4,max=40"`
}

// valid 验证表单
func (service *UserRegisterDto) valid() *serializer.Response {

	//count := int64(0)
	//model.DB.Model(&model.User{}).Where("nickname = ?", service.Nickname).Count(&count)
	//if count > 0 {
	//	return &serializer.Response{
	//		Code: 40001,
	//		Msg:  "昵称被占用",
	//	}
	//}

	count := int64(0)
	model.DB.Model(&model.User{}).Where("username = ?", service.UserName).Count(&count)
	if count > 0 {
		return &serializer.Response{
			Code: 40001,
			Msg:  "用户名已经注册",
		}
	}

	return nil
}

// Register 用户注册
func (service *UserRegisterDto) Register() (model.User, error) {
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
