package service

import (
	"context"
	"database/sql"
	"email-marketing-service/core/handler/senders/dto"
	"email-marketing-service/core/handler/senders/mapper"
	"email-marketing-service/internal/common"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/helper"
	"email-marketing-service/internal/mailer"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net"
	"strings"
	"time"
)

type SenderService struct {
	store db.Store
}

func NewSenderService(store db.Store) *SenderService {
	return &SenderService{
		store: store,
	}
}

func (s *SenderService) CreateSender(ctx context.Context, req *dto.SenderDTO) (*dto.SenderDTO, error) {
	if err := helper.ValidateData(req); err != nil {
		return nil, errors.Join(common.ErrValidatingRequest, err)
	}

	_uuid, err := common.ParseUUIDMap(map[string]string{
		"company": req.CompanyID,
		"user":    req.UserID,
	})

	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	// Extract the domain from the email
	emailParts := strings.Split(req.Email, "@")
	if len(emailParts) != 2 {
		return nil, errors.New("invalid email format")
	}
	domainName := emailParts[1]

	senderExists, err := s.store.CheckSenderExists(ctx, db.CheckSenderExistsParams{
		Email:     req.Email,
		Name:      req.Name,
		CompanyID: _uuid["company"],
	})

	if err != nil {
		return nil, common.ErrFetchingRecord
	}

	if senderExists {
		return nil, common.ErrRecordExists
	}

	senderData := db.CreateSenderParams{
		UserID:    _uuid["user"],
		Email:     req.Email,
		Name:      req.Name,
		CompanyID: _uuid["company"],
	}

	_, err = s.store.FindDomainByNameAndCompany(ctx, db.FindDomainByNameAndCompanyParams{
		Domain:    domainName,
		CompanyID: _uuid["company"],
	})

	if err != nil {
		if err == sql.ErrNoRows {
			if s.HasMXRecord(domainName) {
				senderData.IsSigned = sql.NullBool{Bool: false, Valid: true}
				senderData.Verified = sql.NullBool{Bool: false, Valid: true}

				//trigger a email sending event here.
				s.sendVerificationMail(ctx, _uuid["user"], req.Email)

			} else {
				senderData.IsSigned = sql.NullBool{Bool: false, Valid: true}
				senderData.Verified = sql.NullBool{Bool: false, Valid: true}
			}
			// Don't return error here - sql.ErrNoRows is expected when domain doesn't exist
		} else {
			// Only return error for actual database issues, not for "no rows found"
			return nil, errors.Join(common.ErrFetchingRecord, err)
		}
	} else {
		// Domain exists
		senderData.IsSigned = sql.NullBool{Bool: true, Valid: true}
		senderData.Verified = sql.NullBool{Bool: true, Valid: true}
	}

	_, err = s.store.CreateSender(ctx, senderData)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (s *SenderService) HasMXRecord(domain string) bool {
	mxRecords, err := net.LookupMX(domain)
	return err == nil && len(mxRecords) > 0
}

func (s *SenderService) sendVerificationMail(ctx context.Context, userId uuid.UUID, email string) error {
	user, err := s.store.GetUserByID(ctx, userId)
	if err != nil {
		return common.ErrUserNotFound
	}

	otp := helper.GenerateOTP(20)
	_, err = s.store.CreateOTP(ctx, db.CreateOTPParams{
		UserID:    userId,
		Token:     otp,
		ExpiresAt: sql.NullTime{Time: time.Now().Add(1 * time.Hour), Valid: true},
	})
	if err != nil {
		return err
	}

	go mailer.NewEmailService().VerifySenderMail(user.Fullname, user.Email, email, otp, user.ID.String())

	return nil
}

func (s *SenderService) GetAllSenders(ctx context.Context, req *dto.FetchSenderDTO) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"company": req.CompanyID,
	})

	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	senders_count, err := s.store.CountSenderByCompanyID(ctx, _uuid["company"])
	if err != nil {
		return nil, common.ErrFetchingRecord
	}

	senders, err := s.store.ListSendersByCompanyId(ctx, db.ListSendersByCompanyIdParams{
		CompanyID: _uuid["company"],
		Limit:     int32(req.Limit),
		Offset:    int32(req.Offset),
	})

	if err != nil {
		return nil, common.ErrFetchingRecord
	}

	data := mapper.MapSendersToDTO(senders)

	items := make([]any, len(data))
	for i, v := range data {
		items[i] = v
	}

	response := common.Paginate(int(senders_count), items, int(req.Limit), int(req.Offset))
	return response, nil
}

func (s *SenderService) DeleteSender(ctx context.Context, req dto.FetchSenderDTO) error {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"sender":  req.SenderId,
		"company": req.CompanyID,
	})

	if err != nil {
		return common.ErrInvalidUUID
	}

	err = s.store.SoftDeleteSender(ctx, db.SoftDeleteSenderParams{
		ID:        _uuid["sender"],
		CompanyID: _uuid["company"],
	})

	if err != nil {
		return common.ErrDeletingRecord
	}

	return nil
}

func (s *SenderService) UpdateSender(ctx context.Context, req *dto.SenderDTO) error {
	// Validate the input data
	if err := helper.ValidateData(req); err != nil {
		return errors.Join(common.ErrValidatingRequest, err)
	}

	_uuid, err := common.ParseUUIDMap(map[string]string{
		"company": req.CompanyID,
		"id":      req.SenderId,
		"user":    req.UserID,
	})
	if err != nil {
		return common.ErrInvalidUUID
	}

	// Extract the domain from the email
	emailParts := strings.Split(req.Email, "@")
	if len(emailParts) != 2 {
		return errors.New("invalid email format")
	}
	domainName := emailParts[1]

	// Retrieve the existing sender to make sure it exists
	existingSender, err := s.store.GetSenderById(ctx, db.GetSenderByIdParams{
		CompanyID: _uuid["company"],
		ID:        _uuid["id"],
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return common.ErrRecordNotFound
		}
		return common.ErrFetchingRecord
	}

	// Check if email is being changed and if new email already exists
	if existingSender.Email != req.Email {
		senderExists, err := s.store.CheckSenderExists(ctx, db.CheckSenderExistsParams{
			Email:     req.Email,
			Name:      req.Name,
			CompanyID: _uuid["company"],
		})
		if err != nil {
			return common.ErrFetchingRecord
		}
		if senderExists {
			return common.ErrRecordExists
		}
	}

	// Prepare update data
	updateData := db.UpdateSenderParams{
		ID:        _uuid["id"],
		CompanyID: _uuid["company"],
		UserID:    _uuid["user"],
		Name:      req.Name,
		Email:     req.Email,
	}

	// DNS and domain verification logic
	_, err = s.store.FindDomainByNameAndCompany(ctx, db.FindDomainByNameAndCompanyParams{
		Domain:    domainName,
		CompanyID: _uuid["company"],
	})

	if err != nil {
		if err == sql.ErrNoRows {
			// Domain doesn't exist in our system
			if s.HasMXRecord(domainName) {
				updateData.IsSigned = sql.NullBool{Bool: false, Valid: true}
				updateData.Verified = sql.NullBool{Bool: false, Valid: true}

				// If email changed, send verification email
				if existingSender.Email != req.Email {
					// You might want to get the user ID from the existing sender or from the context
					// For now, I'll assume you have a way to get it
					userID := existingSender.UserID // Assuming UserID is available in the existing sender
					go s.sendVerificationMail(ctx, userID, req.Email)
				}
			} else {
				updateData.IsSigned = sql.NullBool{Bool: false, Valid: true}
				updateData.Verified = sql.NullBool{Bool: false, Valid: true}
			}
		} else {
			// Some other database error
			return common.ErrFetchingRecord
		}
	} else {
		// Domain exists in our system
		updateData.IsSigned = sql.NullBool{Bool: true, Valid: true}
		updateData.Verified = sql.NullBool{Bool: true, Valid: true}
	}

	// Update the sender in the database
	_, err = s.store.UpdateSender(ctx, updateData)
	if err != nil {
		return common.ErrUpdatingRecord
	}

	return nil
}
func (s *SenderService) VerifySender(ctx context.Context, req *dto.VerifySenderDTO) error {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"company": req.CompanyID,
	})

	if err != nil {
		return common.ErrInvalidUUID
	}
	if err := helper.ValidateData(req); err != nil {
		return errors.Join(common.ErrValidatingRequest, err)
	}

	err = s.store.ExecTx(ctx, func(q *db.Queries) error {
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

		err = q.UpdateSenderVerified(ctx, db.UpdateSenderVerifiedParams{
			Verified:  sql.NullBool{Bool: true, Valid: true},
			CompanyID: _uuid["company"],
			Email:     req.Email,
		})

		if err != nil {
			return fmt.Errorf("%w: %v", common.ErrUpdatingRecord, err)
		}

		if err = q.DeleteOTPById(ctx, getOtp.ID); err != nil {
			return fmt.Errorf("%w: %v", common.ErrDeletingOTP, err)
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
