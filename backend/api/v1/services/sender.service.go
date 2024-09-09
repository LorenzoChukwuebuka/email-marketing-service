package services

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/utils"
	"errors"
	"fmt"
	"strings"
)

type SenderServices struct {
	SenderRepo *repository.SenderRepository
	DomainRepo *repository.DomainRepository
}

func NewSenderServices(domainRepo *repository.DomainRepository, senderRepo *repository.SenderRepository) *SenderServices {
	return &SenderServices{
		SenderRepo: senderRepo,
		DomainRepo: domainRepo,
	}
}

func (s *SenderServices) CreateSender(d *dto.SenderDTO) error {

	if err := utils.ValidateData(d); err != nil {
		return fmt.Errorf("invalid plan data: %w", err)
	}

	// Extract the domain from the email
	emailParts := strings.Split(d.Email, "@")
	if len(emailParts) != 2 {
		return errors.New("invalid email format")
	}
	domainName := emailParts[1]

	senderModel := &model.Sender{
		UserID: d.UserID,
		Email:  d.Email,
		Name:   d.Name,
	}

	senderExists, err := s.SenderRepo.CheckIfSenderExists(senderModel)
	if err != nil {
		return err
	}

	if senderExists {
		return fmt.Errorf("the sender details alrady exists")
	}

	getDomain, err := s.DomainRepo.FindDomain(d.UserID, domainName)

	if err != nil {
		if err.Error() == "domain not found" {
			senderModel.IsSigned = false
		} else {
			return err
		}
	}


	if getDomain != nil {
		if getDomain.Verified {
			senderModel.IsSigned = true
		} else {
			senderModel.IsSigned = false
		}
	
	}

	 if err := s.SenderRepo.CreateSender(senderModel); err != nil {
		return err
	 }

	return nil

}
