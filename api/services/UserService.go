package services

import (
	"email-marketing-service/api/custom"
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
	"email-marketing-service/api/utils"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type UserService struct {
	userRepository *repository.UserRepository
	otpService     *OTPService
}

var (
	mail = &custom.Mail{}
)

func NewUserService(userRepo *repository.UserRepository, otpSvc *OTPService) *UserService {
	return &UserService{
		userRepository: userRepo,
		otpService:     otpSvc,
	}
}

func (s *UserService) CreateUser(d *model.User) (map[string]interface{}, error) {

	if err := utils.ValidateData(d); err != nil {
		return nil, err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(d.Password), 14)
	d.Password = string(password)
	d.UUID = uuid.New().String()

	userExists, err := s.userRepository.CheckIfEmailAlreadyExists(d)
	if err != nil {
		return nil, err
	}
	if userExists {
		return nil, fmt.Errorf("user already exists")
	}

	if _, err := s.userRepository.CreateUser(d); err != nil {
		return nil, err
	}

	otp := utils.GenerateOTP(8)

	otpData := &model.OTP{
		UserId: d.ID,
		Token:  otp,
	}
	if err := s.otpService.CreateOTP(otpData); err != nil {
		return nil, err
	}

	if err := mail.SignUpMail(d.Email, d.FullName, otp); err != nil {
		return nil, err
	}

	successMap := map[string]interface{}{
		"message": "Account created successfully. Kindly verify your account",
		"userId":  d.UUID,
	}

	return successMap, nil
}

func (s *UserService) VerifyUser(d *model.OTP) error {

	if err := utils.ValidateData(d); err != nil {
		return err
	}

	//check if token exists in the otp table if yes, retrieve the records
	otpService := s.otpService
	otpData, err := otpService.RetrieveOTP(d)

	if err != nil {
		return err
	}

	var userModel model.User

	userModel.Verified = true
	userModel.ID = otpData.UserId
	userModel.VerifiedAt = time.Now()

	if err = s.userRepository.VerifyUserAccount(&userModel); err != nil {
		return err
	}

	//delete otp from the database
	if err = otpService.DeleteOTP(otpData.Id); err != nil {
		return err
	}

	// Todo

	// 1. get all the plans and check if there is a basic or free plan

	//2. automatically create a basic/free subscription plan for them

	//3. Test and make sure that they can send mails... with their key of course

	return nil
}

func (s *UserService) Login(d *model.LoginModel) (map[string]interface{}, error) {
	if err := utils.ValidateData(d); err != nil {
		return nil, err
	}

	user := &model.User{
		Email:    d.Email,
		Password: d.Password,
	}
	userDetails, err := s.userRepository.Login(user)

	if err != nil {
		return nil, fmt.Errorf("error during login: %w", err)
	}

	if !userDetails.Verified {
		return nil, fmt.Errorf("this account has not been verified")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(userDetails.Password), []byte(d.Password))
	if err != nil {
		return nil, fmt.Errorf("passwords do not match")
	}

	//Generate JWT token
	token, err := utils.JWTEncode(userDetails.ID, userDetails.UUID, userDetails.FullName, userDetails.Email)
	if err != nil {
		return nil, fmt.Errorf("error generating JWT token: %w", err)
	}

	// Marshal the user details back to JSON
	successMap := map[string]interface{}{
		"status":  "login successful",
		"token":   token,
		"details": userDetails,
	}

	return successMap, nil
}

func (s *UserService) ForgetPassword(d *model.ForgetPassword) error {
	if err := utils.ValidateData(d); err != nil {
		return err
	}

	userEmail := &model.User{
		Email: d.Email,
	}

	//check if email exists in db
	userExists, err := s.userRepository.CheckIfEmailAlreadyExists(userEmail)

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
	userDetails, err := s.userRepository.FindUserByEmail(email)

	if err != nil {
		return err
	}

	//generate token
	otp := utils.GenerateOTP(8)

	otpData := &model.OTP{
		UserId: userDetails.ID,
		Token:  otp,
	}

	otpService := s.otpService

	if err = otpService.CreateOTP(otpData); err != nil {
		return err
	}

	if err = mail.ResetPasswordMail(d.Email, userDetails.FullName, otp); err != nil {
		return err
	}

	return nil
}

func (s *UserService) ResetPassword(d *model.ResetPassword) error {
	if err := utils.ValidateData(d); err != nil {
		return err
	}

	data := &model.OTP{
		Token: d.Token,
	}

	otpService := s.otpService

	otpData, err := otpService.RetrieveOTP(data)

	if err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(d.Password), 14)

	user := &model.User{
		ID:       otpData.UserId,
		Password: string(password),
	}

	if err = s.userRepository.ResetPassword(user); err != nil {
		return err
	}

	//delete otp from the database

	if err = otpService.DeleteOTP(otpData.Id); err != nil {
		return err
	}

	return nil
}

func (s *UserService) ChangePassword(userId int, d *model.ChangePassword) error {

	if err := utils.ValidateData(d); err != nil {
		return err
	}

	data := &model.User{
		ID: userId,
	}

	userData, err := s.userRepository.FindUserById(data)

	if err != nil {
		return err
	}

	//compare password
	if err := bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(d.OldPassword)); err != nil {
		return fmt.Errorf("password does not match the records")
	}

	//hash password if it passes test
	password, _ := bcrypt.GenerateFromPassword([]byte(d.NewPassword), 14)

	data.Password = string(password)

	if err := s.userRepository.ChangeUserPassword(data); err != nil {
		return err
	}

	return nil
}

func (s *UserService) EditUser(id int, d *model.User) error {
	if err := s.userRepository.UpdateUserRecords(d); err != nil {
		return err
	}

	return nil
}

func (s *UserService) ResendOTP(d *model.ResendOTP) error {
	if err := utils.ValidateData(d); err != nil {
		return err
	}

	otp := utils.GenerateOTP(8)

	num, err := strconv.Atoi(d.UserId)

	if err != nil {
		return err
	}

	// Store OTP with user details in the database.
	otpData := &model.OTP{
		UserId: num,
		Token:  otp,
	}
	if err := s.otpService.CreateOTP(otpData); err != nil {
		return err
	}

	// Send mail.

	if d.OTPType == "emailVerify" {

		if err := mail.SignUpMail(d.Email, d.Username, otp); err != nil {
			return err
		}
	} else if d.OTPType == "passwordReset" {
		if err := mail.ResetPasswordMail(d.Email, d.Username, otp); err != nil {
			return err
		}

	}

	return nil
}
