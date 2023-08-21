package services

import (
	"email-marketing-service/api/custom"
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
	"email-marketing-service/api/utils"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var err error

func CreateUser(d *model.User) (*model.User, error) {
	validate := validator.New()

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

	//hash password

	password, _ := bcrypt.GenerateFromPassword([]byte(d.Password), 14)

	d.Password = password
	d.UUID = uuid.New().String()

	//check if user already exists

	// userExists, err := repository.CheckIfEmailAlreadyExists(d)

	// if err != nil {
	// 	return nil, err
	// }
	// if userExists {
	// 	return nil, fmt.Errorf("user already exists")
	// }

	_, err = repository.CreateUser(d)

	if err != nil {
		return nil, err
	}

	otp := utils.GenerateOTP(8)

	//send mail

	err = custom.SignUpMail(d.Email, d.UserName, otp)
	if err != nil {
		// Handle the error from sending the mail
		return nil, err

	}
	return d, nil
}
