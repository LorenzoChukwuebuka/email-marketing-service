package mapper

import (
	"database/sql"
	"email-marketing-service/core/handler/support/dto"
	db "email-marketing-service/internal/db/sqlc"
	"time"

	"github.com/google/uuid"
)

// MapSupportTicketToResponse maps a single support ticket to response DTO
func MapSupportTicketToResponse(ticket db.SupportTicket) *dto.SupportTicketResponse {
	return &dto.SupportTicketResponse{
		ID:            ticket.ID.String(),
		UserID:        ticket.UserID.String(),
		Name:          ticket.Name,
		Email:         ticket.Email,
		Subject:       ticket.Subject,
		Description:   mapNullString(ticket.Description),
		TicketNumber:  ticket.TicketNumber,
		Status:        mapNullString(ticket.Status),
		Priority:      mapNullString(ticket.Priority),
		LastReply:     mapNullTime(ticket.LastReply),
		CreatedAt:     ticket.CreatedAt,
		UpdatedAt:     ticket.UpdatedAt,
		TicketFile:    []dto.TicketFile{},    // Empty slice for basic mapping
		TicketMessage: []dto.TicketMessage{}, // Empty slice for basic mapping
	}
}

// MapSupportTicketsToResponse maps multiple support tickets to response DTOs
func MapSupportTicketsToResponse(tickets []db.SupportTicket) []dto.SupportTicketResponse {
	if len(tickets) == 0 {
		return []dto.SupportTicketResponse{}
	}

	responses := make([]dto.SupportTicketResponse, len(tickets))
	for i, ticket := range tickets {
		mapped := MapSupportTicketToResponse(ticket)
		responses[i] = *mapped
	}

	return responses
}

// MapTicketWithDetailsFromRows maps GetTicketWithMessagesRow slice to SupportTicketResponse
func MapTicketWithDetailsFromRows(
	rows []db.GetTicketWithMessagesRow,
	messageFilesMap map[uuid.UUID][]db.TicketFile,
	usersMap map[uuid.UUID]db.GetUserByIDRow, // Add user data map
	adminsMap map[uuid.UUID]db.Admin, // Add admin data map
) (*dto.SupportTicketResponse, error) {
	if len(rows) == 0 {
		return nil, nil
	}

	// Build the ticket from the first row
	firstRow := rows[0]
	ticket := &dto.SupportTicketResponse{
		ID:            firstRow.TicketID.String(),
		UserID:        firstRow.TicketUserID.String(),
		Name:          firstRow.TicketName,
		Email:         firstRow.TicketEmail,
		Subject:       firstRow.TicketSubject,
		Description:   mapNullString(firstRow.TicketDescription),
		TicketNumber:  firstRow.TicketNumber,
		Status:        mapNullString(firstRow.TicketStatus),
		Priority:      mapNullString(firstRow.TicketPriority),
		LastReply:     mapNullTime(firstRow.TicketLastReply),
		CreatedAt:     firstRow.TicketCreatedAt,
		UpdatedAt:     firstRow.TicketUpdatedAt,
		TicketFile:    []dto.TicketFile{},
		TicketMessage: []dto.TicketMessage{},
	}

	// Track processed messages to avoid duplicates
	processedMessages := make(map[uuid.UUID]bool)
	var allFiles []dto.TicketFile

	// Process messages from rows
	for _, row := range rows {
		if row.MessageID.Valid && !processedMessages[row.MessageID.UUID] {
			message := dto.TicketMessage{
				ID:        row.MessageID.UUID.String(),
				TicketID:  row.TicketID.String(),
				UserID:    row.MessageUserID.UUID.String(),
				Message:   row.Message.String,
				IsAdmin:   row.IsAdmin.Bool,
				CreatedAt: row.MessageCreatedAt.Time,
				UpdatedAt: row.MessageUpdatedAt.Time,
			}

			// Populate user/admin information based on IsAdmin flag
			if row.IsAdmin.Bool {
				// Get admin information
				if admin, exists := adminsMap[row.MessageUserID.UUID]; exists {
					message.Admin = &dto.AdminResponse{
						ID:        admin.ID.String(),
						Firstname: admin.Firstname.String,
						Lastname:  admin.Lastname.String,
						Type:      admin.Type,
					}
				}
			} else {
				// Get user information
				if user, exists := usersMap[row.MessageUserID.UUID]; exists {
					message.User = &dto.UserResponse{
						ID:         user.ID.String(),
						Fullname:   user.Fullname,
						Email:      user.Email,
						Verified:   user.Verified,
						Blocked:    user.Blocked,
						VerifiedAt: &user.VerifiedAt.Time,
						CreatedAt:  user.CreatedAt,
					}
				}
			}

			ticket.TicketMessage = append(ticket.TicketMessage, message)
			processedMessages[row.MessageID.UUID] = true

			// Collect files for this message
			if files, exists := messageFilesMap[row.MessageID.UUID]; exists {
				for _, file := range files {
					allFiles = append(allFiles, dto.TicketFile{
						ID:        file.ID.String(),
						MessageID: file.MessageID.String(),
						FileName:  file.FileName,
						FilePath:  file.FilePath,
					})
				}
			}
		}
	}

	ticket.TicketFile = allFiles
	return ticket, nil
}

// Helper functions for null value handling
func mapNullString(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

func mapNullTime(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}
