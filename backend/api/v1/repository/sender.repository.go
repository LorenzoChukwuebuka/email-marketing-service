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
