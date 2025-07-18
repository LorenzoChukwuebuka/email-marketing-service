// mapper/user_mapper.go
package mapper

import (
	"database/sql"
	"email-marketing-service/core/handler/admin/users/dto"
	db "email-marketing-service/internal/db/sqlc"
	"time"

	"github.com/google/uuid"
)

type MappableUser interface {
	GetID() uuid.UUID
	GetFullname() string
	GetEmail() string
	GetPhonenumber() sql.NullString
	GetPicture() sql.NullString
	GetVerified() bool
	GetBlocked() bool
	GetVerifiedAt() sql.NullTime
	GetStatus() string
	GetScheduledForDeletion() bool
	GetScheduledDeletionAt() sql.NullTime
	GetLastLoginAt() sql.NullTime
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	GetCompanyID() uuid.NullUUID
	GetCompanyname() sql.NullString
}

func MapAdminUsers[T MappableUser](users []T) []*dto.AdminUserResponseDTO {
	result := make([]*dto.AdminUserResponseDTO, 0, len(users))

	for _, user := range users {
		result = append(result, &dto.AdminUserResponseDTO{
			ID:                   user.GetID().String(),
			Fullname:             user.GetFullname(),
			Email:                user.GetEmail(),
			Phonenumber:          nullableString(user.GetPhonenumber()),
			Picture:              nullableString(user.GetPicture()),
			Verified:             user.GetVerified(),
			Blocked:              user.GetBlocked(),
			VerifiedAt:           nullableTime(user.GetVerifiedAt()),
			Status:               user.GetStatus(),
			ScheduledForDeletion: user.GetScheduledForDeletion(),
			ScheduledDeletionAt:  nullableTime(user.GetScheduledDeletionAt()),
			LastLoginAt:          nullableTime(user.GetLastLoginAt()),
			CreatedAt:            user.GetCreatedAt(),
			UpdatedAt:            user.GetUpdatedAt(),
			CompanyID:            nullableUUID(user.GetCompanyID()),
			Companyname:          nullableString(user.GetCompanyname()),
		})
	}

	return result
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

func nullableUUID(nu uuid.NullUUID) *string {
	if nu.Valid {
		s := nu.UUID.String()
		return &s
	}
	return nil
}




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

