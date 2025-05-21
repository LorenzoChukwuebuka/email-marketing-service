package mapper

import (
	"email-marketing-service/core/handler/contacts/dto"
	db "email-marketing-service/internal/db/sqlc"
	"github.com/google/uuid"
	"time"
)

func MapContactAllContactResponse(u db.GetAllContactsRow) *dto.ContactResponse {
	var (
		createdAt, updatedAt, deletedAt    *time.Time
		groupCreatedAt, groupUpdatedAt     *time.Time
		description                        *string
		contactGroupID, groupID            *uuid.UUID
		groupCreatorID, userContactGroupID *uuid.UUID
		contactID2                         *uuid.UUID
	)

	// Time values from nullable DB columns
	if u.UcgCreatedAt.Valid {
		createdAt = &u.UcgCreatedAt.Time
	}
	if u.UcgUpdatedAt.Valid {
		updatedAt = &u.UcgUpdatedAt.Time
	}
	if u.UcgDeletedAt.Valid {
		deletedAt = &u.UcgDeletedAt.Time
	}
	if u.GroupCreatedAt.Valid {
		groupCreatedAt = &u.GroupCreatedAt.Time
	}
	if u.GroupUpdatedAt.Valid {
		groupUpdatedAt = &u.GroupUpdatedAt.Time
	}

	// Optional string field
	if u.Description.Valid {
		description = &u.Description.String
	}

	// UUIDs from nullable DB columns
	if u.ContactGroupID.Valid {
		val := u.ContactGroupID.UUID
		contactGroupID = &val
	}
	if u.GroupID.Valid {
		val := u.GroupID.UUID
		groupID = &val
	}
	if u.GroupCreatorID.Valid {
		val := u.GroupCreatorID.UUID
		groupCreatorID = &val
	}
	if u.UserContactGroupID.Valid {
		val := u.UserContactGroupID.UUID
		userContactGroupID = &val
	}
	if u.ContactID_2.Valid {
		val := u.ContactID_2.UUID
		contactID2 = &val
	}

	// Optional bool field
	isSubscribed := false
	if u.IsSubscribed.Valid {
		isSubscribed = u.IsSubscribed.Bool
	}

	// Contact timestamps
	var contactCreatedAt, contactUpdatedAt *time.Time
	if u.ContactCreatedAt.Valid {
		contactCreatedAt = &u.ContactCreatedAt.Time
	}
	if u.ContactUpdatedAt.Valid {
		contactUpdatedAt = &u.ContactUpdatedAt.Time
	}

	// Create the response with core contact info
	response := &dto.ContactResponse{

		ContactID:        u.ContactID,
		CompanyID:        u.CompanyID,
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		Email:            u.Email,
		FromOrigin:       u.FromOrigin,
		IsSubscribed:     isSubscribed,
		UserID:           u.UserID,
		ContactCreatedAt: contactCreatedAt,
		ContactUpdatedAt: contactUpdatedAt,

	}

	// Only add UserContactGroup if we have relevant data
	if userContactGroupID != nil || contactID2 != nil || contactGroupID != nil {
		response.UserContactGroup = &dto.UserContactGroup{
			ID:             userContactGroupID,
			ContactID:      contactID2,
			ContactGroupID: contactGroupID,
			CreatedAt:      createdAt,
			UpdatedAt:      updatedAt,
			DeletedAt:      deletedAt,
		}
	}

	// Only add Group if we have relevant data
	if groupID != nil || u.GroupName.Valid || groupCreatorID != nil {
		var groupName string
		if u.GroupName.Valid {
			groupName = u.GroupName.String
		}

		response.Group = &dto.GroupInfo{
			GroupID:        groupID,
			GroupName:      groupName,
			Description:    description,
			GroupCreatorID: groupCreatorID,
			GroupCreatedAt: groupCreatedAt,
			GroupUpdatedAt: groupUpdatedAt,
		}
	}

	return response
}
