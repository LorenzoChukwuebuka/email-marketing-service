package services

import (
	"email-marketing-service/api/v1/custom"
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	adminrepository "email-marketing-service/api/v1/repository/admin"
	"email-marketing-service/api/v1/utils"
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
	ErrCreatingSMTPKey      = errors.New("error creating smtp key")
	ErrBlocked              = errors.New("your account has been blocked. Kindly contact the admin")
)

const (
	bcryptCost     = 14
	otpLength      = 20
	successMessage = "Account created successfully. Kindly verify your account"
	freePlanName   = "free"
)

var (
	mailer = custom.NewEmailService()
	config = utils.LoadEnv()
)

type UserService struct {
	userRepo              *repository.UserRepository
	otpService            *OTPService
	planRepo              *repository.PlanRepository
	subscriptionRepo      *repository.SubscriptionRepository
	billingRepo           *repository.BillingRepository
	MailUsageRepo         *repository.MailUsageRepository
	smtpKeyRepo           *repository.SMTPKeyRepository
	SenderSVC             *SenderServices
	UserNotificationRepo  *repository.UserNotificationRepository
	AdminNotificationRepo *adminrepository.AdminNotificationRepository
}

func NewUserService(userRepo *repository.UserRepository,
	otpSvc *OTPService,
	planRepo *repository.PlanRepository,
	subscriptionRepo *repository.SubscriptionRepository,
	billingRepo *repository.BillingRepository,
	mailUsageRepo *repository.MailUsageRepository,
	smtpKeyRepo *repository.SMTPKeyRepository,
	sendersvc *SenderServices,
	userNotificationRepo *repository.UserNotificationRepository,
	adminNotificationRepo *adminrepository.AdminNotificationRepository,
) *UserService {
	return &UserService{
		userRepo:              userRepo,
		otpService:            otpSvc,
		planRepo:              planRepo,
		subscriptionRepo:      subscriptionRepo,
		billingRepo:           billingRepo,
		MailUsageRepo:         mailUsageRepo,
		smtpKeyRepo:           smtpKeyRepo,
		SenderSVC:             sendersvc,
		UserNotificationRepo:  userNotificationRepo,
		AdminNotificationRepo: adminNotificationRepo,
	}
}

func (s *UserService) CreateUser(d *dto.User) (map[string]interface{}, error) {
	if err := utils.ValidateData(d); err != nil {
		return nil, fmt.Errorf("invalid user data: %w", err)
	}

	tx := s.userRepo.DB.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			return
		}
	}()

	user := &model.User{
		UUID:     uuid.New().String(),
		FullName: d.FullName,
		Email:    strings.ToLower(d.Email),
	}

	if d.GoogleID != "" {
		user.GoogleID = d.GoogleID
		user.Verified = true
		user.VerifiedAt = utils.Ptr(time.Now())
	}

	if d.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(d.Password), bcryptCost)
		if err != nil {
			return nil, fmt.Errorf("error hashing password: %w", err)
		}
		user.Password = utils.Ptr(string(hashedPassword))
		user.Verified = false
		user.Company = utils.Ptr(d.Company)

		otp := utils.GenerateOTP(otpLength)

		if err := s.createAndSendOTP(user, otp); err != nil {
			tx.Rollback()
			return nil, err
		}

	}

	if err := s.checkUserExists(user); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := s.createUserInDB(user); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := s.createSMTPMasterKey(user.UUID, user.Email); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := s.createTempEmailForUser(user.UUID, user.Email); err != nil {
		tx.Rollback()
		return nil, err
	}

	notificationTitle := fmt.Sprintf("A new member %s has registered", d.FullName)
	link := fmt.Sprintf("/zen/dash/users/detail/%s", user.UUID)
	if err := utils.CreateAdminNotifications(s.AdminNotificationRepo, user.UUID, link, notificationTitle); err != nil {
		tx.Rollback()
		fmt.Printf("Failed to create notification: %v\n", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
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

	errChan := make(chan error)

	go func() {
		if err := mailer.SignUpMail(user.Email, user.FullName, user.UUID, otp); err != nil {
			errChan <- fmt.Errorf("%w: %v", ErrSendingEmail, err)
		}
		close(errChan)
	}()

	if err := <-errChan; err != nil {
		return err
	}

	return nil
}

func (s *UserService) createSuccessResponse(userID string) map[string]interface{} {
	return map[string]interface{}{
		"message": successMessage,
		"userId":  userID,
	}
}

func (s *UserService) createTempEmailForUser(userID string, UserEmail string) error {
	parts := strings.Split(UserEmail, "@")
	if len(parts) > 2 {
		return fmt.Errorf("invalid email format")
	}
	tempMail := parts[0] + "@" + config.DOMAIN
	tempModel := &model.UserTempEmail{UserId: userID, TemporaryEmail: tempMail}
	if err := s.userRepo.CreateTempEmail(tempModel); err != nil {
		return err
	}
	return nil
}

func (s *UserService) VerifyUser(d *model.OTP) error {
	if err := utils.ValidateData(d); err != nil {
		return fmt.Errorf("invalid OTP data: %w", err)
	}

	tx := s.userRepo.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	otpData, err := s.otpService.RetrieveOTP(d)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error retrieving OTP: %w", err)
	}

	htime := time.Now().UTC()

	user := &model.User{
		UUID:       otpData.UserId,
		Verified:   true,
		VerifiedAt: &htime,
	}

	userId, err := s.userRepo.VerifyUserAccount(user)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("unable to verify user: %w", err)
	}

	if err = s.otpService.DeleteOTP(int(otpData.ID)); err != nil {
		tx.Rollback()
		return fmt.Errorf("unable to delete OTP: %w", err)
	}

	if err = s.createSender(d.UserId); err != nil {
		tx.Rollback()
		return fmt.Errorf("unable to create sender")
	}

	return s.createUserBasicPlan(userId)
}

func (s *UserService) createUserBasicPlan(userId uint) error {
	// Start a database transaction
	tx := s.userRepo.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	basicPlan, err := s.findBasicPlan()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to find basic plan: %w", err)
	}
	transactionId := uuid.New().String()
	billing, err := s.createBilling(userId, basicPlan, transactionId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create billing: %w", err)
	}
	err = s.createSubscription(userId, basicPlan, transactionId, int(billing.ID))
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create subscription: %w", err)
	}
	// Commit the transaction if everything is successful
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (s *UserService) findBasicPlan() (*model.PlanResponse, error) {
	plans, err := s.planRepo.GetAllPlans()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch plans: %w", err)
	}

	for _, plan := range plans {
		if strings.ToLower(plan.PlanName) == freePlanName {
			return &plan, nil
		}
	}
	return nil, ErrNoBasicPlan
}

func (s *UserService) createBilling(userId uint, plan *model.PlanResponse, transactionId string) (*model.Billing, error) {
	billing := &model.Billing{
		UUID:          uuid.New().String(),
		UserId:        userId,
		AmountPaid:    plan.Price,
		PlanId:        uint(plan.ID),
		ExpiryDate:    time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
		TransactionId: transactionId,
	}

	bill, err := s.billingRepo.CreateBilling(billing)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCreatingBilling, err)
	}

	return bill, nil
}

func (s *UserService) createSubscription(userId uint, plan *model.PlanResponse, transactionId string, paymentId int) error {
	subscription := &model.Subscription{
		UUID:          uuid.New().String(),
		UserId:        userId,
		PlanId:        uint(plan.ID),
		PaymentId:     paymentId,
		TransactionId: transactionId,
		StartDate:     time.Now(),
		EndDate:       time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
	}

	sub, err := s.subscriptionRepo.CreateSubscription(subscription)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCreatingSubscription, err)
	}

	isPeriodDaily := strings.ToLower(plan.PlanName) == freePlanName
	_, err = s.MailUsageRepo.GetOrCreateCurrentMailUsageRecord(int(sub), plan.MailingLimit.LimitAmount, isPeriodDaily)
	if err != nil {
		return fmt.Errorf("failed to create mail usage record: %w", err)
	}

	return nil
}

func (s *UserService) createSMTPMasterKey(userId string, userEmail string) error {
	smtpkeyModel := &model.SMTPMasterKey{
		UUID:      uuid.New().String(),
		KeyName:   "Master",
		UserId:    userId,
		SMTPLogin: userEmail,
		Password:  utils.GenerateOTP(15),
		Status:    model.KeyStatus("active"),
	}
	err := s.smtpKeyRepo.CreateSMTPMasterKey(smtpkeyModel)
	if err != nil {
		return ErrCreatingSMTPKey
	}
	return nil
}

func (s *UserService) createSender(userId string) error {

	userModel := &model.User{UUID: userId}

	getUser, err := s.userRepo.FindUserById(userModel)

	if err != nil {
		return err
	}

	sender := &dto.SenderDTO{UserID: userId, Email: getUser.Email, Name: "my company"}

	s.SenderSVC.CreateSender(sender)
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

	if user.Blocked {
		return nil, ErrBlocked
	}

	if !user.Verified {
		return nil, ErrAccountNotVerified
	}

	var userPass *string

	if d.Password != "" {
		userPass = user.Password

		if err := bcrypt.CompareHashAndPassword([]byte(*userPass), []byte(d.Password)); err != nil {
			return nil, ErrInvalidCredentials
		}
	}

	token, err := utils.GenerateAccessToken(user.UUID, user.UUID, user.FullName, user.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGeneratingToken, err)
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.UUID, user.UUID, user.FullName, user.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGeneratingToken, err)
	}

	return map[string]interface{}{
		"status":        "login successful",
		"token":         token,
		"details":       user,
		"refresh_token": refreshToken,
	}, nil
}

func (s *UserService) GoogleLogin(d *dto.User) (map[string]interface{}, error) {
	user, err := s.userRepo.Login(&model.User{Email: strings.ToLower(d.Email), GoogleID: d.GoogleID})
	if err != nil {
		return nil, fmt.Errorf("error during login: %w", err)
	}

	if user.Blocked {
		return nil, ErrBlocked
	}

	// Generate access token
	token, err := utils.GenerateAccessToken(user.UUID, user.UUID, user.FullName, user.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGeneratingToken, err)
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.UUID, user.UUID, user.FullName, user.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGeneratingToken, err)
	}

	return map[string]interface{}{
		"status":        "login successful",
		"token":         token,
		"details":       user,
		"refresh_token": refreshToken,
	}, nil
}

func (s *UserService) CreateGoogleUser(d *dto.User) (map[string]interface{}, error) {
	tx := s.userRepo.DB.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			return
		}
	}()

	// Reuse most of your existing CreateUser logic
	userData, err := s.CreateUser(d)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// After successful creation, generate tokens
	user, err := s.userRepo.Login(&model.User{Email: strings.ToLower(d.Email)})
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error fetching user after creation: %w", err)
	}

	if err := s.createUserBasicPlan(user.ID); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create basic plan: %w", err)
	}

	if err := s.createSMTPMasterKey(user.UUID, user.Email); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create key: %w", err)
	}

	if err := s.createSender(user.UUID); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create sender: %w", err)
	}

	token, err := utils.GenerateAccessToken(user.UUID, user.UUID, user.FullName, user.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGeneratingToken, err)
	}

	refreshToken, err := utils.GenerateRefreshToken(user.UUID, user.UUID, user.FullName, user.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGeneratingToken, err)
	}

	return map[string]interface{}{
		"status":        "signup successful",
		"token":         token,
		"details":       user,
		"refresh_token": refreshToken,
		"user_data":     userData,
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

	hashPass := string(hashedPassword)

	if err := s.userRepo.ResetPassword(&model.User{UUID: otpData.UserId, Password: &hashPass}); err != nil {
		return fmt.Errorf("error resetting password: %w", err)
	}

	return s.otpService.DeleteOTP(int(otpData.ID))
}

func (s *UserService) ChangePassword(userId string, d *dto.ChangePassword) error {
	if err := utils.ValidateData(d); err != nil {
		return fmt.Errorf("invalid change password data: %w", err)
	}

	user, err := s.userRepo.FindUserById(&model.User{UUID: userId})
	if err != nil {
		return fmt.Errorf("error finding user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(d.OldPassword)); err != nil {
		return ErrInvalidCredentials
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(d.NewPassword), bcryptCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}
	haspas := string(hashedPassword)
	return s.userRepo.ChangeUserPassword(&model.User{UUID: userId, Password: &haspas})
}

func (s *UserService) EditUser(d *model.User) error {
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

	dailyMailCalc, err := s.MailUsageRepo.GetCurrentMailUsageRecord(currentSub.Id)
	if err != nil {
		return nil, fmt.Errorf("error getting daily mail calculation: %w", err)
	}

	// Check for nil pointers and provide default values
	planName := ""
	mailsPerDay := 0
	remainingMails := 0

	if currentSub != nil && currentSub.Plan != nil {
		planName = currentSub.Plan.PlanName

	}

	if dailyMailCalc != nil {
		remainingMails = dailyMailCalc.RemainingMails
		mailsPerDay = dailyMailCalc.LimitAmount
	}

	return map[string]interface{}{
		"plan":           planName,
		"mailsPerDay":    mailsPerDay,
		"remainingMails": remainingMails,
	}, nil
}

func (s *UserService) GetUserDetails(userId string) (*model.UserResponse, error) {

	userModel := &model.User{UUID: userId}

	userDetails, err := s.userRepo.FindUserById(userModel)

	if err != nil {
		return nil, err
	}

	return &userDetails, nil
}

func (s *UserService) GetAllNotifications(userId string) ([]model.UserNotificationResponse, error) {
	usernotifications, err := s.UserNotificationRepo.GetAllUserNotification(userId)

	if err != nil {
		return nil, err
	}

	return usernotifications, err
}

func (s *UserService) UpdateReadStatus(userId string) error {
	if err := s.UserNotificationRepo.UpdateReadStatus(userId); err != nil {
		return err
	}

	return nil
}

func (s *UserService) MarkUserForDeletion(userId string) error {
	if err := s.userRepo.MarkUserForDeletion(userId); err != nil {
		return err
	}
	return nil
}

func (s *UserService) CancelUserDeletion(userId string) error {
	if err := s.userRepo.CancelUserDeletion(userId); err != nil {
		return err
	}
	return nil
}

// func processPendingDeletions() {
//     users, err := adminRepo.GetPendingDeletions()
//     if err != nil {
//         log.Printf("Failed to get pending deletions: %v", err)
//         return
//     }

//     for _, user := range users {
//         if err := adminRepo.PermanentlyDeleteUser(user.UUID); err != nil {
//             log.Printf("Failed to delete user %s: %v", user.UUID, err)
//             continue
//         }
//     }
// }
