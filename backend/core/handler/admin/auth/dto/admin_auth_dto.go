package dto

import(
	"time"
)

type AdminRequestDTO struct {
	FirstName  string `json:"firstname"`
	MiddleName string `json:"middlename"`
	LastName   string `json:"lastname"`
	Email      string  `json:"email"`
	Password   string  `json:"password"`
	Type       string  `json:"type"`
}

type AdminLoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AdminEmailLogDTO struct {
	Subject string `json:"subject" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type AdminLoginResponse[T any] struct {
    Status       string `json:"status"`
    Token        string `json:"token"`
    Details      T      `json:"details"`
    RefreshToken string `json:"refresh_token"`
    Type         string `json:"type"`
}

type RefreshAccessToken struct {
	RefreshToken string `json:"refresh_token"`
}


 
type AdminResponse struct {
	ID         string     `json:"id"`
	Firstname  string     `json:"firstname"`
	Middlename string     `json:"middlename,omitempty"`
	Lastname   string     `json:"lastname"`
	Email      string     `json:"email"`
	Type       string     `json:"type"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}


