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
	"sync"
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

	// Use a WaitGroup to wait for goroutines to finish.
	var wg sync.WaitGroup

	// Channel to collect errors from goroutines.
	errCh := make(chan error, 2)

	// Check if user already exists concurrently.
	wg.Add(1)
	go func() {
		defer wg.Done()

		//check if user already exists
		userExists, err := s.userRepository.CheckIfEmailAlreadyExists(d)

		if err != nil {
			errCh <- err
			return
		}

		if userExists {
			errCh <- fmt.Errorf("user already exists")
			return
		}

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if _, err := s.userRepository.CreateUser(d); err != nil {
			errCh <- err
			return
		}

		otp := utils.GenerateOTP(8)

		//store otp with user details in db

		otpData := &model.OTP{
			UserId: d.ID,
			Token:  otp,
			UUID:   uuid.New().String(),
		}
		if err := s.otpService.CreateOTP(otpData); err != nil {
			errCh <- err
			return
		}

		//send mail

		if err := custom.SignUpMail(d.Email, d.UserName, otp); err != nil {
			errCh <- err
			return
		}
	}()

	// Close the error channel when all goroutines are done.
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Check for errors from goroutines.
	for err := range errCh {
		if err != nil {
			return nil, err
		}
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
		return nil, err
	}

	//compare password
	if err = bcrypt.CompareHashAndPassword(userDetails.Password, []byte(d.Password)); err != nil {
		return nil, fmt.Errorf("passwords do not match:%w", err)
	}

	token, err := utils.JWTEncode(userDetails.ID, userDetails.UserName, userDetails.Email)

	if err != nil {
		return nil, err
	}

	//marshal the user details back to json

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

// func (s *UserService) ForgetPassword(d *model.ForgetPassword) error {
// 	if err := utils.ValidateData(d); err != nil {
// 		return err
// 	}

// 	userEmail := &model.User{
// 		Email: d.Email,
// 	}

// 	// Check if email exists in the database
// 	userExists, err := s.userRepository.CheckIfEmailAlreadyExists(userEmail)
// 	if err != nil {
// 		return err
// 	}

// 	if !userExists {
// 		return nil
// 	}

// 	email := &model.User{
// 		Email: d.Email,
// 	}

// 	// Use a WaitGroup to wait for goroutines to finish
// 	var wg sync.WaitGroup

// 	var userDetails *model.User
// 	var userDetailsErr error

// 	// Get username and id concurrently
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		userDetails, userDetailsErr = s.userRepository.FindUserByEmail(email)
// 	}()

// 	// Generate token concurrently
// 	otp := utils.GenerateOTP(8)
// 	otpData := &model.OTP{
// 		Token: otp,
// 		UUID:  uuid.New().String(),
// 	}

// 	otpService := s.otpService

// 	var createOTPError error

// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		createOTPError = otpService.CreateOTP(otpData)
// 	}()

// 	// Wait for both goroutines to finish
// 	wg.Wait()

// 	if userDetailsErr != nil {
// 		return userDetailsErr
// 	}

// 	if createOTPError != nil {
// 		return createOTPError
// 	}

// 	// Send reset password email concurrently
// 	var sendEmailError error

// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		sendEmailError = custom.ResetPasswordMail(d.Email, userDetails.UserName, otp)
// 	}()

// 	// Wait for the sendEmail goroutine to finish
// 	wg.Wait()

// 	if sendEmailError != nil {
// 		return sendEmailError
// 	}

// 	return nil
// }

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
