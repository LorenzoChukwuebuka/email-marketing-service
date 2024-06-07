package dto

type User struct {
	FullName string `json:"fullname" validate:"required"`
	Company  string `json:"company" validate:"required"`
	Email    string `json:"email" validate:"required,email" `
	Password string `json:"password" validate:"required" `
	Verified bool   `json:"verified"`
}

type Login struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ForgetPassword struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPassword struct {
	Token    string `json:"token" validated:"required"`
	Password string `json:"password" validate:"required"`
}

type ChangePassword struct {
	UserId      int    `json:"user_id" validated:"required"`
	OldPassword string `json:"old_password" validated:"required"`
	NewPassword string `json:"new_password" validated:"required"`
}

type UserSession struct {
	UserId    string    `json:"user_id"`
	Device    *string `json:"device"`
	IPAddress *string `json:"ip_address"`
	Browser   *string `json:"browser"`
}
