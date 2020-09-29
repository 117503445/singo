package model

import (
	"errors"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name string `json:"name"`
}

func ReadRoleByName(roleName string) (Role, error) {
	var role Role
	result := DB.Where("name = ?", roleName).First(&role)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return role, result.Error
	}

	return role, nil
}
