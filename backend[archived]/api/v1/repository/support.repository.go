package repository

import (
	"email-marketing-service/api/v1/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type SupportRepository struct {
	DB *gorm.DB
}

func NewSupportRepository(db *gorm.DB) *SupportRepository {
	return &SupportRepository{DB: db}
}

func (r *SupportRepository) CreateSupportTicket(ticket *model.SupportTicket) (*model.SupportTicketResponse, error) {
	if err := r.DB.Create(ticket).Error; err != nil {
		return nil, err
	}
	return r.modelToResponse(ticket), nil
}

func (r *SupportRepository) FindTicketByID(id string) (*model.SupportTicketResponse, error) {
	var ticket model.SupportTicket
	if err := r.DB.Preload("Messages.Files").Where("uuid = ?", id).First(&ticket).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ticket not found")
		}
		return nil, err
	}
	return r.modelToResponse(&ticket), nil
}

func (r *SupportRepository) UpdateTicket(ticket *model.SupportTicket) (*model.SupportTicketResponse, error) {
	var existingTicket model.SupportTicket
	if err := r.DB.Where("id = ?", ticket.ID).First(&existingTicket).Error; err != nil {
		return nil, fmt.Errorf("failed to find contact for update: %w", err)
	}

	if ticket.Status != "" {
		existingTicket.Status = ticket.Status
	}

	if ticket.LastReply != nil {
		existingTicket.LastReply = ticket.LastReply
	}

	if err := r.DB.Save(&existingTicket).Error; err != nil {
		return nil, err
	}
	return r.modelToResponse(ticket), nil
}

func (r *SupportRepository) CreateTicketMessage(message *model.TicketMessage) (*model.TicketMessageResponse, error) {
	if err := r.DB.Create(message).Error; err != nil {
		return nil, err
	}
	return r.messageToResponse(message), nil
}

func (r *SupportRepository) CreateTicketFile(file *model.TicketFile) (*model.TicketFileResponse, error) {
	if err := r.DB.Create(file).Error; err != nil {
		return nil, err
	}
	return r.fileToResponse(file), nil
}

func (r *SupportRepository) FindTicketsByUserID(userID string) ([]model.SupportTicketResponse, error) {
	var tickets []model.SupportTicket
	if err := r.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&tickets).Error; err != nil {
		return nil, err
	}
	return r.modelsToResponses(tickets), nil
}

func (r *SupportRepository) FindOpenTickets() ([]model.SupportTicketResponse, error) {
	var tickets []model.SupportTicket
	if err := r.DB.Where("status = ?", model.OpenTicket).Find(&tickets).Error; err != nil {
		return nil, err
	}
	return r.modelsToResponses(tickets), nil
}

func (r *SupportRepository) FindTicketMessagesByTicketID(ticketID uint) ([]model.TicketMessageResponse, error) {
	var messages []model.TicketMessage
	if err := r.DB.Where("ticket_id = ?", ticketID).Order("created_at ASC").Find(&messages).Error; err != nil {
		return nil, err
	}
	return r.messagesToResponses(messages), nil
}

func (r *SupportRepository) GetTicketWithDetails(ticketID string) (*model.SupportTicketResponse, error) {
	var ticket model.SupportTicket
	if err := r.DB.Preload("Messages.Files").Where("uuid = ?", ticketID).First(&ticket).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ticket not found")
		}
		return nil, err
	}
	return r.modelToResponse(&ticket), nil
}

func (r *SupportRepository) GetTicketsByUserID(userID string) ([]model.SupportTicketResponse, error) {
	var tickets []model.SupportTicket
	if err := r.DB.Preload("Messages.Files").Where("user_id = ?", userID).Order("created_at DESC").Find(&tickets).Error; err != nil {
		return nil, err
	}
	return r.modelsToResponses(tickets), nil
}

func (r *SupportRepository) CloseTicket(ticketId string) error {
	// Perform the update using GORM's Save method
	result := r.DB.Model(&model.SupportTicket{}).Where("uuid = ?  ", ticketId).Updates(map[string]interface{}{
		"status": model.CloseTicket,
	})

	// Check for any errors during the update
	if result.Error != nil {
		return fmt.Errorf("failed to close ticket: %w", result.Error)
	}

	// If no rows were affected, return an error indicating no record was found
	if result.RowsAffected == 0 {
		return fmt.Errorf("no ticket record found to update")
	}

	return nil
}

func (r *SupportRepository) GetAllTickets() ([]model.SupportTicketResponse, error) {
	var tickets []model.SupportTicket
	if err := r.DB.Preload("Messages.Files").Order("created_at DESC").Find(&tickets).Error; err != nil {
		return nil, err
	}
	return r.modelsToResponses(tickets), nil
}

func (r *SupportRepository) modelToResponse(ticket *model.SupportTicket) *model.SupportTicketResponse {
	return &model.SupportTicketResponse{
		ID:           ticket.ID,
		UUID:         ticket.UUID,
		UserID:       ticket.UserID,
		Name:         ticket.Name,
		Email:        ticket.Email,
		Subject:      ticket.Subject,
		Description:  ticket.Description,
		TicketNumber: ticket.TicketNumber,
		Status:       ticket.Status,
		Priority:     ticket.Priority,
		LastReply:    ticket.LastReply,
		Messages:     r.messagesToResponses(ticket.Messages),
		CreatedAt:    ticket.CreatedAt,
		UpdatedAt:    ticket.UpdatedAt,
	}
}

func (r *SupportRepository) modelsToResponses(tickets []model.SupportTicket) []model.SupportTicketResponse {
	responses := make([]model.SupportTicketResponse, len(tickets))
	for i, ticket := range tickets {
		responses[i] = model.SupportTicketResponse{
			ID:           ticket.ID,
			UUID:         ticket.UUID,
			UserID:       ticket.UserID,
			Name:         ticket.Name,
			Email:        ticket.Email,
			Subject:      ticket.Subject,
			Description:  ticket.Description,
			TicketNumber: ticket.TicketNumber,
			Status:       ticket.Status,
			Priority:     ticket.Priority,
			LastReply:    ticket.LastReply,
			Messages:     r.messagesToResponses(ticket.Messages),
			CreatedAt:    ticket.CreatedAt,
			UpdatedAt:    ticket.UpdatedAt,
		}
	}
	return responses
}

func (r *SupportRepository) messagesToResponses(messages []model.TicketMessage) []model.TicketMessageResponse {
	responses := make([]model.TicketMessageResponse, len(messages))
	for i, message := range messages {
		responses[i] = model.TicketMessageResponse{
			ID:        message.ID,
			UUID:      message.UUID,
			UserID:    message.UserID,
			Message:   message.Message,
			IsAdmin:   message.IsAdmin,
			CreatedAt: message.CreatedAt,
			Files:     r.filesToResponses(message.Files),
		}
	}
	return responses
}

func (r *SupportRepository) filesToResponses(files []model.TicketFile) []model.TicketFileResponse {
	responses := make([]model.TicketFileResponse, len(files))
	for i, file := range files {
		responses[i] = model.TicketFileResponse{
			UUID:      file.UUID,
			FileName:  file.FileName,
			FilePath:  file.FilePath,
			CreatedAt: file.CreatedAt,
		}
	}
	return responses
}

func (r *SupportRepository) fileToResponse(file *model.TicketFile) *model.TicketFileResponse {
	return &model.TicketFileResponse{
		UUID:      file.UUID,
		FileName:  file.FileName,
		FilePath:  file.FilePath,
		CreatedAt: file.CreatedAt,
	}
}

func (r *SupportRepository) messageToResponse(message *model.TicketMessage) *model.TicketMessageResponse {
	return &model.TicketMessageResponse{
		ID:        message.ID,
		UUID:      message.UUID,
		UserID:    message.UserID,
		Message:   message.Message,
		IsAdmin:   message.IsAdmin,
		CreatedAt: message.CreatedAt,
		Files:     r.filesToResponses(message.Files),
	}
}
