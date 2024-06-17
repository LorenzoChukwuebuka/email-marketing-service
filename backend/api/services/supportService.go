package services

import (
	"email-marketing-service/api/dto"
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
	"fmt"
	"os"
	"path/filepath"
)

type SupportTicketService struct {
	SupportRepository *repository.SupportRepository
	UserRepo          *repository.UserRepository
}

func NewSupportTicketService(supportRepo *repository.SupportRepository, userRepo *repository.UserRepository) *SupportTicketService {
	return &SupportTicketService{
		SupportRepository: supportRepo,
		UserRepo:          userRepo,
	}
}

func (s *SupportTicketService) CreateSupportTicket(d *dto.SupportTicket) error {
	// Log the ticket data (for debugging purposes)
	fmt.Println(d)

	// Define the folder path
	uploadFolder := "uploads/tickets"
	if err := os.MkdirAll(uploadFolder, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Save the uploaded file and set the file path
	if d.TicketFile != nil {
		fileName := d.UserID + "_" + d.Subject
		filePath := filepath.Join(uploadFolder, fileName)
		if err := os.WriteFile(filePath, *d.TicketFile, 0644); err != nil {
			return fmt.Errorf("failed to save file: %v", err)
		}
		d.FilePath = filePath
	}

	userUUID := &model.User{UUID: d.UserID}

	userId, err := s.UserRepo.FindUserById(userUUID)

	if err != nil {
		return err
	}

	supportModel := &model.SupportTicket{
		Subject:     d.Subject,
		Description: d.Description,
		SenderID:    uint(userId.ID),
		Priority:    model.Priority(d.Priority),
		Status:      model.Status(d.Status),
		AssignedTo:  0,
	}

	id, err := s.SupportRepository.CreateSupportTicket(supportModel)

	if err != nil {
		return err
	}

	supportTicketMessageModel := &model.TicketFiles{
		TicketID:   id,
		TicketFile: d.FilePath,
	}

	err = s.SupportRepository.CreateSupportTicketFile(supportTicketMessageModel)
	if err != nil {
		return err
	}

	return nil
}
