package mapper

import (
	"email-marketing-service/core/handler/contacts/dto"
	db "email-marketing-service/internal/db/sqlc"
	"github.com/google/uuid"
	"time"
)

// MapGroupsWithContacts transforms the raw database rows into grouped contact responses
func MapGroupsWithContacts(rows []db.GetGroupsWithContactsRow) []*dto.GroupwithContactResponse {
	// Use a map to track unique groups by their ID
	groupMap := make(map[uuid.UUID]*dto.GroupwithContactResponse)

	// Process each row
	for _, row := range rows {
		// Check if we've already created an entry for this group
		group, exists := groupMap[row.GroupID]

		if !exists {
			// Create new group entry if this is the first time we see this group ID
			var description *string
			if row.Description.Valid {
				description = &row.Description.String
			}

			var groupCreatedAt *time.Time
			if row.GroupCreatedAt.Valid {
				groupCreatedAt = &row.GroupCreatedAt.Time
			}

			// Initialize new group with empty contacts array
			group = &dto.GroupwithContactResponse{
				GroupID:        row.GroupID,
				GroupName:      row.GroupName,
				Description:    description,
				GroupCreatedAt: groupCreatedAt,
				Contacts:       []dto.GroupContactResponse{},
			}

			// Add to our map
			groupMap[row.GroupID] = group
		}

		// Only add contact if contact ID is valid (contact exists)
		if row.ContactID.Valid {
			var firstName, lastName, email, origin *string
			var isSubscribed *bool
			var contactCreatedAt *time.Time

			if row.ContactFirstName.Valid {
				firstName = &row.ContactFirstName.String
			}

			if row.ContactLastName.Valid {
				lastName = &row.ContactLastName.String
			}

			if row.ContactEmail.Valid {
				email = &row.ContactEmail.String
			}

			if row.ContactFromOrigin.Valid {
				origin = &row.ContactFromOrigin.String
			}

			if row.ContactIsSubscribed.Valid {
				isSubscribed = &row.ContactIsSubscribed.Bool
			}

			if row.ContactCreatedAt.Valid {
				contactCreatedAt = &row.ContactCreatedAt.Time
			}

			contact := dto.GroupContactResponse{
				ContactID:           row.ContactID,
				ContactFirstName:    firstName,
				ContactLastName:     lastName,
				ContactEmail:        email,
				ContactFromOrigin:   origin,
				ContactIsSubscribed: isSubscribed,
				ContactCreatedAt:    contactCreatedAt,
			}

			// Add contact to the appropriate group
			group.Contacts = append(group.Contacts, contact)
		}
	}

	// Convert map to slice for the response
	result := make([]*dto.GroupwithContactResponse, 0, len(groupMap))
	for _, group := range groupMap {
		result = append(result, group)
	}

	return result
}

func MapSingleGroupwithContacts(rows []db.GetSingleGroupWithContactsRow) *dto.GroupwithContactResponse {
	if len(rows) == 0 {
		return nil
	}

	group := &dto.GroupwithContactResponse{
		GroupID:        rows[0].GroupID,
		GroupName:      rows[0].GroupName,
		Description:    &rows[0].Description.String,
		GroupCreatedAt: &rows[0].GroupCreatedAt.Time,
	}

	for _, row := range rows {
		if row.ContactID.Valid {
			contact := dto.GroupContactResponse{
				ContactID:           row.ContactID,
				ContactFirstName:    &row.ContactFirstName.String,
				ContactLastName:     &row.ContactLastName.String,
				ContactEmail:        &row.ContactEmail.String,
				ContactFromOrigin:   &row.ContactFromOrigin.String,
				ContactIsSubscribed: &row.ContactIsSubscribed.Bool,
			}
			group.Contacts = append(group.Contacts, contact)
		}
	}

	return group
}
