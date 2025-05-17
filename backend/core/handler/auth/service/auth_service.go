package service

import (
	"context"
	"database/sql"
	"email-marketing-service/core/handler/auth/dto"
	"email-marketing-service/internal/common"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/helper"
	"errors"
	"fmt"
	"github.com/google/uuid"
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

		// Create OTP
		_, err = q.CreateOTP(ctx, db.CreateOTPParams{
			UserID:    user.ID,
			Token:     token,
			ExpiresAt: sql.NullTime{Time: time.Now().Add(1 * time.Hour), Valid: true},
		})
		if err != nil {
			return fmt.Errorf("error creating OTP: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	//send email to the user
	//TODO:send email to user
	//TODO: create smtp master key

	return dto.UsersSignUpResponse{FullName: req.FullName, Company: req.Company, Email: req.Email, Verified: false}, nil
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

		return nil
	})

	return err
}

func (s *Service) ResendVerificationEmail(ctx context.Context, req *dto.ResendOTPRequest) error {
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

func (s *Service) LoginUser() {}
