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
	"strings"
)

type ContactService struct {
	ContactRepo *repository.ContactRepository
}

func NewContactService(contactRepo *repository.ContactRepository) *ContactService {
	return &ContactService{ContactRepo: contactRepo}
}

// CreateContact creates a new contact based on the provided DTO.
func (s *ContactService) CreateContact(d *dto.ContactDTO) (map[string]interface{}, error) {
	if err := utils.ValidateData(d); err != nil {
		return nil, fmt.Errorf("invalid data: %w", err)
	}

	contactModel := s.createContactModel(d)

	if exists, err := s.ContactRepo.CheckIfEmailExists(contactModel); err != nil {
		return nil, fmt.Errorf("error checking contact existence: %w", err)
	} else if exists {
		return nil, fmt.Errorf("contact already exists")
	}

	if err := s.ContactRepo.CreateContact(contactModel); err != nil {
		return nil, fmt.Errorf("error creating contact: %w", err)
	}

	return map[string]interface{}{
		"data":    contactModel,
		"message": "contact added successfully",
	}, nil
}

// createContactModel generates a Contact model from DTO.
func (s *ContactService) createContactModel(d *dto.ContactDTO) *model.Contact {
	return &model.Contact{
		UUID:         uuid.New().String(),
		FirstName:    d.FirstName,
		LastName:     d.LastName,
		Email:        d.Email,
		From:         s.getContactSource(d.From),
		UserId:       d.UserId,
		IsSubscribed: d.IsSubscribed,
	}
}

// getContactSource determines the source of contact with a default value.
func (s *ContactService) getContactSource(source string) string {
	if source == "" {
		return "web"
	}
	return source
}

// UploadContactViaCSV reads a CSV file and uploads contacts in bulk.
func (s *ContactService) UploadContactViaCSV(file multipart.File, filename, userId string) error {
	reader := csv.NewReader(file)

	columnMap, err := s.parseCSVHeader(reader)
	if err != nil {
		return err
	}

	newContacts, err := s.processCSVRecords(reader, columnMap, userId)
	if err != nil {
		return err
	}

	if len(newContacts) > 0 {
		if err := s.ContactRepo.BulkCreateContacts(newContacts); err != nil {
			return fmt.Errorf("error bulk inserting contacts: %w", err)
		}
	}

	return nil
}

// parseCSVHeader reads and validates CSV header, returning a column map.
func (s *ContactService) parseCSVHeader(reader *csv.Reader) (map[string]int, error) {
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
func (s *ContactService) processCSVRecords(reader *csv.Reader, columnMap map[string]int, userId string) ([]model.Contact, error) {
	var newContacts []model.Contact

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV record: %w", err)
		}

		contact := s.createContactFromRecord(record, columnMap, userId)

		exists, err := s.ContactRepo.CheckIfEmailExists(&contact)
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
func (s *ContactService) createContactFromRecord(record []string, columnMap map[string]int, userId string) model.Contact {
	contact := model.Contact{
		UUID:         uuid.New().String(),
		FirstName:    record[columnMap["first name"]],
		LastName:     record[columnMap["last name"]],
		Email:        record[columnMap["email"]],
		From:         "web",
		UserId:       userId,
		IsSubscribed: true,
	}

	if idx, exists := columnMap["from"]; exists && idx < len(record) {
		contact.From = record[idx]
	}

	return contact
}

func (s *ContactService) UpdateContact(d *dto.EditContactDTO) error {
	contactModel := &model.Contact{
		UUID:      d.ContactId,
		UserId:    d.UserId,
		FirstName: d.FirstName,
		LastName:  d.LastName,
		Email:     d.LastName,
		From:      d.From,
	}

	if err := s.ContactRepo.UpdateContact(contactModel); err != nil {
		return err
	}

	return nil
}

func (s *ContactService) GetAllContacts(userId string, page int, pageSize int, searchQuery string) (repository.PaginatedResult, error) {
	paginationParams := repository.PaginationParams{Page: page, PageSize: pageSize}
	contacts, err := s.ContactRepo.GetAllContacts(userId, paginationParams, searchQuery)
	if err != nil {
		return repository.PaginatedResult{}, err
	}

	// If you want to check for empty results, you can do:
	if contacts.TotalCount == 0 {
		return repository.PaginatedResult{}, nil
	}

	return contacts, nil
}

func (s *ContactService) DeleteContact(userId string, contactId string) error {
	if err := s.ContactRepo.DeleteContact(userId, contactId); err != nil {
		return err
	}

	return nil
}

func (s *ContactService) CreateGroup(d *dto.ContactGroupDTO) (map[string]interface{}, error) {
	if err := utils.ValidateData(d); err != nil {
		return nil, fmt.Errorf("invalid  data: %w", err)
	}

	groupModel := &model.ContactGroup{
		UUID:        uuid.New().String(),
		GroupName:   d.GroupName,
		Description: d.Description,
		UserId:      d.UserId,
	}

	groupNameExists, err := s.ContactRepo.CheckIfGroupNameExists(groupModel)
	if err != nil {
		return nil, err
	}

	if groupNameExists {
		return nil, fmt.Errorf("group name already exists")
	}

	if err := s.ContactRepo.CreateGroup(groupModel); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"data":    groupModel,
		"message": "contact group added successfully",
	}, nil

}

func (s *ContactService) AddContactsToGroup(d *dto.AddContactsToGroupDTO) (map[string]interface{}, error) {
	if err := utils.ValidateData(d); err != nil {
		return nil, fmt.Errorf("invalid  data: %w", err)
	}

	getContactId, err := s.ContactRepo.GetASingleContact(d.ContactId, d.UserId)

	if err != nil {
		return nil, err
	}

	getGroupId, err := s.ContactRepo.GetASingleGroup(d.UserId, d.GroupId)

	if err != nil {
		return nil, err
	}

	addToGroupModel := &model.UserContactGroup{
		UUID:           uuid.New().String(),
		UserId:         d.UserId,
		ContactId:      getContactId.ID,
		ContactGroupId: getGroupId.ID,
	}

	err = s.ContactRepo.AddContactsToGroup(addToGroupModel)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"data":    addToGroupModel,
		"message": "contact group added successfully",
	}, nil

}

func (s *ContactService) RemoveContactFromGroup(d *dto.AddContactsToGroupDTO) error {
	if err := utils.ValidateData(d); err != nil {
		return fmt.Errorf("invalid data: %w", err)
	}

	contactId, err := s.ContactRepo.GetASingleContact(d.ContactId, d.UserId)

	if err != nil {
		return err
	}

	groupId, err := s.ContactRepo.GetASingleGroup(d.UserId, d.GroupId)

	if err != nil {
		return err
	}

	if err := s.ContactRepo.RemoveContactFromGroup(int(groupId.ID), d.UserId, int(contactId.ID)); err != nil {
		return err
	}

	return nil
}

func (s *ContactService) UpdateContactGroup(d *dto.ContactGroupDTO, groupId string) error {

	groupModel := &model.ContactGroup{
		UUID:        groupId,
		GroupName:   d.GroupName,
		Description: d.Description,
		UserId:      d.UserId,
	}

	if err := s.ContactRepo.UpdateGroup(groupModel); err != nil {
		return err
	}
	return nil
}

func (s *ContactService) DeleteContactGroup(userId string, groupId string) error {

	res, err := s.ContactRepo.GetASingleGroup(userId, groupId)

	if err != nil {
		return err
	}

	if err := s.ContactRepo.DeleteContactGroup(userId, int(res.ID)); err != nil {
		return err
	}

	return nil
}

func (s *ContactService) GetAllContactGroups(userId string, page int, pageSize int, searchQuery string) (repository.PaginatedResult, error) {
	paginationParams := repository.PaginationParams{Page: page, PageSize: pageSize}

	result, err := s.ContactRepo.GetAllGroups(userId, paginationParams, searchQuery)

	if err != nil {
		return repository.PaginatedResult{}, err
	}

	return result, nil
}

func (s *ContactService) GetASingleGroupWithContacts(userId string, groupId string) (*model.ContactGroupResponse, error) {

	result, err := s.ContactRepo.GetASingleGroupWithContacts(userId, groupId)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *ContactService) GetContactCount(userId string) (map[string]int64, error) {
	result, err := s.ContactRepo.GetContactCount(userId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ContactService) GetContactSubscriptionStatusForDashboard(userId string) (map[string]int64, error) {
	result, err := s.ContactRepo.GetContactSubscriptionStatusForDashboard(userId)

	if err != nil {
		return nil, err
	}

	return result, nil
}
