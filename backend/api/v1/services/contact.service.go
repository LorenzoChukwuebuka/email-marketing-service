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
	return &ContactService{
		ContactRepo: contactRepo,
	}

}

func (s *ContactService) CreateContact(d *dto.ContactDTO) (map[string]interface{}, error) {

	if err := utils.ValidateData(d); err != nil {
		return nil, fmt.Errorf("invalid data: %w", err)
	}

	contactModel := &model.Contact{
		UUID:      uuid.New().String(),
		FirstName: d.FirstName,
		LastName:  d.LastName,
		Email:     d.Email,
		From: func() string {
			if d.From != "" {
				return d.From
			}
			return "web"
		}(),
		UserId:       d.UserId,
		IsSubscribed: d.IsSubscribed || true,
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

func (s *ContactService) GetAllContacts(userId string, page int, pageSize int) (repository.PaginatedResult, error) {
	paginationParams := repository.PaginationParams{Page: page, PageSize: pageSize}
	contacts, err := s.ContactRepo.GetAllContacts(userId, paginationParams)
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

func (s *ContactService) UpdateContactGroup(d *dto.ContactGroupDTO) error {

	groupModel := &model.ContactGroup{

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
	if err := s.ContactRepo.DeleteContactGroup(userId, groupId); err != nil {
		return err
	}

	return nil
}

func (s *ContactService) GetAllContactGroups(userId string, page int, pageSize int) (repository.PaginatedResult, error) {
	paginationParams := repository.PaginationParams{Page: page, PageSize: pageSize}

	result, err := s.ContactRepo.GetAllGroups(userId, paginationParams)

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
