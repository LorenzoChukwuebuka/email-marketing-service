package adminrepository

import (
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	"fmt"

	"gorm.io/gorm"
)

type AdminSupportRepository struct {
	DB *gorm.DB
}

func NewAdminSupportRepository(db *gorm.DB) *AdminSupportRepository {
	return &AdminSupportRepository{
		DB: db,
	}
}

func (r *AdminSupportRepository) FindTicketsByUserID(userID string) ([]model.SupportTicketResponse, error) {
	var tickets []model.SupportTicket
	if err := r.DB.Where("user_id = ?", userID).Find(&tickets).Error; err != nil {
		return nil, err
	}
	return r.modelsToResponses(tickets), nil
}

func (r *AdminSupportRepository) GetAllTickets(search string, params repository.PaginationParams) (repository.PaginatedResult, error) {
	var tickets []model.SupportTicket

	query := r.DB.Preload("Messages.Files").Find(&tickets)

	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ? OR subject ILIKE ? OR ticket_number ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")

	}

	pagination, err := repository.Paginate(query, params, &tickets)

	if err != nil {
		return repository.PaginatedResult{}, fmt.Errorf("failed to paginate contacts: %w", err)
	}

	pagination.Data = r.modelsToResponses(tickets)

	return pagination, nil
}

func (r *AdminSupportRepository) GetPendingTickets(search string, params repository.PaginationParams) (repository.PaginatedResult, error) {
	var tickets []model.SupportTicket

	query := r.DB.Preload("Messages.Files").Where("status = ?", model.PendingTicket).Order("created_at DESC").Find(&tickets)

	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ? OR subject ILIKE ? OR ticket_number ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")

	}

	pagination, err := repository.Paginate(query, params, &tickets)

	if err != nil {
		return repository.PaginatedResult{}, fmt.Errorf("failed to paginate contacts: %w", err)
	}

	pagination.Data = r.modelsToResponses(tickets)

	return pagination, nil
}

func (r *AdminSupportRepository) GetClosedTickets(search string, params repository.PaginationParams) (repository.PaginatedResult, error) {
	var tickets []model.SupportTicket

	query := r.DB.Preload("Messages.Files").Where("status = ?", model.CloseTicket).Order("created_at DESC").Find(&tickets)

	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ? OR subject ILIKE ? OR ticket_number ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")

	}

	pagination, err := repository.Paginate(query, params, &tickets)

	if err != nil {
		return repository.PaginatedResult{}, fmt.Errorf("failed to paginate contacts: %w", err)
	}

	pagination.Data = r.modelsToResponses(tickets)

	return pagination, nil
}

func (r *AdminSupportRepository) modelsToResponses(tickets []model.SupportTicket) []model.SupportTicketResponse {
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

func (r *AdminSupportRepository) messagesToResponses(messages []model.TicketMessage) []model.TicketMessageResponse {
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

func (r *AdminSupportRepository) filesToResponses(files []model.TicketFile) []model.TicketFileResponse {
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
