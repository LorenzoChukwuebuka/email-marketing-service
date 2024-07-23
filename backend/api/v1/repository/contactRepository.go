package repository

import (
	"email-marketing-service/api/v1/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type ContactRepository struct {
	DB *gorm.DB
}

func NewContactRepository(db *gorm.DB) *ContactRepository {
	return &ContactRepository{
		DB: db,
	}
}

func (r *ContactRepository) CreateContact(d *model.Contact) error {
	if err := r.DB.Create(&d).Error; err != nil {
		return fmt.Errorf("failed to insert contact: %w", err)
	}

	return nil
}

func (r *ContactRepository) CheckIfEmailExists(d *model.Contact) (bool, error) {
	result := r.DB.Where("email = ? AND user_id", d.Email, d.UserId).First(&d)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil

}

func (r *ContactRepository) BulkCreateContacts(contacts []model.Contact) error {
	result := r.DB.CreateInBatches(contacts, 100) // Insert in batches of 100
	if result.Error != nil {
		return fmt.Errorf("error bulk creating contacts: %w", result.Error)
	}
	return nil
}

func (r *ContactRepository) GetAllContacts(userId string) ([]model.ContactResponse, error) {
	var contacts []model.Contact
	var contactResponses []model.ContactResponse

	// Query the database to get all contacts for the given user ID
	result := r.DB.Where("user_id = ?", userId).Preload("Groups").Find(&contacts)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get contacts: %w", result.Error)
	}

	// Map the contacts to the ContactResponse type
	for _, contact := range contacts {
		response := model.ContactResponse{
			ID:        contact.ID,
			UUID:      contact.UUID,
			FirstName: contact.FirstName,
			LastName:  contact.LastName,
			Email:     contact.Email,
			From:      contact.From,
			UserId:    contact.UserId,
			CreatedAt: contact.CreatedAt.Format(time.RFC3339),
		}

		// Check if UpdatedAt is not nil before formatting
		if contact.UpdatedAt != nil {
			formatted := contact.UpdatedAt.Format(time.RFC3339)
			response.UpdatedAt = &formatted
		}

		// Check if DeletedAt is not nil before formatting
		if contact.DeletedAt.Valid {
			formatted := contact.DeletedAt.Time.Format(time.RFC3339)
			response.DeletedAt = &formatted
		}

		// Map the groups
		for _, group := range contact.Groups {
			groupResponse := model.ContactGroup{
				ID:          group.ID,
				UUID:        group.UUID,
				GroupName:   group.GroupName,
				Description: group.Description,
				CreatedAt:   group.CreatedAt,
			}

			// Check if UpdatedAt is not nil before assigning
			if group.UpdatedAt != nil {
				groupResponse.UpdatedAt = group.UpdatedAt
			}

			// Assign DeletedAt
			groupResponse.DeletedAt = group.DeletedAt

			response.Groups = append(response.Groups, groupResponse)
		}

		contactResponses = append(contactResponses, response)
	}

	return contactResponses, nil
}

func (r *ContactRepository) DeleteContact(userId string, contactId string) error {

	var existingContact model.Contact
	if err := r.DB.Where("uuid = ? AND user_id", contactId, userId).First(&existingContact).Error; err != nil {
		return fmt.Errorf("failed to find plan for deletion: %w", err)
	}

	// Soft delete by marking the plan as deleted
	if err := r.DB.Delete(&existingContact).Error; err != nil {
		return fmt.Errorf("failed to delete plan: %w", err)
	}

	return nil
}

func (r *ContactRepository) UpdateContact(d *model.Contact) error {
	var existingContact model.Contact
	if err := r.DB.Where("uuid = ? AND user_id", d.UUID, d.UserId).First(&existingContact).Error; err != nil {
		return fmt.Errorf("failed to find plan for deletion: %w", err)
	}

	existingContact.FirstName = d.FirstName
	existingContact.LastName = d.LastName
	existingContact.Email = d.Email
	existingContact.From = d.From
	htime := time.Now().UTC()
	existingContact.UpdatedAt = &htime

	if err := r.DB.Save(&existingContact).Error; err != nil {
		return fmt.Errorf("failed to update plan: %w", err)
	}

	return nil

}

func (r *ContactRepository) CreateGroup(d *model.ContactGroup) error {

	if err := r.DB.Create(&d).Error; err != nil {
		return fmt.Errorf("failed to insert contact group: %w", err)
	}
	return nil
}

func (r *ContactRepository) AddContactsToGroup(d *model.UserContactGroup) error {
	if err := r.DB.Create(&d).Error; err != nil {
		return fmt.Errorf("failed to add user to group: %w", err)
	}
	return nil
}

func (r *ContactRepository) GetAllGroups(userId string) ([]model.ContactGroup, error) {
	var groups []model.ContactGroup
	err := r.DB.Where("user_id = ?", userId).Find(&groups).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all groups: %w", err)
	}
	return groups, nil
}

func (r *ContactRepository) GetASingleGroup(userId string, groupId string) (*model.ContactGroup, error) {
	var group model.ContactGroup
	err := r.DB.Preload("Contacts", func(db *gorm.DB) *gorm.DB {
		return db.Select("contacts.*").
			Joins("JOIN user_contact_groups ON user_contact_groups.contact_id = contacts.id").
			Where("user_contact_groups.user_id = ?", userId)
	}).Where("contact_groups.uuid = ?", groupId).First(&group).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("group not found")
		}
		return nil, fmt.Errorf("failed to fetch group with contacts: %w", err)
	}
	return &group, nil
}

func (r *ContactRepository) DeleteContactGroup(userId string, groupId string) error {
	result := r.DB.Where("user_id = ? AND uuid = ?", userId, groupId).Delete(&model.ContactGroup{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete contact group: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("group not found or already deleted")
	}
	return nil
}

func (r *ContactRepository) RemoveContactFromGroup(groupId string, userId string, contactId uint) error {
	result := r.DB.Where("group_id = ? AND user_id = ? AND contact_id = ?", groupId, userId, contactId).
		Delete(&model.UserContactGroup{})

	if result.Error != nil {
		return fmt.Errorf("failed to remove contact from group: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no matching record found to remove contact from group")
	}

	return nil
}
