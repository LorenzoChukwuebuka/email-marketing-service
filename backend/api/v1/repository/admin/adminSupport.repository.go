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

func (r *AdminSupportRepository) CreateTicketMessage(message *model.TicketMessage) (*model.TicketMessageResponse, error) {
	if err := r.DB.Create(message).Error; err != nil {
		return nil, err
	}
	return r.messageToResponse(message), nil
}

func (r *AdminSupportRepository) CreateTicketFile(file *model.TicketFile) (*model.TicketFileResponse, error) {
	if err := r.DB.Create(file).Error; err != nil {
		return nil, err
	}
	return r.fileToResponse(file), nil
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

	query := r.DB.Preload("Files").Preload("Messages").Find(&tickets)

	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ? OR subject ILIKE ? OR ticket_number ILIKE ? ", "%"+search+"%", "%"+search+"%", "%"+search+"%"+"%"+search+"%")
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

	query := r.DB.Preload("Files").Preload("Messages").Where("status = ?", model.SupportStatus("pending")).Order("created_at DESC").Find(&tickets)

	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ? OR subject ILIKE ? OR ticket_number ILIKE ? ", "%"+search+"%", "%"+search+"%", "%"+search+"%"+"%"+search+"%")
	}

	pagination, err := repository.Paginate(query, params, &tickets)

	if err != nil {
		return repository.PaginatedResult{}, fmt.Errorf("failed to paginate contacts: %w", err)
	}

	pagination.Data = r.modelsToResponses(tickets)

	return pagination, nil
}

func (r *AdminSupportRepository) modelToResponse(ticket *model.SupportTicket) *model.SupportTicketResponse {
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
		Files:        r.filesToResponses(ticket.Files),
		Messages:     r.messagesToResponses(ticket.Messages),
		CreatedAt:    ticket.CreatedAt,
		UpdatedAt:    ticket.UpdatedAt,
	}
}

func (r *AdminSupportRepository) modelsToResponses(tickets []model.SupportTicket) []model.SupportTicketResponse {
	responses := make([]model.SupportTicketResponse, len(tickets))
	for i, ticket := range tickets {
		responses[i] = *r.modelToResponse(&ticket)
	}
	return responses
}

func (r *AdminSupportRepository) fileToResponse(file *model.TicketFile) *model.TicketFileResponse {
	return &model.TicketFileResponse{
		UUID:      file.UUID,
		FileName:  file.FileName,
		CreatedAt: file.CreatedAt,
	}
}

func (r *AdminSupportRepository) filesToResponses(files []model.TicketFile) []model.TicketFileResponse {
	responses := make([]model.TicketFileResponse, len(files))
	for i, file := range files {
		responses[i] = *r.fileToResponse(&file)
	}
	return responses
}

func (r *AdminSupportRepository) messageToResponse(message *model.TicketMessage) *model.TicketMessageResponse {
	return &model.TicketMessageResponse{
		UUID:      message.UUID,
		UserID:    message.UserID,
		Message:   message.Message,
		IsAdmin:   message.IsAdmin,
		CreatedAt: message.CreatedAt,
	}
}

func (r *AdminSupportRepository) messagesToResponses(messages []model.TicketMessage) []model.TicketMessageResponse {
	responses := make([]model.TicketMessageResponse, len(messages))
	for i, message := range messages {
		responses[i] = *r.messageToResponse(&message)
	}
	return responses
}
