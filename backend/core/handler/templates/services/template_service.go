package services

import (
	"context"
	"database/sql"
	"email-marketing-service/core/handler/templates/dto"
	"email-marketing-service/core/handler/templates/mapper"
	"email-marketing-service/internal/common"
	"email-marketing-service/internal/config"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/domain"
	"email-marketing-service/internal/factory/smtpFactory"
	"email-marketing-service/internal/helper"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"strings"
	"sync"
	"sync/atomic"
)

type Service struct {
	store db.Store
}

var (
	cfg = config.LoadEnv()
)

func NewTemplateService(store db.Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) CreateTemplate(ctx context.Context, req *dto.TemplateDTO) (any, error) {
	if err := helper.ValidateData(req); err != nil {
		return nil, errors.Join(common.ErrValidatingRequest, err)
	}

	_uuid, err := common.ParseUUIDMap(map[string]string{
		"user":    req.UserId,
		"company": req.CompanyID,
	})

	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	templateExists, err := s.store.CheckTemplateNameExists(ctx, db.CheckTemplateNameExistsParams{
		UserID:       uuid.NullUUID{UUID: _uuid["user"], Valid: true},
		TemplateName: req.TemplateName,
	})

	if err != nil {
		return nil, err
	}

	if templateExists {
		return nil, common.ErrRecordExists
	}

	// Handle EmailDesign JSON properly
	var emailDesign pqtype.NullRawMessage
	if req.EmailDesign != nil && len(req.EmailDesign) > 0 {
		// Check if it's valid JSON and not just empty brackets
		var temp interface{}
		if err := json.Unmarshal(req.EmailDesign, &temp); err != nil {
			// Invalid JSON, set to null
			emailDesign = pqtype.NullRawMessage{Valid: false}
		} else {
			// Check if it's meaningful content (not just empty array/object)
			jsonStr := string(req.EmailDesign)
			if jsonStr == "[]" || jsonStr == "{}" || strings.TrimSpace(jsonStr) == "" {
				emailDesign = pqtype.NullRawMessage{Valid: false}
			} else {
				emailDesign = pqtype.NullRawMessage{RawMessage: req.EmailDesign, Valid: true}
			}
		}
	} else {
		emailDesign = pqtype.NullRawMessage{Valid: false}
	}

	template, err := s.store.CreateTemplate(ctx, db.CreateTemplateParams{
		UserID:       uuid.NullUUID{UUID: _uuid["user"], Valid: true},
		CompanyID:    uuid.NullUUID{UUID: _uuid["company"], Valid: true},
		TemplateName: req.TemplateName,
		SenderName:   sql.NullString{String: req.SenderName, Valid: true},
		FromEmail:    sql.NullString{String: req.FromEmail, Valid: true},
		Subject:      sql.NullString{String: req.Subject, Valid: true},
		Type:         req.Type,
		EmailHtml:    sql.NullString{String: req.EmailHtml, Valid: true},
		EmailDesign:  emailDesign, // Use the properly handled emailDesign
		IsEditable:   sql.NullBool{Bool: req.IsEditable, Valid: true},
		IsPublished:  sql.NullBool{Bool: req.IsPublished, Valid: true},
		IsPublicTemplate: sql.NullBool{
			Bool:  req.IsPublicTemplate,
			Valid: true,
		},
		IsGalleryTemplate: sql.NullBool{
			Bool:  req.IsGalleryTemplate,
			Valid: true,
		},
		Tags:        sql.NullString{String: req.Tags, Valid: true},
		Description: sql.NullString{String: req.Description, Valid: true},
		ImageUrl:    sql.NullString{String: req.ImageUrl, Valid: true},
		IsActive:    sql.NullBool{Bool: req.IsActive, Valid: true},
		EditorType:  sql.NullString{String: req.EditorType, Valid: true},
	})

	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		return nil, common.ErrCreatingRecord
	}
	data := mapper.MapTemplateToDTO(template)
	return data, nil
}

func (s *Service) GetAllTemplateByType(ctx context.Context, req dto.FetchTemplateDTO) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"user": req.UserId,
	})
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	templates, err := s.store.ListTemplatesByType(ctx, db.ListTemplatesByTypeParams{
		Type:    req.Type,
		UserID:  uuid.NullUUID{UUID: _uuid["user"], Valid: true},
		Limit:   int32(req.Limit),
		Offset:  int32(req.Offset),
		Column5: req.SearchQuery,
	})
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		return nil, common.ErrFetchingRecord
	}

	count_templates, err := s.store.CountTemplatesByUserID(ctx, uuid.NullUUID{UUID: _uuid["user"], Valid: true})
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		return nil, common.ErrFetchingRecord
	}

	response := mapper.MapTemplateResponse(templates)

	items := make([]any, len(response))
	for i, v := range response {
		items[i] = v
	}

	data := common.Paginate(int(count_templates), items, req.Offset, req.Limit)
	return data, nil
}

func (s *Service) GetTemplateByID(ctx context.Context, req dto.FetchTemplateDTO) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"user":     req.UserId,
		"template": req.TemplateId,
	})
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	template, err := s.store.GetTemplateByID(ctx, db.GetTemplateByIDParams{
		UserID: uuid.NullUUID{UUID: _uuid["user"], Valid: true},
		TemplateID:     _uuid["template"],
		Type:   req.Type,
	})
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		return nil, common.ErrFetchingRecord
	}

	response := mapper.MapSingleTemplateResponse(template)
	return response, nil
}

func (s *Service) UpdateTemplate(ctx context.Context, req *dto.TemplateDTO) (*dto.TemplateDTO, error) {
	if err := helper.ValidateData(req); err != nil {
		return nil, errors.Join(common.ErrValidatingRequest, err)
	}
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"user":     req.UserId,
		"template": req.TemplateID,
	})

	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	_, err = s.store.UpdateTemplate(ctx, db.UpdateTemplateParams{
		UserID:       uuid.NullUUID{UUID: _uuid["user"], Valid: true},
		ID:           _uuid["template"],
		TemplateName: req.TemplateName,
		SenderName:   sql.NullString{String: req.SenderName, Valid: true},
		FromEmail:    sql.NullString{String: req.FromEmail, Valid: true},
		Subject:      sql.NullString{String: req.Subject, Valid: true},
		Type:         req.Type,
		EmailHtml:    sql.NullString{String: req.EmailHtml, Valid: true},
		EmailDesign:  pqtype.NullRawMessage{RawMessage: req.EmailDesign, Valid: true},
		IsEditable:   sql.NullBool{Bool: req.IsEditable, Valid: true},
		IsPublished:  sql.NullBool{Bool: req.IsPublished, Valid: true},
		IsPublicTemplate: sql.NullBool{
			Bool:  req.IsPublicTemplate,
			Valid: true,
		},
		IsGalleryTemplate: sql.NullBool{
			Bool:  req.IsGalleryTemplate,
			Valid: true,
		},
		Tags:        sql.NullString{String: req.Tags, Valid: true},
		Description: sql.NullString{String: req.Description, Valid: true},
		ImageUrl:    sql.NullString{String: req.ImageUrl, Valid: true},
		IsActive:    sql.NullBool{Bool: req.IsActive, Valid: true},
		EditorType:  sql.NullString{String: req.EditorType, Valid: true},
	})
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		return nil, common.ErrUpdatingRecord
	}
	return req, nil
}

func (s *Service) DeleteTemplate(ctx context.Context, req dto.FetchTemplateDTO) error {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"user":     req.UserId,
		"template": req.TemplateId,
	})
	if err != nil {
		return common.ErrInvalidUUID
	}

	err = s.store.SoftDeleteTemplate(ctx, db.SoftDeleteTemplateParams{
		ID:     _uuid["template"],
		UserID: uuid.NullUUID{UUID: _uuid["user"], Valid: true},
	})
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		return common.ErrDeletingRecord
	}
	return nil
}

func (s *Service) SendTestMail(ctx context.Context, d *dto.SendTestMailDTO) error {
	if err := helper.ValidateData(d); err != nil {
		return fmt.Errorf("invalid data: %w", err)
	}

	_uuid, err := common.ParseUUIDMap(map[string]string{
		"user":     d.UserId,
		"template": d.TemplateId,
	})
	if err != nil {
		return common.ErrInvalidUUID
	}

	emails := strings.Split(d.EmailAddress, ",")

	user, err := s.store.GetUserByID(context.Background(), _uuid["user"])
	if err != nil {
		return fmt.Errorf("error fetching user: %w", err)
	}

	mailUsageRecord, err := s.store.GetCurrentEmailUsage(ctx, user.CompanyID)
	if err != nil {
		return fmt.Errorf("error fetching or creating mail usage record: %w", err)
	}

	if mailUsageRecord.RemainingEmails.Int32 == 0 {
		return fmt.Errorf("you have exceeded your plan limit")
	}

	template, err := s.store.GetTemplateByID(ctx, db.GetTemplateByIDParams{
		UserID: uuid.NullUUID{UUID: _uuid["user"], Valid: true},
		TemplateID:     _uuid["template"],
		Type:   d.Type,
	})
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		return common.ErrFetchingRecord
	}

	var wg sync.WaitGroup
	var errChan = make(chan error, len(emails))
	var successCount int32 // Track successful emails

	for _, email := range emails {
		wg.Add(1)
		go func(email string) {
			defer wg.Done()

			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("Error sending batch: %v\n", err)
					errChan <- fmt.Errorf("panic: %v", err)
				}
			}()

			if err := s.proccessEmail(template.EmailHtml.String, d.Subject, email, user.Email, user.Fullname); err != nil {
				errChan <- fmt.Errorf("error processing email %s: %w", email, err)
				return
			}

			// Increment successful count using atomic operation
			atomic.AddInt32(&successCount, 1)
		}(email)
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	// Update the database with the total successful count in one operation
	if successCount > 0 {
		_, err = s.store.UpdateEmailsSentAndRemaining(ctx, db.UpdateEmailsSentAndRemainingParams{
			CompanyID:  user.CompanyID,
			EmailsSent: sql.NullInt32{Int32: successCount, Valid: true},
			ID:         mailUsageRecord.ID,
		})
		if err != nil {
			return fmt.Errorf("error updating mail usage: %w", err)
		}
	}
	return nil
}

func (s *Service) proccessEmail(design string, subject string, email string, from string, fromName string) error {
	valid := helper.IsValidEmail(email)
	if !valid {
		return nil
	}

	sender := &domain.Sender{Email: from, Name: &fromName}
	recipient := domain.Recipient{Email: email}
	request := &domain.EmailRequest{
		Sender:      *sender,
		To:          recipient,
		Subject:     subject,
		HtmlContent: &design,
	}
	println(cfg.MAIL_PROCESSOR)
	mailS, err := smtpfactory.MailFactory(cfg.MAIL_PROCESSOR)
	if err != nil {
		return fmt.Errorf("failed to create mail factory: %w", err)
	}

	if err := mailS.HandleSendMail(request); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
