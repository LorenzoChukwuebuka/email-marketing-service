package repository

import (
	"email-marketing-service/api/v1/model"
	"errors"
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

func (r *ContactRepository) GetAllContacts(userId string, params PaginationParams, searchQuery string) (PaginatedResult, error) {
	var contacts []model.Contact

	query := r.DB.Model(&model.Contact{}).Where("user_id = ?", userId).Order("created_at DESC").Preload("Groups")

	if searchQuery != "" {
		query = query.Where("first_name LIKE ? OR last_name LIKE ? OR email LIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

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

	return paginatedResult, nil
}

func mapContactToResponse(contact model.Contact) model.ContactResponse {
	response := model.ContactResponse{
		ID:           contact.ID,
		UUID:         contact.UUID,
		FirstName:    contact.FirstName,
		LastName:     contact.LastName,
		Email:        contact.Email,
		From:         contact.From,
		UserId:       contact.UserId,
		IsSubscribed: contact.IsSubscribed,
		CreatedAt:    contact.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    contact.UpdatedAt.Format(time.RFC3339),
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

func mapGroupToResponse(group model.ContactGroup) model.ContactGroupResponse {
	groupResponse := model.ContactGroupResponse{
		ID:          group.ID,
		UUID:        group.UUID,
		GroupName:   group.GroupName,
		UserId:      group.UserId,
		Description: group.Description,
		CreatedAt:   group.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   group.UpdatedAt.Format(time.RFC3339),
	}

	if group.DeletedAt.Valid {
		formatted := group.DeletedAt.Time.Format(time.RFC3339)
		groupResponse.DeletedAt = &formatted
	}

	for _, contact := range group.Contacts {
		contactResponse := mapContactToResponse(contact)
		groupResponse.Contacts = append(groupResponse.Contacts, contactResponse)
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

	if err := r.DB.Where("uuid = ? AND user_id = ?", d.UUID, d.UserId).First(&existingContact).Error; err != nil {
		return fmt.Errorf("failed to find contact for update: %w", err)
	}

	if d.FirstName != "" {
		existingContact.FirstName = d.FirstName
	}
	if d.LastName != "" {
		existingContact.LastName = d.LastName
	}
	if d.Email != "" {
		existingContact.Email = d.Email
	}
	if d.From != "" {
		existingContact.From = d.From
	}

	existingContact.IsSubscribed = d.IsSubscribed
	existingContact.UpdatedAt = time.Now().UTC()

	if err := r.DB.Save(&existingContact).Error; err != nil {
		return fmt.Errorf("failed to update contact: %w", err)
	}

	return nil
}

func (r *ContactRepository) UpdateSubscriptionStatus(email string) error {
	var existingContact model.Contact

	if err := r.DB.Where("email = ?", email).First(&existingContact).Error; err != nil {
		return fmt.Errorf("failed to find contact for update: %w", err)
	}

	existingContact.IsSubscribed = false

	if err := r.DB.Save(&existingContact).Error; err != nil {
		return fmt.Errorf("failed to update contact: %w", err)
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

func mapToContactGroupResponse(group model.ContactGroup) model.ContactGroupResponse {
	response := model.ContactGroupResponse{
		ID:          group.ID,
		UUID:        group.UUID,
		GroupName:   group.GroupName,
		UserId:      group.UserId,
		Description: group.Description,
		CreatedAt:   group.CreatedAt.Format(time.RFC3339), // Format time as needed
		UpdatedAt:   group.UpdatedAt.Format(time.RFC3339),
		Contacts:    mapContactsToResponse(group.Contacts),
	}

	if group.DeletedAt.Valid {
		formatted := group.DeletedAt.Time.Format(time.RFC3339)
		response.DeletedAt = &formatted
	}
	return response
}

func mapContactsToResponse(contacts []model.Contact) []model.ContactResponse {
	var contactResponses []model.ContactResponse
	for _, contact := range contacts {

		contactResponses = append(contactResponses, model.ContactResponse{
			ID:        contact.ID,
			UUID:      contact.UUID,
			FirstName: contact.FirstName,
			LastName:  contact.LastName,
			Email:     contact.Email,
			From:      contact.From,
			UserId:    contact.UserId,
			CreatedAt: contact.CreatedAt.Format(time.RFC3339),
			UpdatedAt: contact.UpdatedAt.Format(time.RFC3339),
			DeletedAt: func() *string {
				var htime string
				if contact.DeletedAt.Valid {
					htime = contact.DeletedAt.Time.Format(time.RFC3339)
				}
				return &htime
			}(),
		})

	}
	return contactResponses
}

func (r *ContactRepository) GetAllGroups(userId string, params PaginationParams, searchQuery string) (PaginatedResult, error) {
	var groups []model.ContactGroup

	query := r.DB.
		Preload("Contacts", func(db *gorm.DB) *gorm.DB {
			return db.Select("contacts.*, user_contact_groups.contact_group_id").
				Joins("LEFT JOIN user_contact_groups ON user_contact_groups.contact_id = contacts.id").
				Where("contacts.deleted_at IS NULL").
				Where("user_contact_groups.deleted_at IS NULL")
		}).
		Where("contact_groups.user_id = ?", userId).
		Where("contact_groups.deleted_at IS NULL")

	if searchQuery != "" {
		query = query.Where("contact_groups.group_name LIKE ?", "%"+searchQuery+"%")
	}

	query = query.Order("contact_groups.created_at DESC")

	paginatedResult, err := Paginate(query, params, &groups)
	if err != nil {
		return PaginatedResult{}, fmt.Errorf("failed to paginate contacts: %w", err)
	}

	var response []model.ContactGroupResponse

	for _, group := range groups {
		response = append(response, mapGroupToResponse(group))
	}

	paginatedResult.Data = response

	return paginatedResult, nil
}

func (r *ContactRepository) GetASingleGroup(userId string, groupId string) (*model.ContactGroup, error) {
	var groups model.ContactGroup
	err := r.DB.Where("user_id = ? AND uuid = ?", userId, groupId).First(&groups).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all groups: %w", err)
	}
	return &groups, nil
}

func (r *ContactRepository) GetASingleGroupWithContacts(userId string, groupId string) (*model.ContactGroupResponse, error) {
	var group model.ContactGroup

	err := r.DB.Preload("Contacts", func(db *gorm.DB) *gorm.DB {
		return db.Joins("JOIN user_contact_groups ON user_contact_groups.contact_id = contacts.id").
			Where("user_contact_groups.user_id = ?", userId).
			Where("contacts.deleted_at IS NULL").
			Where("user_contact_groups.deleted_at IS NULL")
	}).
		Where("contact_groups.uuid = ?", groupId).
		Where("contact_groups.deleted_at IS NULL").
		First(&group).Error

	if err != nil {
		return nil, err
	}

	response := mapToContactGroupResponse(group)

	return &response, nil
}

func (r *ContactRepository) GetGroupById(userId string, groupId int) (*model.ContactGroupResponse, error) {
	var group model.ContactGroup
	err := r.DB.Preload("Contacts", func(db *gorm.DB) *gorm.DB {
		return db.Joins("JOIN user_contact_groups ON user_contact_groups.contact_id = contacts.id").
			Where("user_contact_groups.user_id = ?", userId).
			Where("contacts.deleted_at IS NULL").
			Where("user_contact_groups.deleted_at IS NULL")
	}).
		Where("contact_groups.id = ?", groupId).
		Where("contact_groups.deleted_at IS NULL").
		First(&group).Error

	if err != nil {
		return nil, err
	}

	response := mapToContactGroupResponse(group)

	return &response, nil
}

func (r *ContactRepository) DeleteContactGroup(userId string, groupId int) error {

	tx := r.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	// Check if the group exists in user contact groups
	var userContactGroup model.UserContactGroup
	exists := tx.Where("contact_group_id = ?", groupId).First(&userContactGroup).Error

	if exists != nil && !errors.Is(exists, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return fmt.Errorf("failed to check user contact groups: %w", exists)
	}

	// If the group exists, delete associated user contact groups
	if exists == nil {
		result := tx.Where("contact_group_id = ?", groupId).Delete(&model.UserContactGroup{})
		if result.Error != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete user contact groups: %w", result.Error)
		}
	}

	// Delete the contact group
	result := tx.Where("user_id = ? AND id = ?", userId, groupId).Delete(&model.ContactGroup{})
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete contact group: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("group not found or already deleted")
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *ContactRepository) RemoveContactFromGroup(groupId int, userId string, contactId int) error {
	result := r.DB.Where("contact_group_id = ? AND user_id = ? AND contact_id = ?", groupId, userId, contactId).
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

	if err := r.DB.Where("uuid = ? AND user_id = ?", d.UUID, d.UserId).First(&existingContactGroup).Error; err != nil {
		return fmt.Errorf("failed to find group: %w", err)
	}

	existingContactGroup.GroupName = d.GroupName
	existingContactGroup.Description = d.Description

	if err := r.DB.Save(&existingContactGroup).Error; err != nil {
		return fmt.Errorf("failed to update group: %w", err)
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

func (r *ContactRepository) GetContactCount(userId string) (map[string]int64, error) {
	contactCounts := make(map[string]int64)
	var totalCount int64
	var recentCount int64

	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	if err := r.DB.Model(&model.Contact{}).Where("user_id = ?", userId).Count(&totalCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count total contacts: %w", err)
	}

	if err := r.DB.Model(&model.Contact{}).Where("user_id = ? AND created_at >= ?", userId, thirtyDaysAgo).Count(&recentCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count recent contacts: %w", err)
	}

	contactCounts["total"] = totalCount
	contactCounts["recent"] = recentCount

	return contactCounts, nil
}

func (r *ContactRepository) GetContactSubscriptionStatusForDashboard(userId string) (map[string]int64, error) {
	result := make(map[string]int64)

	// 1. Total counts of unsubscribed contacts
	var unsubscribedCount int64
	if err := r.DB.Model(&model.Contact{}).Where("user_id = ? AND is_subscribed = ?", userId, false).Count(&unsubscribedCount).Error; err != nil {
		return nil, err
	}
	result["unsubscribed"] = unsubscribedCount

	// 2. Total counts of contacts
	var totalCount int64
	if err := r.DB.Model(&model.Contact{}).Where("user_id = ?", userId).Count(&totalCount).Error; err != nil {
		return nil, err
	}
	result["total"] = totalCount

	// 3. New contacts (contacts less than 10 days old)
	var newContactsCount int64
	tenDaysAgo := time.Now().AddDate(0, 0, -10)
	if err := r.DB.Model(&model.Contact{}).Where("user_id = ? AND created_at >= ?", userId, tenDaysAgo).Count(&newContactsCount).Error; err != nil {
		return nil, err
	}
	result["new"] = newContactsCount

	// 4. Engaged subscribers (contacts who opened, clicked, or converted in any campaign)
	var engagedCount int64
	if err := r.DB.
		Table("email_campaign_results").
		Select("COUNT(DISTINCT contacts.id)").
		Joins("JOIN contacts ON contacts.email = email_campaign_results.recipient_email").
		Where("contacts.user_id = ?", userId).
		Where("email_campaign_results.opened_at IS NOT NULL OR email_campaign_results.clicked_at IS NOT NULL OR email_campaign_results.conversion_at IS NOT NULL").
		Count(&engagedCount).Error; err != nil {
		return nil, err
	}
	result["engaged"] = engagedCount

	return result, nil
}
