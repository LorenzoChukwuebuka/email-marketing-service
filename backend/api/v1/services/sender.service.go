package services

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/utils"
	"errors"
	"fmt"
	"net"
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
		return fmt.Errorf("the sender details already exists")
	}

	getDomain, err := s.DomainRepo.FindDomain(d.UserID, domainName)

	if err != nil {
		if err.Error() == "domain not found" {

			if s.HasMXRecord(domainName) {
				senderModel.IsSigned = false
				senderModel.Verified = true
			} else {
				senderModel.IsSigned = false
				senderModel.Verified = false
			}

		} else {
			return err
		}
	}

	if getDomain != nil {
		if getDomain.Verified {
			senderModel.IsSigned = true
			senderModel.Verified = true
		} else {
			senderModel.IsSigned = false
			senderModel.Verified = false
		}

	}

	if err := s.SenderRepo.CreateSender(senderModel); err != nil {
		return err
	}

	return nil

}

func (s *SenderServices) HasMXRecord(domain string) bool {
	mxRecords, err := net.LookupMX(domain)
	return err == nil && len(mxRecords) > 0
}

func (s *SenderServices) GetAllSenders(userId string, page int, pageSize int, searchQuery string) (repository.PaginatedResult, error) {
	paginationParams := repository.PaginationParams{Page: page, PageSize: pageSize}
	getSenders, err := s.SenderRepo.GetAllSenders(userId, searchQuery, paginationParams)

	if err != nil {
		return repository.PaginatedResult{}, err
	}

	return getSenders, nil
}

func (s *SenderServices) DeleteSender() {}

func (s *SenderServices) UpdateSender() {}
