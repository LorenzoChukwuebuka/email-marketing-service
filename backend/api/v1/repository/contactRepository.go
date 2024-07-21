package repository

import (
	"email-marketing-service/api/v1/model"
	"fmt"

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

func (r *ContactRepository) createcontactResponse() {

}

func (r *ContactRepository) CreateContact(d *model.Contact) error {
	if err := r.DB.Create(&d).Error; err != nil {
		return fmt.Errorf("failed to insert contact: %w", err)
	}

	return nil
}

func (r *ContactRepository) CheckIfEmailExists(d *model.Contact) (bool, error) {
	result := r.DB.Where("email = ?", d.Email).First(&d)
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

func (r *ContactRepository) GetAllContacts(userId string) {

}
