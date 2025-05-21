package mapper

import (
	"email-marketing-service/core/handler/admin/auth/dto"
	db "email-marketing-service/internal/db/sqlc"
	"time"
)



func MapAdminToResponse(admin db.Admin) dto.AdminResponse {
	var deletedAt *time.Time
	if admin.DeletedAt.Valid {
		deletedAt = &admin.DeletedAt.Time
	}

	return dto.AdminResponse{
		ID:         admin.ID.String(),
		Firstname:  admin.Firstname.String,
		Middlename: admin.Middlename.String,
		Lastname:   admin.Lastname.String,
		Email:      admin.Email,
		Type:       admin.Type,
		CreatedAt:  admin.CreatedAt,
		UpdatedAt:  admin.UpdatedAt.Time, 
		DeletedAt:  deletedAt,
	}
}
