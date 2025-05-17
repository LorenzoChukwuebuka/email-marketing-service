package common

import (
	"errors"
)

var (
	ErrDecodingRequestBody     = errors.New("Error decoding request body")
	ErrPasswordHashingFailed   = errors.New("Error hashing password")
	ErrValidatingRequest       = errors.New("Error validating request")
	ErrFetchingUser            = errors.New("Error Fetching User")
	ErrUserAlreadyExists       = errors.New("User already exists")
	ErrUserNotFound            = errors.New("User not found")
	ErrCreatingUser            = errors.New("error creating user")
	ErrCreatingOTP             = errors.New("error creating OTP token")
	ErrFetchingOTP             = errors.New("error creating otp")
	ErrVerificationCodeExpired = errors.New("verification code has expired")
	ErrDeletingOTP             = errors.New("error deleting otp")
	ErrInvalidUUID = errors.New("error parsing uuid")
)
