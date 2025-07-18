package mapper

import (
	"database/sql"
	db "email-marketing-service/internal/db/sqlc"
	"github.com/google/uuid"
	"time"
)

//
// GetAllUsersAdapter
//

type GetAllUsersAdapter db.GetAllUsersRow

func (u GetAllUsersAdapter) GetID() uuid.UUID                     { return u.ID }
func (u GetAllUsersAdapter) GetFullname() string                  { return u.Fullname }
func (u GetAllUsersAdapter) GetEmail() string                     { return u.Email }
func (u GetAllUsersAdapter) GetPhonenumber() sql.NullString       { return u.Phonenumber }
func (u GetAllUsersAdapter) GetPicture() sql.NullString           { return u.Picture }
func (u GetAllUsersAdapter) GetVerified() bool                    { return u.Verified }
func (u GetAllUsersAdapter) GetBlocked() bool                     { return u.Blocked }
func (u GetAllUsersAdapter) GetVerifiedAt() sql.NullTime          { return u.VerifiedAt }
func (u GetAllUsersAdapter) GetStatus() string                    { return u.Status }
func (u GetAllUsersAdapter) GetScheduledForDeletion() bool        { return u.ScheduledForDeletion }
func (u GetAllUsersAdapter) GetScheduledDeletionAt() sql.NullTime { return u.ScheduledDeletionAt }
func (u GetAllUsersAdapter) GetLastLoginAt() sql.NullTime         { return u.LastLoginAt }
func (u GetAllUsersAdapter) GetCreatedAt() time.Time              { return u.CreatedAt }
func (u GetAllUsersAdapter) GetUpdatedAt() time.Time              { return u.UpdatedAt }
func (u GetAllUsersAdapter) GetCompanyID() uuid.NullUUID          { return u.CompanyID }
func (u GetAllUsersAdapter) GetCompanyname() sql.NullString       { return u.Companyname }

//
// GetVerifiedUsersAdapter
//

type GetVerifiedUsersAdapter db.GetVerifiedUsersRow

func (u GetVerifiedUsersAdapter) GetID() uuid.UUID                     { return u.ID }
func (u GetVerifiedUsersAdapter) GetFullname() string                  { return u.Fullname }
func (u GetVerifiedUsersAdapter) GetEmail() string                     { return u.Email }
func (u GetVerifiedUsersAdapter) GetPhonenumber() sql.NullString       { return u.Phonenumber }
func (u GetVerifiedUsersAdapter) GetPicture() sql.NullString           { return u.Picture }
func (u GetVerifiedUsersAdapter) GetVerified() bool                    { return u.Verified }
func (u GetVerifiedUsersAdapter) GetBlocked() bool                     { return u.Blocked }
func (u GetVerifiedUsersAdapter) GetVerifiedAt() sql.NullTime          { return u.VerifiedAt }
func (u GetVerifiedUsersAdapter) GetStatus() string                    { return u.Status }
func (u GetVerifiedUsersAdapter) GetScheduledForDeletion() bool        { return u.ScheduledForDeletion }
func (u GetVerifiedUsersAdapter) GetScheduledDeletionAt() sql.NullTime { return u.ScheduledDeletionAt }
func (u GetVerifiedUsersAdapter) GetLastLoginAt() sql.NullTime         { return u.LastLoginAt }
func (u GetVerifiedUsersAdapter) GetCreatedAt() time.Time              { return u.CreatedAt }
func (u GetVerifiedUsersAdapter) GetUpdatedAt() time.Time              { return u.UpdatedAt }
func (u GetVerifiedUsersAdapter) GetCompanyID() uuid.NullUUID          { return u.CompanyID }
func (u GetVerifiedUsersAdapter) GetCompanyname() sql.NullString       { return u.Companyname }

//
// GetUnVerifiedUsersAdapter
//

type GetUnVerifiedUsersAdapter db.GetUnVerifiedUsersRow

func (u GetUnVerifiedUsersAdapter) GetID() uuid.UUID               { return u.ID }
func (u GetUnVerifiedUsersAdapter) GetFullname() string            { return u.Fullname }
func (u GetUnVerifiedUsersAdapter) GetEmail() string               { return u.Email }
func (u GetUnVerifiedUsersAdapter) GetPhonenumber() sql.NullString { return u.Phonenumber }
func (u GetUnVerifiedUsersAdapter) GetPicture() sql.NullString     { return u.Picture }
func (u GetUnVerifiedUsersAdapter) GetVerified() bool              { return u.Verified }
func (u GetUnVerifiedUsersAdapter) GetBlocked() bool               { return u.Blocked }
func (u GetUnVerifiedUsersAdapter) GetVerifiedAt() sql.NullTime    { return u.VerifiedAt }
func (u GetUnVerifiedUsersAdapter) GetStatus() string              { return u.Status }
func (u GetUnVerifiedUsersAdapter) GetScheduledForDeletion() bool  { return u.ScheduledForDeletion }
func (u GetUnVerifiedUsersAdapter) GetScheduledDeletionAt() sql.NullTime {
	return u.ScheduledDeletionAt
}
func (u GetUnVerifiedUsersAdapter) GetLastLoginAt() sql.NullTime   { return u.LastLoginAt }
func (u GetUnVerifiedUsersAdapter) GetCreatedAt() time.Time        { return u.CreatedAt }
func (u GetUnVerifiedUsersAdapter) GetUpdatedAt() time.Time        { return u.UpdatedAt }
func (u GetUnVerifiedUsersAdapter) GetCompanyID() uuid.NullUUID    { return u.CompanyID }
func (u GetUnVerifiedUsersAdapter) GetCompanyname() sql.NullString { return u.Companyname }
