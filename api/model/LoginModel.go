package model

type LoginModel struct {
	Email    string `json:"email" validate:"required,email"`
	Password []byte `json:"password" validate:"required"`
}
