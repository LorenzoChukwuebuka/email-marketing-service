package services

import (
	"database/sql"
	"email-marketing-service/api/custom"
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
	"email-marketing-service/api/utils"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
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

	fmt.Println(d)
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

	var userModel model.User

	userModel.Verified = true
	userModel.ID = otpData.UserId
	userModel.VerifiedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	err = repository.VerifyUserAccount(&userModel)

	if err != nil {
		return err
	}

	//delete otp from the database

	err = DeleteOTP(otpData.Id)

	if err != nil {
		return err
	}

	//maybe send onboarding mail to them I don't know

	return nil
}

func Login(d *model.LoginModel) (map[string]interface{}, error) {
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

	var user model.User

	user.Email = d.Email
	user.Password = d.Password
	userDetails, err := repository.Login(&user)

	if err != nil {
		return nil, err
	}

	//compare password
	err = bcrypt.CompareHashAndPassword(userDetails.Password, []byte(d.Password))

	if err != nil {
		return nil, fmt.Errorf("passwords do not match:%w", err)
	}

	token, err := utils.JWTEncode(userDetails.ID, userDetails.UserName, userDetails.Email)

	if err != nil {
		return nil, err
	}

	successMap := map[string]interface{}{
		"status": "login successful",
		"token":  token,
	}

	return successMap, nil
}

func ForgetPassword() {}

func ResetPassword() {}
