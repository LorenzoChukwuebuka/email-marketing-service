package controllers

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
	"email-marketing-service/api/utils"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world")
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var reqdata *model.User

	utils.DecodeRequestBody(r, &reqdata)

	validate := validator.New()

	err := validate.Struct(reqdata)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		// Construct a response with validation errors
		errorMap := make(map[string]string)
		for _, e := range validationErrors {
			errorMap[e.Field()] = e.Tag()
		}

		errorResponse := map[string]interface{}{"errors": errorMap}

		utils.ErrorResponse(w, errorResponse)

		return
	}

	//hash password

	password, _ := bcrypt.GenerateFromPassword([]byte(reqdata.Password), 14)

	reqdata.Password = password
	reqdata.UUID = uuid.New().String()

	userCreate, err := repository.CreateUser(reqdata)

	if err != nil {
		// Return a generic error response for any repository error
		utils.ErrorResponse(w, err.Error())
		return
	}

	// //send mail

	utils.SuccessResponse(w, 200, userCreate)

}

func VerifyUser(w http.ResponseWriter, r *http.Request) {

}

func ForgetPassword(w http.ResponseWriter, r *http.Request) {

}
