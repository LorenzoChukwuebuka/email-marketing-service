package services

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/model"
	adminmodel "email-marketing-service/api/v1/model/admin"
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

type SupportTicketService struct {
	SupportRepository     *repository.SupportRepository
	UserRepository        *repository.UserRepository
	UserNotification      *repository.UserNotificationRepository
	AdminNotificationRepo *adminrepository.AdminNotificationRepository
	AdminRepo             *adminrepository.AdminRepository
}

func NewSupportTicketService(supportRepo *repository.SupportRepository,
	userRepo *repository.UserRepository,
	userNotificationRepo *repository.UserNotificationRepository, adminNotificationRepo *adminrepository.AdminNotificationRepository, adminRepo *adminrepository.AdminRepository) *SupportTicketService {
	return &SupportTicketService{
		SupportRepository:     supportRepo,
		UserRepository:        userRepo,
		UserNotification:      userNotificationRepo,
		AdminNotificationRepo: adminNotificationRepo,
		AdminRepo:             adminRepo,
	}
}

func (s *SupportTicketService) CreateSupportTicket(userID string, req *dto.CreateSupportTicketRequest) (*dto.CreateSupportTicketResponse, error) {
	// Start a database transaction
	tx := s.SupportRepository.DB.Begin()

	// Ensure the transaction rolls back if there's a panic or unexpected error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Find the user by their userID
	userModel := &model.User{UUID: userID}
	user, err := s.UserRepository.FindUserById(userModel)
	if err != nil {
		// Rollback the transaction if user not found
		tx.Rollback()
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Initialize the SupportTicket model with details from the request and user info
	ticket := &model.SupportTicket{
		UserID:       userID,                               // Set the user ID
		Name:         user.FullName,                        // Set the user's full name
		Email:        user.Email,                           // Set the user's email
		Subject:      req.Subject,                          // Set the ticket's subject
		TicketNumber: utils.GenerateUniqueRandomNumbers(8), // Generate a random ticket number
		Description:  req.Description,                      // Set the ticket description
		Status:       model.PendingTicket,                  // Set initial status as 'pending'
		Priority:     model.Priority(req.Priority),         // Set the ticket priority from request
	}

	// Create the initial message for the support ticket
	initialMessage := model.TicketMessage{
		UserID:  userID,      // Set the user ID who created the ticket
		Message: req.Message, // Set the initial message content
		IsAdmin: false,       // Mark the message as not from an admin
	}

	// Check if there are any files attached, if yes, save each one
	if req.File != nil {
		for _, file := range req.File {
			// Save the file to the filesystem
			filePath, err := s.saveFile(file, userID)
			if err != nil {
				tx.Rollback() // Rollback if file save fails
				return nil, fmt.Errorf("failed to save file: %w", err)
			}
			// Associate the saved file with the message
			initialMessage.Files = append(initialMessage.Files, model.TicketFile{
				FileName: file.Filename, // Store the original file name
				FilePath: filePath,      // Store the file path where it was saved
			})
		}
	}

	// Attach the message (and files if any) to the ticket
	ticket.Messages = []model.TicketMessage{initialMessage}

	// Attempt to create the support ticket in the repository
	createdTicket, err := s.SupportRepository.CreateSupportTicket(ticket)
	if err != nil {
		tx.Rollback() // Rollback the transaction if ticket creation fails
		return nil, fmt.Errorf("failed to create support ticket: %w", err)
	}

	// Commit the transaction once all steps are successful
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Send a notification to the user indicating that their ticket was created
	notificationTitle := fmt.Sprintf("You have opened a ticket with subject %s. Our agent will respond to you soon", createdTicket.Subject)
	if err := utils.CreateNotification(s.UserNotification, userID, notificationTitle); err != nil {
		fmt.Printf("Failed to create notification: %v\n", err)
	}

	// Return a success response with the created ticket ID and a message
	return &dto.CreateSupportTicketResponse{
		TicketID: createdTicket.ID,                                                              // Return the newly created ticket ID
		Message:  "Your ticket has been successfully created. An agent will reply to you soon.", // Success message
	}, nil
}

func (s *SupportTicketService) ReplyToTicket(ticketID string, userID string, req *dto.ReplyTicketRequest) (*dto.ReplyTicketResponse, error) {
	// Find the ticket by its ID
	ticket, err := s.SupportRepository.FindTicketByID(ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to find ticket: %w", err)
	}

	// Create a new ticket message object with the user reply
	message := &model.TicketMessage{
		TicketID: ticket.ID,   // Associate the message with the ticket
		UserID:   userID,      // Set the ID of the user who is replying
		Message:  req.Message, // Set the content of the reply message
		IsAdmin:  false,       // Mark the message as from a user (not admin)
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

	// Send a notification to the user indicating their reply has been received
	notificationTitle := fmt.Sprintf("You have replied to the ticket with subject %s. Our agent will respond to you soon", ticket.Subject)
	if err := utils.CreateNotification(s.UserNotification, userID, notificationTitle); err != nil {
		fmt.Printf("Failed to create notification: %v\n", err)
	}

	// Return a success response indicating that the reply has been added to the ticket
	return &dto.ReplyTicketResponse{
		MessageID: msg.ID,                                                  // Return the ID of the newly created message
		Message:   "Your reply has been successfully added to the ticket.", // Success message
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

	//Populate user information for non-admin messages
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

		//for admin
		if message.IsAdmin {
			admin, err := s.AdminRepo.FindAdminById(message.UserID)

			if err != nil {
				return nil, fmt.Errorf("failed to get admin for message: %w", err)
			}

			ticket.Messages[i].Admin = &adminmodel.AdminResponse{
				UUID:      admin.UUID,
				FirstName: admin.FirstName,
				LastName:  admin.LastName,
				Type:      admin.Type,
			}

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

func (s *SupportTicketService) CloseTicket(ticketId string) error {

	if err := s.SupportRepository.CloseTicket(ticketId); err != nil {
		return err
	}
	return nil
}
