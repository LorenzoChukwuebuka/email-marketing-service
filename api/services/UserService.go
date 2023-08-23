package services

import (
	"database/sql"
	"email-marketing-service/api/custom"
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
	"email-marketing-service/api/utils"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var err error

func CreateUser(d *model.User) (*model.User, error) {

	err = utils.ValidateData(d)

	if err != nil {
		return nil, err
	}
	password, _ := bcrypt.GenerateFromPassword([]byte(d.Password), 14)

	d.Password = password
	d.UUID = uuid.New().String()

	//check if user already exists
	userExists, err := repository.CheckIfEmailAlreadyExists(d)

	if err != nil {
		return nil, err
	}

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
	err = utils.ValidateData(d)

	if err != nil {
		return err
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

func Login(d *model.LoginModel) (map[string]string, error) {
	err = utils.ValidateData(d)

	if err != nil {
		return nil, err
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

	successMap := map[string]string{
		"status": "login successful",
		"token":  token,
	}

	return successMap, nil
}

func ForgetPassword(d *model.ForgetPassword) error {
	err = utils.ValidateData(d)

	if err != nil {
		return err
	}

	userEmail := &model.User{
		Email: d.Email,
	}

	//check if email exists in db
	userExists, err := repository.CheckIfEmailAlreadyExists(userEmail)

	if err != nil {
		return err
	}

	if !userExists {
		return nil
	}

	email := &model.User{
		Email: d.Email,
	}

	//get username and id and append them to the email and otp services
	userDetails, err := repository.FindUserByEmail(email)

	if err != nil {
		return err
	}

	//generate token
	otp := utils.GenerateOTP(8)

	otpData := &model.OTP{
		UserId: userDetails.ID,
		Token:  otp,
		UUID:   uuid.New().String(),
	}

	err = CreateOTP(otpData)

	if err != nil {
		return err
	}

	err = custom.ResetPasswordMail(d.Email, userDetails.UserName, otp)

	if err != nil {
		return err
	}

	return nil
}

func ResetPassword(d *model.ResetPassword) error {
	return nil
}


