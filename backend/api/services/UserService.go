package services

import (
	"email-marketing-service/api/custom"
	"email-marketing-service/api/dto"
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
	"email-marketing-service/api/utils"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

var (
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrCreatingUser         = errors.New("error creating user")
	ErrCreatingOTP          = errors.New("error creating OTP")
	ErrSendingEmail         = errors.New("error sending email")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrAccountNotVerified   = errors.New("this account has not been verified")
	ErrGeneratingToken      = errors.New("error generating JWT token")
	ErrNoBasicPlan          = errors.New("no basic or free plan found")
	ErrCreatingBilling      = errors.New("failed to create billing")
	ErrCreatingSubscription = errors.New("failed to create subscription")
	ErrInvalidOTPType       = errors.New("invalid OTP type")
)

const (
	bcryptCost     = 14
	otpLength      = 8
	successMessage = "Account created successfully. Kindly verify your account"
	basicPlanName  = "basic"
	freePlanName   = "free"
)

var (
	mailer = &custom.Mail{}
)

type UserService struct {
	userRepo          *repository.UserRepository
	otpService        *OTPService
	planRepo          *repository.PlanRepository
	subscriptionRepo  *repository.SubscriptionRepository
	billingRepo       *repository.BillingRepository
	dailyMailCalcRepo *repository.DailyMailCalcRepository
}

func NewUserService(userRepo *repository.UserRepository,
	otpSvc *OTPService,
	planRepo *repository.PlanRepository,
	subscriptionRepo *repository.SubscriptionRepository,
	billingRepo *repository.BillingRepository,
	dailyMailCalcRepo *repository.DailyMailCalcRepository) *UserService {
	return &UserService{
		userRepo:          userRepo,
		otpService:        otpSvc,
		planRepo:          planRepo,
		subscriptionRepo:  subscriptionRepo,
		billingRepo:       billingRepo,
		dailyMailCalcRepo: dailyMailCalcRepo,
	}
}

func (s *UserService) CreateUser(d *dto.User) (map[string]interface{}, error) {
	if err := utils.ValidateData(d); err != nil {
		return nil, fmt.Errorf("invalid user data: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(d.Password), bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	user := &model.User{
		UUID:     uuid.New().String(),
		FullName: d.FullName,
		Company:  d.Company,
		Email:    strings.ToLower(d.Email),
		Password: string(hashedPassword),
		Verified: false,
	}

	if err := s.checkUserExists(user); err != nil {
		return nil, err
	}

	if err := s.createUserInDB(user); err != nil {
		return nil, err
	}

	otp := utils.GenerateOTP(otpLength)
	if err := s.createAndSendOTP(user, otp); err != nil {
		return nil, err
	}

	return s.createSuccessResponse(user.UUID), nil
}

func (s *UserService) checkUserExists(user *model.User) error {
	exists, err := s.userRepo.CheckIfEmailAlreadyExists(user)
	if err != nil {
		return fmt.Errorf("error checking user existence: %w", err)
	}
	if exists {
		return ErrUserAlreadyExists
	}
	return nil
}

func (s *UserService) createUserInDB(user *model.User) error {
	_, err := s.userRepo.CreateUser(user)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCreatingUser, err)
	}
	return nil
}

func (s *UserService) createAndSendOTP(user *model.User, otp string) error {
	otpData := &model.OTP{
		UserId: user.UUID,
		Token:  otp,
	}
	if err := s.otpService.CreateOTP(otpData); err != nil {
		return fmt.Errorf("%w: %v", ErrCreatingOTP, err)
	}

	if err := mailer.SignUpMail(user.Email, user.FullName, user.UUID, otp); err != nil {
		return fmt.Errorf("%w: %v", ErrSendingEmail, err)
	}
	return nil
}

func (s *UserService) createSuccessResponse(userID string) map[string]interface{} {
	return map[string]interface{}{
		"message": successMessage,
		"userId":  userID,
	}
}

func (s *UserService) VerifyUser(d *model.OTP) error {
	if err := utils.ValidateData(d); err != nil {
		return fmt.Errorf("invalid OTP data: %w", err)
	}

	otpData, err := s.otpService.RetrieveOTP(d)
	if err != nil {
		return fmt.Errorf("error retrieving OTP: %w", err)
	}

	user := &model.User{
		UUID:       otpData.UserId,
		Verified:   true,
		VerifiedAt: time.Now(),
	}

	userId, err := s.userRepo.VerifyUserAccount(user)
	if err != nil {
		return fmt.Errorf("unable to verify user: %w", err)
	}

	if err = s.otpService.DeleteOTP(otpData.Id); err != nil {
		return fmt.Errorf("unable to delete OTP: %w", err)
	}

	return s.createUserBasicPlan(userId)
}

func (s *UserService) createUserBasicPlan(userId int) error {
	basicPlan, err := s.findBasicPlan()
	if err != nil {
		return err
	}

	transactionId := uuid.New().String()

	billing, err := s.createBilling(userId, basicPlan, transactionId)
	if err != nil {
		return err
	}

	return s.createSubscription(userId, basicPlan, transactionId, billing.Id)
}

func (s *UserService) findBasicPlan() (*model.PlanResponse, error) {
	plans, err := s.planRepo.GetAllPlans()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch plans: %w", err)
	}

	for _, plan := range plans {
		if strings.ToLower(plan.PlanName) == basicPlanName || strings.ToLower(plan.PlanName) == freePlanName {
			return &plan, nil
		}
	}

	return nil, ErrNoBasicPlan
}

func (s *UserService) createBilling(userId int, plan *model.PlanResponse, transactionId string) (*model.Billing, error) {
	billing := &model.Billing{
		UUID:          uuid.New().String(),
		UserId:        userId,
		AmountPaid:    plan.Price,
		PlanId:        plan.ID,
		ExpiryDate:    time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
		TransactionId: transactionId,
		CreatedAt:     time.Now(),
	}

	bill, err := s.billingRepo.CreateBilling(billing)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCreatingBilling, err)
	}

	return bill, nil
}

func (s *UserService) createSubscription(userId int, plan *model.PlanResponse, transactionId string, paymentId int) error {
	subscription := &model.Subscription{
		UserId:        userId,
		PlanId:        plan.ID,
		StartDate:     time.Now(),
		EndDate:       time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
		TransactionId: transactionId,
		CreatedAt:     time.Now(),
		PaymentId:     paymentId,
	}

	err := s.subscriptionRepo.CreateSubscription(subscription)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCreatingSubscription, err)
	}

	return nil
}

func (s *UserService) Login(d *dto.Login) (map[string]interface{}, error) {
	if err := utils.ValidateData(d); err != nil {
		return nil, fmt.Errorf("invalid login data: %w", err)
	}

	user, err := s.userRepo.Login(&model.User{Email: strings.ToLower(d.Email)})
	if err != nil {
		return nil, fmt.Errorf("error during login: %w", err)
	}

	if !user.Verified {
		return nil, ErrAccountNotVerified
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(d.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := utils.JWTEncode(user.UUID, user.UUID, user.FullName, user.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGeneratingToken, err)
	}

	return map[string]interface{}{
		"status":  "login successful",
		"token":   token,
		"details": user,
	}, nil
}

func (s *UserService) ForgetPassword(d *dto.ForgetPassword) error {
	if err := utils.ValidateData(d); err != nil {
		return fmt.Errorf("invalid forget password data: %w", err)
	}

	user, err := s.userRepo.FindUserByEmail(&model.User{Email: strings.ToLower(d.Email)})
	if err != nil {
		return fmt.Errorf("error finding user: %w", err)
	}

	otp := utils.GenerateOTP(otpLength)
	if err := s.otpService.CreateOTP(&model.OTP{UserId: user.UUID, Token: otp}); err != nil {
		return fmt.Errorf("%w: %v", ErrCreatingOTP, err)
	}

	return mailer.ResetPasswordMail(user.Email, user.FullName, otp)
}

func (s *UserService) ResetPassword(d *dto.ResetPassword) error {
	if err := utils.ValidateData(d); err != nil {
		return fmt.Errorf("invalid reset password data: %w", err)
	}

	otpData, err := s.otpService.RetrieveOTP(&model.OTP{Token: d.Token})
	if err != nil {
		return fmt.Errorf("error retrieving OTP: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(d.Password), bcryptCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	if err := s.userRepo.ResetPassword(&model.User{UUID: otpData.UserId, Password: string(hashedPassword)}); err != nil {
		return fmt.Errorf("error resetting password: %w", err)
	}

	return s.otpService.DeleteOTP(otpData.Id)
}

func (s *UserService) ChangePassword(userId int, d *dto.ChangePassword) error {
	if err := utils.ValidateData(d); err != nil {
		return fmt.Errorf("invalid change password data: %w", err)
	}

	user, err := s.userRepo.FindUserById(&model.User{ID: userId})
	if err != nil {
		return fmt.Errorf("error finding user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(d.OldPassword)); err != nil {
		return ErrInvalidCredentials
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(d.NewPassword), bcryptCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	return s.userRepo.ChangeUserPassword(&model.User{ID: userId, Password: string(hashedPassword)})
}

func (s *UserService) EditUser(id int, d *model.User) error {
	return s.userRepo.UpdateUserRecords(d)
}

func (s *UserService) ResendOTP(d *dto.ResendOTP) error {
	if err := utils.ValidateData(d); err != nil {
		return fmt.Errorf("invalid resend OTP data: %w", err)
	}

	otp := utils.GenerateOTP(otpLength)
	if err := s.otpService.CreateOTP(&model.OTP{UserId: d.UserId, Token: otp}); err != nil {
		return fmt.Errorf("%w: %v", ErrCreatingOTP, err)
	}

	switch d.OTPType {
	case "emailVerify":
		return mailer.SignUpMail(d.Email, d.Username, d.UserId, otp)
	case "passwordReset":
		return mailer.ResetPasswordMail(d.Email, d.Username, otp)
	default:
		return ErrInvalidOTPType
	}
}

func (s *UserService) GetUserCurrentRunningSubscriptionWithMailsRemaining(userId string) (map[string]interface{}, error) {
	user, err := s.userRepo.FindUserById(&model.User{UUID: userId})
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	currentSub, err := s.subscriptionRepo.GetUsersCurrentSubscription(user.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting current subscription: %w", err)
	}

	dailyMailCalc, err := s.dailyMailCalcRepo.GetUserActiveCalculation(currentSub.Id)
	if err != nil {
		return nil, fmt.Errorf("error getting daily mail calculation: %w", err)
	}

	return map[string]interface{}{
		"plan":           currentSub.Plan.PlanName,
		"mailsPerDay":    currentSub.Plan.NumberOfMailsPerDay,
		"remainingMails": dailyMailCalc.RemainingMails,
	}, nil
}
