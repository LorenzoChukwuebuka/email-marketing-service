package dto

type OTP struct {
	UserId int    `json:"user_id"`
	Token  string `json:"token" validated:"required"`
}

type ResendOTP struct {
	UserId   string `json:"user_id" validated:"required" `
	Username string `json:"username" validated:"required"`
	Email    string `json:"email" validated:"required"`
	OTPType  string `json:"otp_type" validated:"required"`
}
