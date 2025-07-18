package mapper

import (
	"email-marketing-service/core/handler/admin/support/dto"
	db "email-marketing-service/internal/db/sqlc"
)

func MapSupportTicketsToResponse(tickets []db.SupportTicket) []*dto.SupportTicketResponse {

	if tickets == nil {
		return nil
	}

	responses := make([]*dto.SupportTicketResponse, 0, len(tickets))

	for _, ticket := range tickets {
		response := &dto.SupportTicketResponse{
			ID:           ticket.ID.String(),
			UserID:       ticket.UserID.String(),
			Name:         ticket.Name,
			Email:        ticket.Email,
			Subject:      ticket.Subject,
			TicketNumber: ticket.TicketNumber,
			CreatedAt:    ticket.CreatedAt,
			UpdatedAt:    ticket.UpdatedAt,
		}

		// Handle nullable description
		if ticket.Description.Valid {
			response.Description = &ticket.Description.String
		}

		// Handle nullable status
		if ticket.Status.Valid {
			response.Status = &ticket.Status.String
		}

		// Handle nullable priority
		if ticket.Priority.Valid {
			response.Priority = &ticket.Priority.String
		}

		// Handle nullable last_reply
		if ticket.LastReply.Valid {
			response.LastReply = &ticket.LastReply.Time
		}

		responses = append(responses, response)
	}

	return responses

}
