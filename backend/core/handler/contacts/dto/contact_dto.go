package dto

type ContactRequestDTO struct {
	FirstName    string `json:"first_name" validate:"required"`
	LastName     string `json:"last_name" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	UserId       string `json:"user_id"`
	From         string `json:"from"`
	IsSubscribed bool   `json:"is_subscribed"`
	CompanyID    string `json:"company_id" validate:"required"`
}

type EditContactDTO struct {
	ContactId    string `json:"contact_id"`
	UserID       string `json:"user_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	From         string `json:"from"`
	IsSubscribed bool   `json:"is_subscribed"`
}

type ContactGroupDTO struct {
	UserId      string `json:"user_id" validate:"required"`
	GroupName   string `json:"group_name" valiate:"required"`
	Description string `json:"description"`
	CompanyID   string `json:"company_id" validate:"required"`
}

type AddContactsToGroupDTO struct {
	UserId    string `json:"user_id" validate:"required"`
	GroupId   string `json:"group_id" validate:"required"`
	ContactId string `json:"contact_id" validate:"required"`
}

type ToggleContactSubDTO struct {
	ContactId    string `json:"contact_id"`
	UserId       string `json:"user_id"`
	IsSubscribed bool   `json:"is_subscribed"`
}

type FetchContactDTO struct {
	UserId      string `json:"user_id"`
	CompanyID   string `json:"company_id"`
	Offset      int    `json:"offset"`
	Limit       int    `json:"limit"`
	SearchQuery string `json:"search_query"`
}


type FetchContactGroupDTO struct {
	UserId    string `json:"user_id"`
	CompanyID   string `json:"company_id"`
	Offset      int    `json:"offset"`
	Limit       int    `json:"limit"`
	SearchQuery string `json:"search_query"`
}
