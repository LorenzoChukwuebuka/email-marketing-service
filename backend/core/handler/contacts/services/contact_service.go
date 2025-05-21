package services

import (
	"context"
	"database/sql"
	"email-marketing-service/core/handler/contacts/dto"
	"email-marketing-service/core/handler/contacts/mapper"
	"email-marketing-service/internal/common"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/helper"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	store db.Store
}

func NewContactService(store db.Store) *Service {
	return &Service{store: store}
}

// CreateContact creates a new contact based on the provided DTO.
func (s *Service) CreateContact(ctx context.Context, d *dto.ContactRequestDTO) (any, error) {
	if err := helper.ValidateData(d); err != nil {
		return nil, errors.Join(common.ErrValidatingRequest, err)
	}

	userID, err := uuid.Parse(d.UserId)
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	companyId, err := uuid.Parse(d.CompanyID)
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	exists, err := s.store.CheckIfContactEmailExists(ctx, db.CheckIfContactEmailExistsParams{
		Email:  d.Email,
		UserID: userID,
	})

	if err != nil {
		return nil, fmt.Errorf("error checking contact existence: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("contact already exists")
	}

	if err := s.store.CreateContact(ctx, db.CreateContactParams{
		CompanyID:    companyId,
		UserID:       userID,
		FirstName:    d.FirstName,
		LastName:     d.LastName,
		Email:        d.Email,
		FromOrigin:   s.getContactSource(d.From),
		IsSubscribed: sql.NullBool{Bool: true, Valid: true},
		CreatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
	}); err != nil {
		return nil, common.ErrCreatingRecord
	}

	return map[string]any{
		"data":    d,
		"message": "contact added successfully",
	}, nil
}

// getContactSource determines the source of contact with a default value.
func (s *Service) getContactSource(source string) string {
	if source == "" {
		return "web"
	}
	return source
}

// UploadContactViaCSV reads a CSV file and uploads contacts in bulk.
func (s *Service) UploadContactViaCSV(ctx context.Context, file multipart.File, filename, userId string, companyId string) error {
	reader := csv.NewReader(file)
	columnMap, err := s.parseCSVHeader(reader)
	if err != nil {
		return err
	}

	newContacts, err := s.processCSVRecords(ctx, reader, columnMap, userId, companyId)
	if err != nil {
		return err
	}

	if len(newContacts) > 0 {
		userID, err := uuid.Parse(userId)
		if err != nil {
			return common.ErrInvalidUUID
		}

		companyUUID, err := uuid.Parse(companyId)
		if err != nil {
			return common.ErrInvalidUUID
		}

		for _, contact := range newContacts {
			if err := s.store.CreateContact(ctx, db.CreateContactParams{
				CompanyID:    companyUUID,
				UserID:       userID,
				FirstName:    contact.FirstName,
				LastName:     contact.LastName,
				Email:        contact.Email,
				FromOrigin:   s.getContactSource(contact.From),
				IsSubscribed: sql.NullBool{Bool: contact.IsSubscribed, Valid: true},
				CreatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
			}); err != nil {
				return fmt.Errorf("error creating contact: %w", err)
			}
		}
	}

	return nil
}

// parseCSVHeader reads and validates CSV header, returning a column map.
func (s *Service) parseCSVHeader(reader *csv.Reader) (map[string]int, error) {
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV header: %w", err)
	}

	columnMap := make(map[string]int)
	for i, column := range header {
		columnMap[strings.ToLower(strings.TrimSpace(column))] = i
	}

	requiredColumns := []string{"first name", "last name", "email"}
	for _, col := range requiredColumns {
		if _, exists := columnMap[col]; !exists {
			return nil, fmt.Errorf("required column '%s' is missing from the CSV", col)
		}
	}

	return columnMap, nil
}

// processCSVRecords reads each record and creates a new contact list for bulk insert.
func (s *Service) processCSVRecords(ctx context.Context, reader *csv.Reader, columnMap map[string]int, userId string, companyId string) ([]dto.ContactRequestDTO, error) {
	var newContacts []dto.ContactRequestDTO

	userID, err := uuid.Parse(userId)
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV record: %w", err)
		}

		contact := s.createContactFromRecord(record, columnMap, userId, companyId)

		exists, err := s.store.CheckIfContactEmailExists(ctx, db.CheckIfContactEmailExistsParams{
			Email:  contact.Email,
			UserID: userID,
		})
		if err != nil {
			return nil, fmt.Errorf("error checking email existence: %w", err)
		}

		if !exists {
			newContacts = append(newContacts, contact)
		}
	}

	return newContacts, nil
}

// createContactFromRecord generates a contact from a CSV record.
func (s *Service) createContactFromRecord(record []string, columnMap map[string]int, userId string, companyId string) dto.ContactRequestDTO {
	contact := dto.ContactRequestDTO{
		FirstName:    record[columnMap["first name"]],
		LastName:     record[columnMap["last name"]],
		Email:        record[columnMap["email"]],
		From:         "web",
		UserId:       userId,
		CompanyID:    companyId,
		IsSubscribed: true,
	}

	if idx, exists := columnMap["from"]; exists && idx < len(record) {
		contact.From = record[idx]
	}

	return contact
}

func (s *Service) UpdateContact(ctx context.Context, req *dto.EditContactDTO) error {
	contactUUID, err := uuid.Parse(req.ContactId)
	if err != nil {
		return common.ErrInvalidUUID
	}

	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		return common.ErrInvalidUUID
	}

	err = s.store.UpdateContact(ctx, db.UpdateContactParams{
		ID:           contactUUID,
		UserID:       userUUID,
		IsSubscribed: sql.NullBool{Bool: req.IsSubscribed, Valid: true},
		FirstName:    sql.NullString{String: req.FirstName, Valid: true},
		LastName:     sql.NullString{String: req.LastName, Valid: true},
		Email:        sql.NullString{String: req.Email, Valid: true},
		FromOrigin:   sql.NullString{String: req.From, Valid: true},
	})

	if err != nil {
		return common.ErrUpdatingRecord
	}

	return nil
}

func (s *Service) DeleteContact(ctx context.Context, userId string, contactId string) error {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return common.ErrInvalidUUID
	}

	contactUUID, err := uuid.Parse(contactId)
	if err != nil {
		return common.ErrInvalidUUID
	}

	err = s.store.DeleteContact(ctx, db.DeleteContactParams{
		ID:     contactUUID,
		UserID: userUUID,
	})
	if err != nil {
		return common.ErrDeletingRecord
	}
	return nil
}

func (s *Service) GetAllContacts(ctx context.Context, req dto.FetchContactDTO) (any, error) {
	userUUID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	companyUUID, err := uuid.Parse(req.CompanyID)
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	contacts, err := s.store.GetAllContacts(ctx, db.GetAllContactsParams{
		UserID:    userUUID,
		CompanyID: companyUUID,
		Limit:     int32(req.Limit),
		Offset:    int32(req.Offset),
	})

	if err != nil {
		return nil, common.ErrFetchingRecord
	}

	contactsResponse := make([]any, 0)

	for _, contact := range contacts {
		value := mapper.MapContactAllContactResponse(contact)
		contactsResponse = append(contactsResponse, *value)
	}

	data := common.Paginate(10, contactsResponse, req.Offset, req.Limit)

	return data, nil

}
