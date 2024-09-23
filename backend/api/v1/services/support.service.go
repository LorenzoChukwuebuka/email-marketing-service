package services

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	adminrepository "email-marketing-service/api/v1/repository/admin"
	"email-marketing-service/api/v1/utils"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

type SupportTicketService struct {
	SupportRepository     *repository.SupportRepository
	UserRepository        *repository.UserRepository
	UserNotification      *repository.UserNotificationRepository
	AdminNotificationRepo *adminrepository.AdminNotificationRepository
}

func NewSupportTicketService(supportRepo *repository.SupportRepository,
	userRepo *repository.UserRepository,
	userNotificationRepo *repository.UserNotificationRepository,adminNotificationRepo *adminrepository.AdminNotificationRepository) *SupportTicketService {
	return &SupportTicketService{
		SupportRepository: supportRepo,
		UserRepository:    userRepo,
		UserNotification:  userNotificationRepo,
		AdminNotificationRepo: adminNotificationRepo,
	}
}

func (s *SupportTicketService) CreateSupportTicket(userID string, req *dto.CreateSupportTicketRequest) (*dto.CreateSupportTicketResponse, error) {

	tx := s.SupportRepository.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	userModel := &model.User{UUID: userID}

	user, err := s.UserRepository.FindUserById(userModel)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	ticket := &model.SupportTicket{
		UserID:       userID,
		Name:         user.FullName,
		Email:        user.Email,
		Subject:      req.Subject,
		TicketNumber: utils.GenerateUniqueRandomNumbers(8),
		Description:  req.Description,
		Status:       model.PendingTicket,
		Priority:     model.Priority(req.Priority),
	}

	if req.File != nil {
		filePath, err := s.saveFile(req.File, userID)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to save file: %w", err)
		}
		ticket.Files = []model.TicketFile{{
			FileName: req.File.Filename,
			FilePath: filePath,
		}}
	}

	if req.Message != "" {
		ticket.Messages = []model.TicketMessage{{
			UserID:  userID,
			Message: req.Message,
			IsAdmin: false,
		}}
	}

	createdTicket, err := s.SupportRepository.CreateSupportTicket(ticket)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create support ticket: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	notificationTitle := fmt.Sprintf("You have opened a ticket with subject %s. Our agent will respond to you soon", createdTicket.Name)
	if err := utils.CreateNotification(s.UserNotification, userID, notificationTitle); err != nil {
		fmt.Printf("Failed to create notification: %v\n", err)
	}

	return &dto.CreateSupportTicketResponse{
		TicketID: createdTicket.ID,
		Message:  "Your ticket has been successfully created. An agent will reply to you soon.",
	}, nil
}

func (s *SupportTicketService) ReplyToTicket(ticketID string, userID string, req *dto.ReplyTicketRequest) (*dto.ReplyTicketResponse, error) {
	ticket, err := s.SupportRepository.FindTicketByID(ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to find ticket: %w", err)
	}

	message := &model.TicketMessage{
		TicketID: ticket.ID,
		UserID:   userID,
		Message:  req.Message,
		IsAdmin:  false, // Assuming this is a user reply, not an admin
	}

	_, err = s.SupportRepository.CreateTicketMessage(message)
	if err != nil {
		return nil, fmt.Errorf("failed to create ticket message: %w", err)
	}

	if req.File != nil {
		filePath, err := s.saveFile(req.File, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to save file: %w", err)
		}

		fileModel := &model.TicketFile{
			TicketID: ticket.ID,
			FileName: req.File.Filename,
			FilePath: filePath,
		}

		_, err = s.SupportRepository.CreateTicketFile(fileModel)
		if err != nil {
			return nil, fmt.Errorf("failed to create ticket file record: %w", err)
		}
	}

	now := time.Now()
	_, err = s.SupportRepository.UpdateTicket(&model.SupportTicket{
		Model:     gorm.Model{ID: ticket.ID},
		LastReply: &now,
		Status:    model.OpenTicket,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update ticket: %w", err)
	}

	notificationTitle := fmt.Sprintf("You have replied a ticket with subject %s. Our agent will respond to you soon", ticket.Name)
	if err := utils.CreateNotification(s.UserNotification, userID, notificationTitle); err != nil {
		fmt.Printf("Failed to create notification: %v\n", err)
	}

	return &dto.ReplyTicketResponse{

		Message: "Your reply has been successfully added to the ticket.",
	}, nil
}

func (s *SupportTicketService) saveFile(file *multipart.FileHeader, userID string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	uploadFolder := filepath.Join("uploads", "tickets", userID)
	if err := os.MkdirAll(uploadFolder, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	filePath := filepath.Join(uploadFolder, fileName)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	return filePath, nil
}

func (s *SupportTicketService) GetTicketWithDetails(ticketID string) (*model.SupportTicketResponse, error) {
	ticket, err := s.SupportRepository.GetTicketWithDetails(ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	// Populate user information for non-admin messages
	for i, message := range ticket.Messages {
		if !message.IsAdmin {
			user, err := s.UserRepository.FindUserById(&model.User{UUID: message.UserID})
			if err != nil {
				return nil, fmt.Errorf("failed to get user for message: %w", err)
			}
			ticket.Messages[i].User = &model.UserResponse{
				UUID:        user.UUID,
				FullName:    user.FullName,
				Email:       user.Email,
				Company:     user.Company,
				PhoneNumber: user.PhoneNumber,
				Verified:    user.Verified,
				Blocked:     user.Blocked,
				VerifiedAt:  user.VerifiedAt,
				CreatedAt:   user.CreatedAt,
			}
		}

		if message.IsAdmin {
			fmt.Println("admin for now")
		}
	}

	return ticket, nil
}

func (s *SupportTicketService) GetTicketsByUserID(userID string) ([]model.SupportTicketResponse, error) {
	tickets, err := s.SupportRepository.GetTicketsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user tickets: %w", err)
	}
	return tickets, nil
}


