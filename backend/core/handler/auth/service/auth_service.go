package service

import (
	"context"
	"database/sql"
	"email-marketing-service/core/handler/auth/dto"
	"email-marketing-service/internal/common"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/enums"
	"email-marketing-service/internal/helper"
	"email-marketing-service/internal/mailer"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Service struct {
	store db.Store
}

var (
	otpLength = 15
)

func NewAuthService(store db.Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) SignUp(ctx context.Context, req *dto.UserSignUpRequest) (any, error) {
	if err := helper.ValidateData(req); err != nil {
		return nil, errors.Join(common.ErrValidatingRequest, err)
	}
	// Check if a user already exists with this email
	_, err := s.store.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, common.ErrUserAlreadyExists
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, common.ErrFetchingUser
	}

	companyNullString := sql.NullString{
		String: req.Company,
		Valid:  req.Company != "",
	}

	hashedPass, err := common.HashPassword(req.Password)
	if err != nil {
		return nil, common.ErrPasswordHashingFailed
	}

	token := helper.GenerateOTP(otpLength)

	err = s.store.ExecTx(ctx, func(q *db.Queries) error {
		// Create company
		company, err := q.CreateCompany(ctx, companyNullString)
		if err != nil {
			return fmt.Errorf("error creating company: %w", err)
		}

		// Create user
		userParams := db.CreateUserParams{
			Fullname:  req.FullName,
			CompanyID: company.ID,
			Email:     req.Email,
			Password: sql.NullString{
				String: hashedPass,
				Valid:  true,
			},
			Verified: false,
		}
		user, err := q.CreateUser(ctx, userParams)
		if err != nil {
			return fmt.Errorf("error creating user: %w", err)
		}

		//create smtp master key
		if err = s.createSMTPMasterKey(ctx, q, user.ID, user.CompanyID, user.Email); err != nil {
			return err
		}

		// Create OTP
		_, err = q.CreateOTP(ctx, db.CreateOTPParams{
			UserID:    user.ID,
			Token:     token,
			ExpiresAt: sql.NullTime{Time: time.Now().Add(1 * time.Hour), Valid: true},
		})
		if err != nil {
			return fmt.Errorf("error creating OTP: %w", err)
		}

		if err = s.sendCreateUserMail(user, token); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return dto.UsersSignUpResponse{FullName: req.FullName, Company: req.Company, Email: req.Email, Verified: false}, nil
}

func (s *Service) sendCreateUserMail(user db.User, token string) error {
	errChan := make(chan error)
	go func() {
		if err := mailer.NewEmailService().SignUpMail(user.Email, user.Fullname, user.ID, token); err != nil {
			errChan <- fmt.Errorf("%w: %v", common.ErrSendingEmail, err)
		}
		close(errChan)
	}()
	if err := <-errChan; err != nil {
		return err
	}
	return nil
}

func (s *Service) createSMTPMasterKey(ctx context.Context, q *db.Queries, userId uuid.UUID, companyId uuid.UUID, userEmail string) error {
	_, err := q.CreateSMTPMasterKey(ctx, db.CreateSMTPMasterKeyParams{
		UserID:    userId,
		CompanyID: companyId,
		KeyName:   "Master",
		SmtpLogin: userEmail,
		Password:  helper.GenerateOTP(15),
		Status:    string(enums.STMPMasterKeyStatus("active")),
	})
	if err != nil {
		return errors.Join(common.ErrCreatingSMTPKey, err)
	}
	return nil
}

func (s *Service) VerifyUser(ctx context.Context, req *dto.VerifyUserRequest) error {
	if err := helper.ValidateData(req); err != nil {
		return errors.Join(common.ErrValidatingRequest, err)
	}

	err := s.store.ExecTx(ctx, func(q *db.Queries) error {
		getOtp, err := q.GetOTPByToken(ctx, req.Token)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("%w: OTP not found for token", common.ErrFetchingOTP)
			}
			return fmt.Errorf("%w: %v", common.ErrFetchingOTP, err)
		}

		if time.Now().UTC().After(getOtp.ExpiresAt.Time) {
			return common.ErrVerificationCodeExpired
		}

		err = q.VerifyUser(ctx, getOtp.UserID)
		if err != nil {
			return fmt.Errorf("error verifying user: %w", err)
		}

		if err = q.DeleteOTPById(ctx, getOtp.ID); err != nil {
			return fmt.Errorf("%w: %v", common.ErrDeletingOTP, err)
		}

		//get user
		user, err := q.GetUserByID(ctx, getOtp.UserID)
		if err != nil {
			return common.ErrFetchingUser
		}

		//create subscription
		if err = s.createUserSubscription(ctx, q, user); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (s *Service) createUserSubscription(ctx context.Context, q *db.Queries, user db.User) error {
	plan, err := s.store.GetPlanByName(ctx, "Free")
	if err != nil {
		return err
	}

	//create user subscription plan
	subscription, err := q.CreateSubscription(ctx, db.CreateSubscriptionParams{
		CompanyID:       user.CompanyID,
		PlanID:          plan.ID,
		Amount:          decimal.NewFromInt(0),
		BillingCycle:    sql.NullString{String: "monthly", Valid: true},
		TrialStartsAt:   sql.NullTime{Time: time.Now().UTC(), Valid: true},
		TrialEndsAt:     sql.NullTime{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true},
		NextBillingDate: sql.NullTime{Time: time.Now().Add(31 * 24 * time.Hour), Valid: true},
		AutoRenew:       sql.NullBool{Bool: false, Valid: true},
	})

	if err != nil {
		return common.ErrCreatingSubscription
	}

	//create payment
	payment, err := q.CreatePayment(ctx, db.CreatePaymentParams{
		CompanyID:            user.CompanyID,
		UserID:               user.ID,
		SubscriptionID:       subscription.ID,
		Amount:               decimal.NewFromInt(0),
		Currency:             sql.NullString{String: "NGN", Valid: true},
		PaymentMethod:        sql.NullString{String: "None", Valid: true},
		Status:               sql.NullString{String: "active", Valid: true},
		Notes:                sql.NullString{String: "everything was successful", Valid: true},
		TransactionReference: sql.NullString{String: "none", Valid: true},
		PaymentDate:          sql.NullTime{Time: time.Now(), Valid: true},
		BillingPeriodStart:   sql.NullTime{Time: time.Now(), Valid: true},
		BillingPeriodEnd:     sql.NullTime{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true},
	})

	//update payment integrity

	paymentHash, err := common.GeneratePaymentHash(payment.ID, user.ID, 0, subscription.ID)

	if err != nil {
		return err
	}

	err = q.UpdatePaymentHash(ctx, db.UpdatePaymentHashParams{
		IntegrityHash: sql.NullString{String: paymentHash, Valid: true},
		ID:            payment.ID,
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ResendEmail(ctx context.Context, req *dto.ResendOTPRequest) error {
	if err := helper.ValidateData(req); err != nil {
		return fmt.Errorf("invalid resend OTP data: %w", err)
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return common.ErrInvalidUUID
	}

	token := helper.GenerateOTP(otpLength)

	_, err = s.store.CreateOTP(ctx, db.CreateOTPParams{
		UserID:    userID,
		Token:     token,
		ExpiresAt: sql.NullTime{Time: time.Now().Add(1 * time.Hour), Valid: true},
	})
	if err != nil {
		return fmt.Errorf("error creating OTP: %w", err)
	}
	// switch req.OTPType {
	// case "emailVerify":
	// 	return mailer.SignUpMail(d.Email, d.Username, d.UserId, otp)
	// case "passwordReset":
	// 	return mailer.ResetPasswordMail(d.Email, d.Username, otp)
	// default:
	// 	return ErrInvalidOTPType
	// }

	fmt.Print("sending mail")
	return nil
}

func (s *Service) LoginUser(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse[db.User], error) {
	if err := helper.ValidateData(req); err != nil {
		return nil, errors.Join(common.ErrValidatingRequest, err)
	}
	// Check if a user exists
	user, err := s.store.GetUserByEmail(ctx, req.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.ErrUserNotFound
		}
		return nil, common.ErrFetchingUser
	}

	//check if user is verified
	if user.Blocked {
		return nil, common.ErrBlocked
	}

	if !user.Verified {
		return nil, common.ErrAccountNotVerified
	}

	err = common.CheckPassword(req.Password, user.Password.String)

	if err != nil {
		return nil, common.ErrCheckingPasswordHash
	}

	token, err := helper.GenerateAccessToken(user.ID.String(), user.ID.String(), user.Fullname, user.Email)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := helper.GenerateRefreshToken(user.ID.String(), user.ID.String(), user.Fullname, user.Email)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse[db.User]{
		Status:       "success",
		Token:        token,
		Details:      user,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) ForgetPassword(ctx context.Context, req *dto.ForgetPassword) error {
	if err := helper.ValidateData(req); err != nil {
		return fmt.Errorf("invalid forget password data: %w", err)
	}

	user, err := s.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return fmt.Errorf("error finding user: %w", err)
	}

	token := helper.GenerateOTP(otpLength)

	_, err = s.store.CreateOTP(ctx, db.CreateOTPParams{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: sql.NullTime{Time: time.Now().Add(1 * time.Hour), Valid: true},
	})
	if err != nil {
		return fmt.Errorf("error creating OTP: %w", err)
	}

	return mailer.NewEmailService().ResetPasswordMail(req.Email, user.Fullname, token)
}

func (s *Service) ResetPassword(ctx context.Context, req *dto.ResetPassword) error {
	if err := helper.ValidateData(req); err != nil {
		return fmt.Errorf("invalid reset password data: %w", err)
	}

	getOtp, err := s.store.GetOTPByToken(ctx, req.Token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: OTP not found for token", common.ErrFetchingOTP)
		}
		return fmt.Errorf("%w: %v", common.ErrFetchingOTP, err)
	}

	if time.Now().UTC().After(getOtp.ExpiresAt.Time) {
		return common.ErrVerificationCodeExpired
	}

	hashedPassword, err := common.HashPassword(req.Password)
	if err != nil {
		return common.ErrPasswordHashingFailed
	}

	hashPass := string(hashedPassword)

	if err := s.store.ResetUserPassword(ctx, db.ResetUserPasswordParams{
		Password: sql.NullString{String: hashPass, Valid: true},
		ID:       getOtp.UserID,
	}); err != nil {
		return fmt.Errorf("error resetting password: %w", err)
	}

	return s.store.DeleteOTPById(ctx, getOtp.ID)
}

func (s *Service) ChangePassword(ctx context.Context, userId string, req *dto.ChangePassword) error {
	if err := helper.ValidateData(req); err != nil {
		return fmt.Errorf("invalid change password data: %w", err)
	}

	userID, err := uuid.Parse(userId)
	if err != nil {
		return common.ErrInvalidUUID
	}

	user, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("error finding user: %w", err)
	}

	if err := common.CheckPassword(req.OldPassword, user.Password.String); err != nil {
		return common.ErrCheckingPasswordHash
	}

	hashedPassword, err := common.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	return s.store.ResetUserPassword(ctx, db.ResetUserPasswordParams{
		Password: sql.NullString{String: hashedPassword, Valid: true},
		ID:       user.ID,
	})

}
