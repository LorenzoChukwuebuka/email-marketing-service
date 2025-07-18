package common

import (
	"errors"
	"fmt"
	"runtime"
)

var (
	ErrDecodingRequestBody       = errors.New("error decoding request body")
	ErrPasswordHashingFailed     = errors.New("error hashing password")
	ErrValidatingRequest         = errors.New("error validating request")
	ErrFetchingUser              = errors.New("error Fetching User")
	ErrUserAlreadyExists         = errors.New("user already exists")
	ErrUserNotFound              = errors.New("user not found")
	ErrCreatingUser              = errors.New("error creating user")
	ErrCreatingOTP               = errors.New("error creating OTP token")
	ErrFetchingOTP               = errors.New("error creating otp")
	ErrVerificationCodeExpired   = errors.New("verification code has expired")
	ErrDeletingOTP               = errors.New("error deleting otp")
	ErrInvalidUUID               = errors.New("error parsing uuid")
	ErrFetchingAdmin             = errors.New("error encountered while fetching admin")
	ErrAdminNotFound             = errors.New("admin does not exist")
	ErrCheckingPasswordHash      = errors.New("password does not match")
	ErrCreatingSMTPKey           = errors.New("error creating smtp key")
	ErrBlocked                   = errors.New("your account has been blocked. Kindly contact the admin")
	ErrSendingEmail              = errors.New("error sending email")
	ErrAccountNotVerified        = errors.New("account has not been verified")
	ErrCreatingSubscription      = errors.New("error creating subscription")
	ErrInvalidOTPType            = errors.New("invalid otp type")
	ErrFetchingSubscription      = errors.New("error fetching company's subscription")
	ErrParsingFile               = errors.New("error parsing file")
	ErrRetrievingFile            = errors.New("error retrieving file")
	ErrUpdatingRecord            = errors.New("error updating record")
	ErrDeletingRecord            = errors.New("error deleting record")
	ErrCreatingRecord            = errors.New("error creating record")
	ErrFetchingRecord            = errors.New("error fetching record")
	ErrFetchingCount             = errors.New("error fetching counts")
	ErrRecordExists              = errors.New("record already exists")
	ErrPaymentMethodNotSupported = errors.New("payment method not supported")
	ErrRecordNotFound            = errors.New("record not found")
	ErrUserPendingDeletion       = errors.New("user is pending deletion")
	ErrUserDeleted               = errors.New("user has been deleted")
	ErrUserBlocked               = errors.New("user has been blocked")
	ErrUserAlreadyVerified       = errors.New("user has already been verified")
)

func TraceError(err error) error {
	if err == nil {
		return nil
	}
	_, file, line, _ := runtime.Caller(1)
	return fmt.Errorf("%s:%d: %w", file, line, err)
}
