package dto

type UserSignUpRequest struct {
	FullName string `json:"fullname" validate:"required"`
	Company  string `json:"company"`
	Email    string `json:"email" validate:"required,email" `
	Password string `json:"password"`
	Verified bool   `json:"verified"`
	GoogleID string `json:"google_id"`
}

type UsersSignUpResponse struct {
	FullName string `json:"fullname" validate:"required"`
	Company  string `json:"company"`
	Email    string `json:"email" validate:"required,email" `
	Verified bool   `json:"verified"`
}

type VerifyUserRequest struct {
	UserID string `json:"user_id"`
	Token  string `json:"token" validate:"required"`
}

type EditUserDTO struct {
	UUID        string `json:"uuid"`
	FullName    string `json:"fullname"`
	Company     string `json:"company"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phonenumber"`
}

type Login struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password"`
	GoogleID string `json:"google_id"`
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
	UserId    string  `json:"user_id"`
	Device    *string `json:"device"`
	IPAddress *string `json:"ip_address"`
	Browser   *string `json:"browser"`
}

type RefreshAccessToken struct {
	RefreshToken string `json:"refresh_token"`
}

type ResendOTPRequest struct {
	UserId   string `json:"user_id" validated:"required" `
	Username string `json:"username" validated:"required"`
	Email    string `json:"email" validated:"required"`
	OTPType  string `json:"otp_type" validated:"required"`
}

