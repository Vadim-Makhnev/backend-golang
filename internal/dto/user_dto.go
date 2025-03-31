package dto

type UserRegister struct {
	Username string `json:"username" validate:"required,min=5,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=7,max=20"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=7,max=20"`
}