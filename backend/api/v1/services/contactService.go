package services

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/utils"
	"encoding/csv"
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	//"strconv"
	"strings"
	"time"
)

type ContactService struct {
	ContactRepo *repository.ContactRepository
}

func NewContactService(contactRepo *repository.ContactRepository) *ContactService {
	return &ContactService{
		ContactRepo: contactRepo,
	}

}

func (s *ContactService) CreateContact(d *dto.ContactDTO) (map[string]interface{}, error) {

	if err := utils.ValidateData(d); err != nil {
		return nil, fmt.Errorf("invalid plan data: %w", err)
	}

	contactModel := &model.Contact{
		UUID:      uuid.New().String(),
		FirstName: d.FirstName,
		LastName:  d.LastName,
		Email:     d.Email,
		From:      d.From,
		UserId:    d.UserId,
		CreatedAt: time.Now().UTC(),
	}

	checkIfUserExists, err := s.ContactRepo.CheckIfEmailExists(contactModel)

	if err != nil {
		return nil, err
	}

	if checkIfUserExists {
		return nil, fmt.Errorf("contact already exists")
	}

	if err := s.ContactRepo.CreateContact(contactModel); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"data":    contactModel,
		"message": "contact added successfully",
	}, nil

}

func (s *ContactService) UploadContactViaCSV(file multipart.File, filename string, userId string) error {
	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read the header
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("error reading CSV header: %w", err)
	}

	// Create a map of column indices
	columnMap := make(map[string]int)
	for i, column := range header {
		columnMap[strings.ToLower(strings.TrimSpace(column))] = i
	}

	// Validate required columns
	requiredColumns := []string{"first name", "last name", "email"}
	for _, col := range requiredColumns {
		if _, exists := columnMap[col]; !exists {
			return fmt.Errorf("required column '%s' is missing from the CSV", col)
		}
	}

	// Process the records
	var newContacts []model.Contact
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading CSV record: %w", err)
		}

		contact := model.Contact{
			UUID:      uuid.New().String(),
			FirstName: record[columnMap["first name"]],
			LastName:  record[columnMap["last name"]],
			Email:     record[columnMap["email"]],
			From:      "web",
			UserId:    userId,
			CreatedAt: time.Now(),
		}

		if idx, exists := columnMap["from"]; exists && idx < len(record) {
			contact.From = record[idx]
		}

		// Check if email already exists
		exists, err := s.ContactRepo.CheckIfEmailExists(&contact)
		if err != nil {
			return fmt.Errorf("error checking email existence: %w", err)
		}

		if !exists {
			newContacts = append(newContacts, contact)
		}
	}

	// Batch insert new contacts
	if len(newContacts) > 0 {
		err = s.ContactRepo.BulkCreateContacts(newContacts)
		if err != nil {
			return fmt.Errorf("error bulk inserting contacts: %w", err)
		}
	}

	return nil
}

func (s *ContactService) UpdateContact() {}

func (s *ContactService) GetAllContacts(userId string) ([]model.ContactResponse, error) {
	contacts, err := s.ContactRepo.GetAllContacts(userId)
	if err != nil {
		return nil, err
	}

	if len(contacts) == 0 {
		return []model.ContactResponse{}, nil
	}

	return contacts, nil
}

func (s *ContactService) DeleteContact(userId string,contactId string) {
	
}

func (s *ContactService) CreateGroup() {}

func (s *ContactService) AddContactsToGroup() {}

func (s *ContactService) RemoveContactFromGroup() {}

func (s *ContactService) DeleteContactGroup() {}
