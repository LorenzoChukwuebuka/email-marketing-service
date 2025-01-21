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

/** this function automatically closes out a support if it has been open for more than 48 hours without any reply
1. Get all the tickets
2. loop through to get all open tickets
3 check if the current time is 48hrs or more greater than the time the ticket was created...
4. send a mail or notify them of this
5. if it is more than 72 hours and unanswered. Close the ticket.
**/

func (s *SupportTicketService) AutomaticallyCloseTickets() error {
	tickets, err := s.SupportRepository.GetAllTickets()
	if err != nil {
		return nil
	}
	for _, ticket := range tickets {

		// Convert ticket.CreatedAt to ISO 8601 format
		isoTime := ticket.CreatedAt.Format("2006-01-02T15:04:05Z07:00")

		// Print the ISO 8601 formatted time
		fmt.Println("ISO 8601 Time:", isoTime)

		// Check if the ticket is older than 48 hours
		if time.Since(ticket.CreatedAt) > 48*time.Hour && ticket.Status != model.CloseTicket  {
			// Send a mail about the current status of their ticket
			fmt.Printf("Ticket created at %s is older than 48 hours\n", isoTime)
			// Add your mail sending logic here
		}

		//check if ticket is older than 72 hours
		if time.Since(ticket.CreatedAt) > 72*time.Hour && ticket.Status != model.CloseTicket {
			//get the ticket id, userId and all and send them a mail

			fmt.Printf("Ticket created at %s is older than 72 hours\n", isoTime)
		}
	}

	return nil
}

func (s *SupportTicketService) getUploadBasePath() string {
	if os.Getenv("SERVER_MODE") == "production" {
		return "/app/backend/uploads"
	}
	return "./uploads"
}

func (s *SupportTicketService) ensureUploadDirectory(userID string) (string, error) {
	basePath := s.getUploadBasePath()
	uploadFolder := filepath.Join(basePath, "tickets", userID)

	// Create all necessary directories with proper permissions
	if err := os.MkdirAll(uploadFolder, 0777); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Ensure proper permissions on the created directory
	if err := os.Chmod(uploadFolder, 0777); err != nil {
		return "", fmt.Errorf("failed to set directory permissions: %w", err)
	}

	return uploadFolder, nil
}

func (s *SupportTicketService) saveFile(file *multipart.FileHeader, userID string) (string, error) {
	// Open source file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Ensure upload directory exists with proper permissions
	uploadFolder, err := s.ensureUploadDirectory(userID)
	if err != nil {
		return "", err
	}

	// Generate unique filename with timestamp
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	filePath := filepath.Join(uploadFolder, fileName)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Set proper file permissions
	if err := os.Chmod(filePath, 0666); err != nil {
		return "", fmt.Errorf("failed to set file permissions: %w", err)
	}

	// Copy file contents
	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	// Return relative path from base directory
	basePath := s.getUploadBasePath()
	relativePath, err := filepath.Rel(basePath, filePath)
	if err != nil {
		// Fallback to full path if relative path calculation fails
		return filePath, nil
	}

	return relativePath, nil
}

// Helper method to get full file path from relative path
func (s *SupportTicketService) getFullFilePath(relativePath string) string {
	return filepath.Join(s.getUploadBasePath(), relativePath)
}

// Add this method to help with file deletion if needed
func (s *SupportTicketService) deleteFile(relativePath string) error {
	fullPath := s.getFullFilePath(relativePath)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
