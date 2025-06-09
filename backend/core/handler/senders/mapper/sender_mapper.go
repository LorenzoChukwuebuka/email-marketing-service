package mapper

import (
	"email-marketing-service/core/handler/senders/dto"
	db "email-marketing-service/internal/db/sqlc"
	"time"
)

func MapSenderToDTO(sender *db.Sender) *dto.SendersResponse {
	var createdAt, updatedAt, deletedAt *time.Time
	var domainId string

	// Handle nullable time fields
	if sender.CreatedAt.Valid {
		createdAt = &sender.CreatedAt.Time
	}

	if sender.UpdatedAt.Valid {
		updatedAt = &sender.UpdatedAt.Time
	}

	if sender.DeletedAt.Valid {
		deletedAt = &sender.DeletedAt.Time
	}

	if sender.DomainID.Valid {
		domainId = sender.DomainID.UUID.String()
	}

	return &dto.SendersResponse{
		ID:        sender.ID.String(),
		UserID:    sender.UserID.String(),
		CompanyID: sender.CompanyID.String(),
		Name:      sender.Name,
		Email:     sender.Email,
		Verified:  sender.Verified.Bool, // This will be false if null
		IsSigned:  sender.IsSigned.Bool, // This will be false if null
		DomainID:  domainId,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
	}
}

func MapSendersToDTO(senders []db.Sender) []*dto.SendersResponse {
	if len(senders) == 0 {
		return []*dto.SendersResponse{}
	}

	result := make([]*dto.SendersResponse, len(senders))
	for i, sender := range senders {
		result[i] = MapSenderToDTO(&sender)
	}

	return result
}
