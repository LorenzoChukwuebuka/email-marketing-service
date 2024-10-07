package repository

import (
	"email-marketing-service/api/v1/model"
	"fmt"
	"gorm.io/gorm"
)

type SenderRepository struct {
	DB *gorm.DB
}

func NewSenderRepository(db *gorm.DB) *SenderRepository {
	return &SenderRepository{
		DB: db,
	}
}

func (r *SenderRepository) CreateSender(d *model.Sender) error {
	if err := r.DB.Create(&d).Error; err != nil {
		return fmt.Errorf("failed to insert domain: %w", err)
	}
	return nil
}

func (r *SenderRepository) CheckIfSenderExists(d *model.Sender) (bool, error) {
	result := r.DB.Where("email = ? AND name = ? AND user_id =?", d.Email, d.Name, d.UserID).First(&d)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}

func (r *SenderRepository) GetAllSenders(userId string, searchQuery string, params PaginationParams) (PaginatedResult, error) {
	var senders []model.Sender
	query := r.DB.Model(&model.Sender{}).Where("user_id = ?", userId)

	if searchQuery != "" {
		query = query.Where("email ILIKE ? OR name ILIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	query.Order("created_at DESC")

	paginatedResult, err := Paginate(query, params, &senders)
	if err != nil {
		return PaginatedResult{}, fmt.Errorf("failed to paginate contacts: %w", err)
	}

	var response []model.SenderResponse

	for _, sender := range senders {
		senderResponse := model.SenderResponse{
			UUID:      sender.UUID,
			UserID:    sender.UserID,
			Name:      sender.Name,
			Email:     sender.Email,
			Verified:  sender.Verified,
			IsSigned:  sender.IsSigned,
			CreatedAt: sender.CreatedAt.Format("2006-01-02 15:04:05"), // Format time to string
			UpdatedAt: sender.UpdatedAt.Format("2006-01-02 15:04:05"),
			DeletedAt: func() *string {
				if sender.DeletedAt.Valid {
					formatted := sender.DeletedAt.Time.Format("2006-01-02 15:04:05")
					return &formatted
				}
				return nil
			}(),
		}
		response = append(response, senderResponse)
	}

	paginatedResult.Data = response

	return paginatedResult, nil
}

func (r *SenderRepository) DeleteSender(uuid string, userId string) error {
	var sender model.Sender

	// Query the sender by uuid and userId to ensure that the sender belongs to the correct user.
	if err := r.DB.Where("uuid = ? AND user_id = ?", uuid, userId).First(&sender).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return err
		} else {
			fmt.Printf("Error querying database: %v\n", err)
			return err
		}
	}

	// Delete the sender record
	if err := r.DB.Delete(&sender).Error; err != nil {
		fmt.Printf("Error deleting sender: %v\n", err)
		return err
	}

	return nil
}

func (r *SenderRepository) FindSenderByID(id string, userId string) (*model.Sender, error) {
	var sender model.Sender
	// Query the sender by email and user ID
	result := r.DB.Where("uuid = ? AND user_id = ?", id, userId).First(&sender)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("sender not found")
		}
		return nil, result.Error
	}
	return &sender, nil
}

func (r *SenderRepository) UpdateSender(d *model.Sender) error {
	// Perform the update using GORM's Save method
	result := r.DB.Model(&model.Sender{}).Where("uuid = ? AND user_id = ?", d.UUID, d.UserID).Updates(map[string]interface{}{
		"name":      d.Name,
		"email":     d.Email,
		"verified":  d.Verified,
		"is_signed": d.IsSigned,
	})

	// Check for any errors during the update
	if result.Error != nil {
		return fmt.Errorf("failed to update sender: %w", result.Error)
	}

	// If no rows were affected, return an error indicating no record was found
	if result.RowsAffected == 0 {
		return fmt.Errorf("no sender record found to update")
	}

	return nil
}
