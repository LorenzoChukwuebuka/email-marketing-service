package repository

import (
	"email-marketing-service/api/v1/model"
	"fmt"
	"gorm.io/gorm"
	"time"
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
	result := r.DB.Where("email = ? AND user_id =?", d.Email, d.UserId).First(&d)
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

func (r *ContactRepository) GetASingleContact(contactId string, userId string) (*model.ContactResponse, error) {
	var contact model.Contact

	result := r.DB.First(&contact, "uuid = ? AND user_id =?", contactId, userId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // No contact found with the given ID
		}
		return nil, fmt.Errorf("failed to retrieve contact: %w", result.Error)
	}

	//htime := contact.UpdatedAt

	contactResponse := &model.ContactResponse{
		ID:        contact.ID,
		Email:     contact.Email,
		UserId:    contact.UserId,
		CreatedAt: contact.CreatedAt.Format(time.RFC3339),
	}

	return contactResponse, nil
}

func (r *ContactRepository) GetAllContacts(userId string, params PaginationParams) (PaginatedResult, error) {
	var contacts []model.Contact

	query := r.DB.Model(&model.Contact{}).Where("user_id = ?", userId).Preload("Groups")

	paginatedResult, err := Paginate(query, params, &contacts)
	if err != nil {
		return PaginatedResult{}, fmt.Errorf("failed to paginate contacts: %w", err)
	}

	// Map the contacts to ContactResponse
	var contactResponses []model.ContactResponse
	for _, contact := range contacts {
		response := mapContactToResponse(contact)
		contactResponses = append(contactResponses, response)
	}

	paginatedResult.Data = contactResponses

	// Log the final result
	fmt.Printf("Paginated result: %+v\n", paginatedResult)

	return paginatedResult, nil
}

func mapContactToResponse(contact model.Contact) model.ContactResponse {
	response := model.ContactResponse{
		ID:        contact.ID,
		UUID:      contact.UUID,
		FirstName: contact.FirstName,
		LastName:  contact.LastName,
		Email:     contact.Email,
		From:      contact.From,
		UserId:    contact.UserId,
		CreatedAt: contact.CreatedAt.Format(time.RFC3339),
		UpdatedAt: contact.UpdatedAt.Format(time.RFC3339),
	}

	if contact.DeletedAt.Valid {
		formatted := contact.DeletedAt.Time.Format(time.RFC3339)
		response.DeletedAt = &formatted
	}

	for _, group := range contact.Groups {
		groupResponse := mapGroupToResponse(group)
		response.Groups = append(response.Groups, groupResponse)
	}

	return response
}

func mapGroupToResponse(group model.ContactGroup) model.ContactGroup {
	groupResponse := model.ContactGroup{
		UUID:        group.UUID,
		GroupName:   group.GroupName,
		Description: group.Description,
	}

	return groupResponse
}

func (r *ContactRepository) DeleteContact(userId string, contactId string) error {

	var existingContact model.Contact
	if err := r.DB.Where("uuid = ? AND user_id =?", contactId, userId).First(&existingContact).Error; err != nil {
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
	if err := r.DB.Where("uuid = ? AND user_id =?", d.UUID, d.UserId).First(&existingContact).Error; err != nil {
		return fmt.Errorf("failed to find plan for deletion: %w", err)
	}

	existingContact.FirstName = d.FirstName
	existingContact.LastName = d.LastName
	existingContact.Email = d.Email
	existingContact.From = d.From
	htime := time.Now().UTC()
	existingContact.UpdatedAt = htime

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
	var groups model.ContactGroup
	err := r.DB.Where("user_id = ?", userId).Find(&groups).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all groups: %w", err)
	}
	return &groups, nil
}

func (r *ContactRepository) GetASingleGroupWithContacts(userId string, groupId string) (*model.ContactGroup, error) {
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

func (r *ContactRepository) UpdateGroup(d *model.ContactGroup) error {

	var existingContactGroup model.ContactGroup

	if err := r.DB.Where("uuid = ? AND user_id", d.UUID, d.UserId).First(&existingContactGroup).Error; err != nil {
		return fmt.Errorf("failed to find plan for deletion: %w", err)
	}

	existingContactGroup.GroupName = d.GroupName
	existingContactGroup.Description = d.Description

	if err := r.DB.Save(&existingContactGroup).Error; err != nil {
		return fmt.Errorf("failed to update plan: %w", err)
	}

	return nil
}

func (r *ContactRepository) CheckIfGroupNameExists(d *model.ContactGroup) (bool, error) {
	result := r.DB.Where("group_name = ? AND user_id = ?", d.GroupName, d.UserId).First(&d)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil

}
