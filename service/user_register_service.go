package service

import (
	"gorm.io/gorm"
	"singo/dto"
	"singo/model"
)

//// valid 验证表单
//func (service *dto.UserRegisterIn) valid() *serializer.Response {
//
//	//count := int64(0)
//	//model.DB.Model(&model.User{}).Where("nickname = ?", service.Nickname).Count(&count)
//	//if count > 0 {
//	//	return &serializer.Response{
//	//		Code: 40001,
//	//		Msg:  "昵称被占用",
//	//	}
//	//}
//
//	count := int64(0)
//	model.DB.Model(&model.User{}).Where("username = ?", service.UserName).Count(&count)
//	if count > 0 {
//		return &serializer.Response{
//			Code: 40001,
//			Msg:  "用户名已经注册",
//		}
//	}
//
//	return nil
//}

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
