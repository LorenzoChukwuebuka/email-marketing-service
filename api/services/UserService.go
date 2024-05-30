package services

import (
	"email-marketing-service/api/custom"
	"email-marketing-service/api/dto"
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
	"email-marketing-service/api/utils"
	"fmt"
	"strings"
	"time"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository   *repository.UserRepository
	otpService       *OTPService
	PlanRepo         *repository.PlanRepository
	SubscriptionRepo *repository.SubscriptionRepository
	BillingRepo      *repository.BillingRepository
}

var (
	mail = &custom.Mail{}
)

func NewUserService(userRepo *repository.UserRepository, otpSvc *OTPService, planRepo *repository.PlanRepository, subscriptionRepo *repository.SubscriptionRepository, billingRepo *repository.BillingRepository) *UserService {
	return &UserService{
		userRepository:   userRepo,
		otpService:       otpSvc,
		PlanRepo:         planRepo,
		SubscriptionRepo: subscriptionRepo,
		BillingRepo:      billingRepo,
	}
}

func (s *UserService) CreateUser(d *dto.User) (map[string]interface{}, error) {

	if err := utils.ValidateData(d); err != nil {
		return nil, err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(d.Password), 14)
	d.Password = string(password)

	usermodel := &model.User{
		FullName: d.FullName,
		Company:  d.Company,
		Email:    d.Email,
		Password: d.Password,
		Verified: false,
	}

	usermodel.UUID = uuid.New().String()

	userExists, err := s.userRepository.CheckIfEmailAlreadyExists(usermodel)
	if err != nil {
		return nil, err
	}
	if userExists {
		return nil, fmt.Errorf("user already exists")
	}

	if _, err := s.userRepository.CreateUser(usermodel); err != nil {
		return nil, err
	}

	otp := utils.GenerateOTP(8)

	otpData := &model.OTP{
		UserId: usermodel.UUID,
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
		"userId":  usermodel.UUID,
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
	userModel.UUID = otpData.UserId
	userModel.VerifiedAt = time.Now()

	userId, err := s.userRepository.VerifyUserAccount(&userModel)

	if err != nil {
		return fmt.Errorf("unable to verify user :%w", err)
	}

	//delete otp from the database
	if err = otpService.DeleteOTP(otpData.Id); err != nil {
		return fmt.Errorf("unable to delete otp:%w", err)
	}

	err = s.createUserBasicPlan(userId)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) createUserBasicPlan(userId int) error {
	plans, err := s.PlanRepo.GetAllPlans()

	if err != nil {
		return fmt.Errorf("failed to fetch plans: %w", err)
	}

	var basicPlan *model.PlanResponse
	for _, plan := range plans {
		if strings.ToLower(plan.PlanName) == "basic" || strings.ToLower(plan.PlanName) == "free" {
			basicPlan = &plan
			break
		}
	}

	if basicPlan == nil {
		return fmt.Errorf("no basic or free plan found")
	}

	transactionId := uuid.New().String()

	billing := &model.Billing{
		UUID:          uuid.New().String(),
		UserId:        userId,
		AmountPaid:    basicPlan.Price,
		PlanId:        basicPlan.ID,
		PaymentMethod: "",
		Duration:      "",
		Email:         "",
		ExpiryDate:    time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
		Reference:     "",
		TransactionId: transactionId,
		CreatedAt:     time.Now(),
	}

	bill, err := s.BillingRepo.CreateBilling(billing)

	if err != nil {
		return fmt.Errorf("failed to create billing for user %d: %w", userId, err)
	}

	subscription := &model.Subscription{
		UserId:        userId,
		PlanId:        basicPlan.ID,
		Expired:       false,
		StartDate:     time.Now(),
		EndDate:       time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC), // 1 month subscription
		TransactionId: transactionId,
		Cancelled:     false,
		CreatedAt:     time.Now(),
		PaymentId:     bill.Id,
	}

	err = s.SubscriptionRepo.CreateSubscription(subscription)
	if err != nil {
		return fmt.Errorf("failed to create subscription for user %d: %w", userId, err)
	}

	//send onboarding mail

	return nil
}

func (s *UserService) Login(d *dto.Login) (map[string]interface{}, error) {
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
	token, err := utils.JWTEncode(userDetails.UUID, userDetails.UUID, userDetails.FullName, userDetails.Email)
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

func (s *UserService) ForgetPassword(d *dto.ForgetPassword) error {
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
		UserId: userDetails.UUID,
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

func (s *UserService) ResetPassword(d *dto.ResetPassword) error {
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
		UUID:     otpData.UserId,
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

func (s *UserService) ChangePassword(userId int, d *dto.ChangePassword) error {

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

func (s *UserService) ResendOTP(d *dto.ResendOTP) error {
	if err := utils.ValidateData(d); err != nil {
		return err
	}

	otp := utils.GenerateOTP(8)

	// Store OTP with user details in the database.
	otpData := &model.OTP{
		UserId: d.UserId,
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
