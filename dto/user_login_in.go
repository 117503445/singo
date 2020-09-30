package dto

type UserLoginIn struct {
	UserName string `json:"username" validate:"required,min=5,max=30"`
	Password string `json:"password" validate:"required,min=4,max=40"`
}
