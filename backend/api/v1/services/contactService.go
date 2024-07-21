package services

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/utils"
	"fmt"
	"github.com/google/uuid"
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
		UUID: uuid.New().String(),
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

	return nil, nil

}

func (s *ContactService) UploadContactViaCSV() {}

func (s *ContactService) UpdateContact() {}

func (s *ContactService) GetAllContacts() {}

func (s *ContactService) DeleteContact() {}

func (s *ContactService) CreateGroup() {}

func (s *ContactService) AddContactsToGroup() {}

func (s *ContactService) RemoveContactFromGroup() {}

func (s *ContactService) DeleteContactGroup() {}
