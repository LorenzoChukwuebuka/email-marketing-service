package services

import (
	"email-marketing-service/api/custom"
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
	"email-marketing-service/api/utils"
	"fmt"
	"time"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var err error

var validate = validator.New()

func CreateUser(d *model.User) (*model.User, error) {

	err = validate.Struct(d)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		// Construct a response with validation errors
		errorMap := make(map[string]string)
		for _, e := range validationErrors {
			errorMap[e.Field()] = e.Tag()
		}
		errorResponse := map[string]interface{}{"errors": errorMap}

		return nil, fmt.Errorf("validation errors: %v", errorResponse)

	}
	password, _ := bcrypt.GenerateFromPassword([]byte(d.Password), 14)

	d.Password = password
	d.UUID = uuid.New().String()

	//check if user already exists
	userExists, err := repository.CheckIfEmailAlreadyExists(d)

	if err != nil {
		return nil, err
	}

	fmt.Println(userExists)
	if userExists {
		return nil, fmt.Errorf("user already exists")
	}

	_, err = repository.CreateUser(d)

	if err != nil {
		return nil, err
	}

	otp := utils.GenerateOTP(8)

	//store otp with user details in db

	var otpData model.OTP

	otpData.UserId = d.ID
	otpData.Token = otp
	otpData.UUID = uuid.New().String()

	err = CreateOTP(&otpData)

	if err != nil {
		return nil, err
	}

	//send mail

	err = custom.SignUpMail(d.Email, d.UserName, otp)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func VerifyUser(d *model.OTP) error {
	err = validate.Struct(d)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		// Construct a response with validation errors
		errorMap := make(map[string]string)
		for _, e := range validationErrors {
			errorMap[e.Field()] = e.Tag()
		}
		errorResponse := map[string]interface{}{"errors": errorMap}

		return fmt.Errorf("validation errors: %v", errorResponse)

	}

	//check if token exists in the otp table if yes, retrieve the records

	otpData, err := RetrieveOTP(d)

	if err != nil {
		return err
	}

	var userModel *model.User

	userModel.Verified = true
	userModel.ID = otpData.UserId
	userModel.VerifiedAt = time.Now()

	

	return nil
}

func Login() {

}

func ForgetPassword() {}

func ResetPassword() {}
