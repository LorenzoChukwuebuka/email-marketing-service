package adminservice

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	adminrepository "email-marketing-service/api/v1/repository/admin"
	"email-marketing-service/api/v1/utils"
	"fmt"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type AdminSupportService struct {
	AdminSupportRepo      *adminrepository.AdminSupportRepository
	AdminNotificationRepo *adminrepository.AdminNotificationRepository
	UserNotificationRepo  *repository.UserNotificationRepository
	SupportRepository     *repository.SupportRepository
}

func NewAdminSupportService(adminsupportRepo *adminrepository.AdminSupportRepository,
	adminNotificationRepo *adminrepository.AdminNotificationRepository,
	usernotificationRepo *repository.UserNotificationRepository, supportRepo *repository.SupportRepository) *AdminSupportService {
	return &AdminSupportService{
		AdminSupportRepo:      adminsupportRepo,
		AdminNotificationRepo: adminNotificationRepo,
		UserNotificationRepo:  usernotificationRepo,
		SupportRepository:     supportRepo,
	}
}

// GetAllTickets fetches all support tickets with optional search filter and pagination
func (s *AdminSupportService) GetAllTickets(search string, page int, pageSize int) (repository.PaginatedResult, error) {
	params := repository.PaginationParams{Page: page, PageSize: pageSize}
	tickets, err := s.AdminSupportRepo.GetAllTickets(search, params)
	if err != nil {
		return repository.PaginatedResult{}, err
	}
	return tickets, nil
}

func (s *AdminSupportService) ReplyToTicket(ticketID string, userID string, req *dto.ReplyTicketRequest) (*dto.ReplyTicketResponse, error) {
	// Find the ticket by its ID
	ticket, err := s.SupportRepository.FindTicketByID(ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to find ticket: %w", err)
	}

	// Create a new ticket message object, assuming this is an admin reply
	message := &model.TicketMessage{
		TicketID: ticket.ID,   // Associate the message with the ticket
		UserID:   userID,      // Set the ID of the admin who is replying
		Message:  req.Message, // Set the content of the reply message
		IsAdmin:  true,        // Mark the message as from an admin
	}

	// Save the message in the database
	msg, err := s.SupportRepository.CreateTicketMessage(message)
	if err != nil {
		return nil, fmt.Errorf("failed to create ticket message: %w", err)
	}

	// Check if there are any files attached in the reply request
	if req.File != nil {
		for _, file := range req.File {
			// Save the file to the filesystem
			filePath, err := s.saveFile(file, userID)
			if err != nil {
				return nil, fmt.Errorf("failed to save file: %w", err)
			}

			// Create a new record for the saved file and associate it with the message
			fileModel := &model.TicketFile{
				MessageID: msg.ID,        // Associate the file with the message
				FileName:  file.Filename, // Store the original file name
				FilePath:  filePath,      // Store the file path where it was saved
			}

			// Save the file record in the database
			_, err = s.SupportRepository.CreateTicketFile(fileModel)
			if err != nil {
				return nil, fmt.Errorf("failed to create ticket file record: %w", err)
			}
		}
	}

	// Update the ticket's status and last reply time after the message is added
	now := time.Now()
	_, err = s.SupportRepository.UpdateTicket(&model.SupportTicket{
		Model:     gorm.Model{ID: ticket.ID}, // Update the ticket with its ID
		LastReply: &now,                      // Set the last reply time to the current time
		Status:    model.OpenTicket,          // Change the status to 'open' to indicate a reply
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update ticket: %w", err)
	}

	// Send a notification to the user indicating their ticket has been replied to
	notificationTitle := fmt.Sprintf("Your ticket with subject %s has been replied to", ticket.Subject)
	if err := utils.CreateNotification(s.UserNotificationRepo, ticket.UserID, notificationTitle); err != nil {
		fmt.Printf("Failed to create notification: %v\n", err)
	}

	// Send a notification to the admin confirming they replied to the ticket
	adminNotificationTitle := fmt.Sprintf("You have replied to a ticket with subject %s.", ticket.Subject)
	if err := utils.CreateAdminNotifications(s.AdminNotificationRepo, userID, "#", adminNotificationTitle); err != nil {
		fmt.Printf("Failed to create admin notification: %v\n", err)
	}

	// Return a success response indicating that the reply has been added to the ticket
	return &dto.ReplyTicketResponse{
		Message: "Your reply has been successfully added to the ticket.", // Success message
	}, nil
}

// CreateTicketFile adds a file to the support ticket and returns the file response
func (s *AdminSupportService) CreateTicketFile(file *model.TicketFile) (*model.TicketFileResponse, error) {
	fileResponse, err := s.SupportRepository.CreateTicketFile(file)
	if err != nil {
		return nil, err
	}
	return fileResponse, nil
}

// FindTicketsByUserID fetches all tickets for a given user by their user ID
func (s *AdminSupportService) FindTicketsByUserID(userID string) ([]model.SupportTicketResponse, error) {
	tickets, err := s.AdminSupportRepo.FindTicketsByUserID(userID)
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

// GetPendingTickets retrieves tickets with a pending status and supports search and pagination
func (s *AdminSupportService) GetPendingTickets(search string, page int, pageSize int) (repository.PaginatedResult, error) {
	params := repository.PaginationParams{Page: page, PageSize: pageSize}
	tickets, err := s.AdminSupportRepo.GetPendingTickets(search, params)
	if err != nil {
		return repository.PaginatedResult{}, err
	}
	return tickets, nil
}

func (s *AdminSupportService) GetClosedTickets(search string, page int, pageSize int) (repository.PaginatedResult, error) {
	params := repository.PaginationParams{Page: page, PageSize: pageSize}
	tickets, err := s.AdminSupportRepo.GetClosedTickets(search, params)
	if err != nil {
		return repository.PaginatedResult{}, err
	}
	return tickets, nil
}

func (s *AdminSupportService) saveFile(file *multipart.FileHeader, userID string) (string, error) {
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

func (s *AdminSupportService) NotifyAdmin() {}

func (s *AdminSupportService) NotifyUser() {}
