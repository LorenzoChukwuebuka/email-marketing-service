package mapper

import (
	"database/sql"
	"email-marketing-service/core/handler/user/dto"
	db "email-marketing-service/internal/db/sqlc"
	"time"
)

func MapUserResponse(r db.GetUserByIDRow) *dto.UserResponse {
	return &dto.UserResponse{
		ID:                   r.ID.String(),
		Fullname:             r.Fullname,
		CompanyID:            r.CompanyID.String(),
		Email:                r.Email,
		Phonenumber:          nullableString(r.Phonenumber),
		Password:             nullableString(r.Password),
		GoogleID:             nullableString(r.GoogleID),
		Picture:              nullableString(r.Picture),
		Verified:             r.Verified,
		Blocked:              r.Blocked,
		VerifiedAt:           nullableTime(r.VerifiedAt),
		Status:               r.Status,
		ScheduledForDeletion: r.ScheduledForDeletion,
		ScheduledDeletionAt:  nullableTime(r.ScheduledDeletionAt),
		LastLoginAt:          nullableTime(r.LastLoginAt),
		CreatedAt:            r.CreatedAt,
		UpdatedAt:            r.UpdatedAt,
		DeletedAt:            nullableTime(r.DeletedAt),
		Company: dto.CompanyResponse{
			CompanyID:        r.CompanyID_2.String(),
			Companyname:      nullableString(r.Companyname),
			CompanyCreatedAt: r.CompanyCreatedAt,
			CompanyUpdatedAt: r.CompanyUpdatedAt,
			CompanyDeletedAt: nullableTime(r.CompanyDeletedAt),
		},
	}
}

func nullableString(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func nullableTime(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}
