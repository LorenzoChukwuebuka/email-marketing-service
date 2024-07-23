package dto

type ContactDTO struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	UserId    string `json:"user_id"`
	GroupId   uint   `json:"group_id" gorm:"default:null"`
	From      string `json:"from"`
}

type EditContactDTO struct {
	UserId    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	From      string `json:"from"`
}

type ContactGroupDTO struct {
}
