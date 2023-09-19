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

type UserService struct {
	userRepository *repository.UserRepository
	otpService     *OTPService
}

func NewUserService(userRepo *repository.UserRepository, otpSvc *OTPService) *UserService {
	return &UserService{
		userRepository: userRepo,
		otpService:     otpSvc,
	}
}

// CreateUser creates a new user, sends an OTP email, and stores OTP data.
func (s *UserService) CreateUser(d *model.User) (*model.User, error) {

	if err := utils.ValidateData(d); err != nil {
		return nil, err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(d.Password), 14)

	d.Password = password
	d.UUID = uuid.New().String()

	// Check if user already exists.
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

	// Store OTP with user details in the database.
	otpData := &model.OTP{
		UserId: d.ID,
		Token:  otp,
		UUID:   uuid.New().String(),
	}
	if err := s.otpService.CreateOTP(otpData); err != nil {
		return nil, err
	}

	// Send mail.
	if err := custom.SignUpMail(d.Email, d.UserName, otp); err != nil {
		return nil, err
	}

	return d, nil
}

// VerifyUser verifies a user account using OTP.
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
	userModel.VerifiedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	if err = s.userRepository.VerifyUserAccount(&userModel); err != nil {
		return err
	}

	//delete otp from the database

	if err = otpService.DeleteOTP(otpData.Id); err != nil {
		return err
	}

	//maybe send onboarding mail to them I don't know

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
	err = bcrypt.CompareHashAndPassword(userDetails.Password, []byte(d.Password))
	if err != nil {
		return nil, fmt.Errorf("passwords do not match")
	}

	// Generate JWT token
	token, err := utils.JWTEncode(userDetails.ID, userDetails.UUID, userDetails.UserName, userDetails.Email)
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
		UUID:   uuid.New().String(),
	}

	otpService := s.otpService

	if err = otpService.CreateOTP(otpData); err != nil {
		return err
	}

	if err = custom.ResetPasswordMail(d.Email, userDetails.UserName, otp); err != nil {
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
		Password: password,
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
	if err := bcrypt.CompareHashAndPassword([]byte(userData.Password), d.OldPassword); err != nil {
		return fmt.Errorf("password does not match the records")
	}

	//hash password if it passes test
	password, _ := bcrypt.GenerateFromPassword([]byte(d.NewPassword), 14)

	data.Password = password

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
