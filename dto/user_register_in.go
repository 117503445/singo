package dto

type UserRegisterIn struct {
	UserName string `form:"username" json:"username" binding:"required,min=5,max=30"`
	Password string `form:"password" json:"password" binding:"required,min=4,max=40"`
	Avatar   string `json:"avatar" gorm:"size:1000"`
}
