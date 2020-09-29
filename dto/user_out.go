package dto

import (
	"github.com/devfeel/mapper"
	"singo/model"
	"singo/util"
)

type UserOut struct {
	ID       uint     `json:"id"`
	Username string   `json:"username"`
	Role     []string `json:"role"`
	Avatar   string   `json:"avatar" gorm:"size:1000"`
}

func UserToUserOut(user *model.User) (*UserOut, error) {
	userOut := &UserOut{}
	if err := mapper.AutoMapper(user, userOut); err != nil {
		util.Log().Error("mapper user -> userOut failed", err)
		return nil, err
	} else {
		model.DB.Preload("Roles").Find(&user)

		userOut.ID = user.ID

		roles := make([]string, len(user.Roles))
		for i, role := range user.Roles {
			roles[i] = role.Name
		}
		userOut.Role = roles

		return userOut, nil
	}
}
