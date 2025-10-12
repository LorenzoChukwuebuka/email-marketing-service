package service

import (
	"context"
	"database/sql"
	adminmapper "email-marketing-service/core/handler/admin/support/mapper"
	"email-marketing-service/core/handler/support/dto"
	"email-marketing-service/core/handler/support/mapper"
	"email-marketing-service/internal/common"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/enums"
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type AdminSupportService struct {
	store db.Store
}

func NewAdminSupportService(store db.Store) *AdminSupportService {
	return &AdminSupportService{store: store}
}

func (s *AdminSupportService) ReplyToTicket(ctx context.Context, ticketID string, adminID string, req *dto.ReplyTicketRequest) (*dto.ReplyTicketResponse, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"ticket": ticketID,
		"admin":  adminID,
	})
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	var result *dto.ReplyTicketResponse

	// Start a transaction
	err = s.store.ExecTx(ctx, func(q *db.Queries) error {
		// Find the ticket by its ID
		ticket, err := q.FindTicketByID(ctx, _uuid["ticket"])
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("ticket not found")
			}
			return fmt.Errorf("failed to find ticket: %w", err)
		}

		// Create a new ticket message with the admin reply
		message, err := q.CreateTicketMessage(ctx, db.CreateTicketMessageParams{
			TicketID: ticket.ID,
			UserID:   _uuid["admin"],
			Message:  req.Message,
			IsAdmin:  sql.NullBool{Bool: true, Valid: true}, // Mark as admin reply
		})
		if err != nil {
			return fmt.Errorf("failed to create ticket message: %w", err)
		}

		// Check if there are any files attached in the reply request
		if req.File != nil {
			for _, file := range req.File {
				// Save the file to the filesystem
				filePath, err := s.saveFile(file, adminID)
				if err != nil {
					return fmt.Errorf("failed to save file: %w", err)
				}

				// Create a ticket file record in the database
				_, err = q.CreateTicketFile(ctx, db.CreateTicketFileParams{
					MessageID: message.ID,
					FileName:  file.Filename,
					FilePath:  filePath,
				})
				if err != nil {
					return fmt.Errorf("failed to create ticket file record: %w", err)
				}
			}
		}

		// Update the ticket's status and last reply time
		now := time.Now()
		_, err = q.UpdateTicketStatus(ctx, db.UpdateTicketStatusParams{
			ID:        ticket.ID,
			Status:    sql.NullString{String: string(enums.OpenTicket), Valid: true},
			LastReply: sql.NullTime{Time: now, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to update ticket: %w", err)
		}

		// Send notification to the user
		notificationTitle := fmt.Sprintf("Your ticket with subject '%s' has been replied to by our support team", ticket.Subject)
		_, err = q.CreateUserNotification(ctx, db.CreateUserNotificationParams{
			UserID:          ticket.UserID,
			Title:           notificationTitle,
			AdditionalField: "support_reply",
		})
		if err != nil {
			return fmt.Errorf("failed to create user notification: %w", err)
		}

		// Send notification to the admin
		adminNotificationTitle := fmt.Sprintf("You have replied to ticket '%s'", ticket.Subject)
		_, err = q.CreateAdminNotification(ctx, db.CreateAdminNotificationParams{
			UserID:  ticket.UserID,
			Title:  adminNotificationTitle,
		})
		if err != nil {
			return fmt.Errorf("failed to create admin notification: %w", err)
		}

		// Set the result to return
		result = &dto.ReplyTicketResponse{
			MessageID: uint(message.ID.ID()), // Convert UUID to uint if needed
			Message:   "Your reply has been successfully added to the ticket.",
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *AdminSupportService) GetAllTickets(ctx context.Context, search string, page int, pageSize int) (any, error) {
	result, err := s.store.GetAllTicketsWithPagination(ctx, db.GetAllTicketsWithPaginationParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(page),
	})

	if err != nil {
		return nil, common.ErrFetchingRecord
	}

	ticket_count, err := s.store.GetAllTicketsCount(ctx)
	if err != nil {
		return nil, common.ErrFetchingCount
	}

	data := adminmapper.MapSupportTicketsToResponse(result)

	items := make([]any, len(data))

	for i, v := range data {
		items[i] = v
	}

	response := common.Paginate(int(ticket_count), items, page, pageSize)

	return response, nil
}

func (s *AdminSupportService) GetPendingTickets(ctx context.Context, search string, page int, pageSize int) (any, error) {
	result, err := s.store.GetPendingTicketsWithPagination(ctx, db.GetPendingTicketsWithPaginationParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(page),
	})
	if err != nil {
		return nil, common.ErrFetchingRecord
	}

	ticket_count, err := s.store.GetPendingTicketsCount(ctx)
	if err != nil {
		return nil, common.ErrFetchingCount
	}

	data := adminmapper.MapSupportTicketsToResponse(result)

	items := make([]any, len(data))

	for i, v := range data {
		items[i] = v
	}

	response := common.Paginate(int(ticket_count), items, page, pageSize)

	return response, nil
}

func (s *AdminSupportService) GetClosedTickets(ctx context.Context, search string, page int, pageSize int) (any, error) {
	result, err := s.store.GetClosedTicketsWithPagination(ctx, db.GetClosedTicketsWithPaginationParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(page),
	})
	if err != nil {
		return nil, common.ErrFetchingRecord
	}

	ticket_count, err := s.store.GetClosedTicketsCount(ctx)
	if err != nil {
		return nil, common.ErrFetchingCount
	}

	data := adminmapper.MapSupportTicketsToResponse(result)
	items := make([]any, len(data))

	for i, v := range data {
		items[i] = v
	}

	response := common.Paginate(int(ticket_count), items, page, pageSize)

	return response, nil
}

func (s *AdminSupportService) GetTicketWithDetails(ctx context.Context, ticketID string) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"ticket": ticketID,
	})
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	// Get ticket with messages (reusing the same logic as user service)
	rows, err := s.store.GetTicketWithMessages(ctx, _uuid["ticket"])
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket with messages: %w", err)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("ticket not found")
	}

	// Get all message IDs to fetch files
	messageIDs := make([]uuid.UUID, 0)
	messageIDSet := make(map[uuid.UUID]bool)

	for _, row := range rows {
		if row.MessageID.Valid && !messageIDSet[row.MessageID.UUID] {
			messageIDs = append(messageIDs, row.MessageID.UUID)
			messageIDSet[row.MessageID.UUID] = true
		}
	}

	// Get files for all messages
	messageFilesMap := make(map[uuid.UUID][]db.TicketFile)
	for _, messageID := range messageIDs {
		files, err := s.store.GetMessageFiles(ctx, messageID)
		if err == nil {
			messageFilesMap[messageID] = files
		}
	}

	userIDs := make([]uuid.UUID, 0)
	adminIDs := make([]uuid.UUID, 0)
	userIDSet := make(map[uuid.UUID]bool)
	adminIDSet := make(map[uuid.UUID]bool)

	for _, row := range rows {
		if row.MessageID.Valid {
			if row.IsAdmin.Bool && !adminIDSet[row.MessageUserID.UUID] {
				adminIDs = append(adminIDs, row.MessageUserID.UUID)
				adminIDSet[row.MessageUserID.UUID] = true
			} else if !row.IsAdmin.Bool && !userIDSet[row.MessageUserID.UUID] {
				userIDs = append(userIDs, row.MessageUserID.UUID)
				userIDSet[row.MessageUserID.UUID] = true
			}
		}
	}

	// Fetch users and admins
	usersMap := make(map[uuid.UUID]db.GetUserByIDRow)
	adminsMap := make(map[uuid.UUID]db.Admin)

	for _, userID := range userIDs {
		user, err := s.store.GetUserByID(ctx, userID)
		if err == nil {
			usersMap[userID] = user
		}
	}

	for _, adminID := range adminIDs {
		admin, err := s.store.GetAdminByID(ctx, adminID)
		if err == nil {
			adminsMap[adminID] = admin
		}
	}

	// Use mapper to convert to response
	ticket, err := mapper.MapTicketWithDetailsFromRows(rows, messageFilesMap, usersMap, adminsMap)
	if err != nil {
		return nil, fmt.Errorf("failed to map ticket details: %w", err)
	}

	return ticket, nil
}

func (s *AdminSupportService) FindTicketsByUserID(ctx context.Context, userID string) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"user": userID,
	})
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	tickets, err := s.store.GetTicketsByUserID(ctx, _uuid["user"])
	if err != nil {
		return nil, fmt.Errorf("failed to get user tickets: %w", err)
	}

	data := mapper.MapSupportTicketsToResponse(tickets)
	return data, nil
}

func (s *AdminSupportService) CloseTicket(ctx context.Context, ticketID string) error {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"ticket": ticketID,
	})
	if err != nil {
		return common.ErrInvalidUUID
	}

	_, err = s.store.CloseTicketByID(ctx, _uuid["ticket"])
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("ticket not found")
		}
		return fmt.Errorf("failed to close ticket: %w", err)
	}

	return nil
}

func (s *AdminSupportService) CloseStaleTickets(ctx context.Context) (any, error) {
	tickets, err := s.store.CloseStaleTickets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to close stale tickets: %w", err)
	}

	return map[string]interface{}{
		"closed_tickets_count": len(tickets),
		"message":              fmt.Sprintf("Successfully closed %d stale tickets", len(tickets)),
	}, nil
}

// File handling methods (copied from user service)
func (s *AdminSupportService) getUploadBasePath() string {
	if os.Getenv("SERVER_MODE") == "production" {
		return "/app/backend/uploads"
	}
	return "./uploads"
}

func (s *AdminSupportService) ensureUploadDirectory(adminID string) (string, error) {
	basePath := s.getUploadBasePath()
	uploadFolder := filepath.Join(basePath, "tickets", "admin", adminID)

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

func (s *AdminSupportService) saveFile(file *multipart.FileHeader, adminID string) (string, error) {
	// Open source file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Ensure upload directory exists with proper permissions
	uploadFolder, err := s.ensureUploadDirectory(adminID)
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
func (s *AdminSupportService) GetFullFilePath(relativePath string) string {
	return filepath.Join(s.getUploadBasePath(), relativePath)
}

// Add this method to help with file deletion if needed
func (s *AdminSupportService) DeleteFile(relativePath string) error {
	fullPath := s.GetFullFilePath(relativePath)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
