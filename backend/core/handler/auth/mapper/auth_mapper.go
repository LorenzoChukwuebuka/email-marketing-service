package mapper

import (
	"email-marketing-service/core/handler/auth/dto"
	db "email-marketing-service/internal/db/sqlc"
	"github.com/google/uuid"
	"time"
)

func MapPublicUser(u db.GetUserByEmailRow) dto.PublicUser {
	return dto.PublicUser{
		ID:        u.ID,
		Fullname:  u.Fullname,
		Email:     u.Email,
		CompanyID: u.CompanyID,
		Verified:  u.Verified,
		Blocked:   u.Blocked,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Company: struct {
			ID          uuid.UUID `json:"id"`
			CompanyName string    `json:"company_name"`
			CreatedAt   time.Time `json:"created_at"`
			UpdatedAt   time.Time `json:"updated_at"`
		}{
			ID:          u.CompanyID,
			CompanyName: u.Companyname.String,
			CreatedAt:   u.CompanyCreatedAt,
			UpdatedAt:   u.CompanyUpdatedAt,
		},
	}
}
