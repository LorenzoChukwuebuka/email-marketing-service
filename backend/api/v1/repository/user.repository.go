package repository

import (
	"email-marketing-service/api/v1/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

const (
	DeletionGracePeriod   = 30 * 24 * time.Hour // 30 days
	StatusActive          = "active"
	StatusPendingDeletion = "pending_deletion"
	StatusDeleted         = "deleted"
)

func (r *UserRepository) createUserResponse(user model.User) model.UserResponse {

	htime := user.VerifiedAt.String()

	response := model.UserResponse{
		ID:          user.ID,
		UUID:        user.UUID,
		FullName:    user.FullName,
		Email:       user.Email,
		Company:     user.Company,
		PhoneNumber: user.PhoneNumber,
		Password:    user.Password, // Note: Make sure you have a good reason to include the password in the response
		Verified:    user.Verified,
		CreatedAt:   user.CreatedAt.String(),
		VerifiedAt:  &htime,
		UpdatedAt:   user.UpdatedAt.String(),
	}

	if user.DeletedAt.Valid {
		formatted := user.DeletedAt.Time.Format(time.RFC3339)
		response.DeletedAt = &formatted
	}

	return response
}

func (r *UserRepository) CreateUser(d *model.User) (*model.User, error) {
	if err := r.DB.Create(&d).Error; err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}
	return d, nil
}

func (r *UserRepository) CheckIfEmailAlreadyExists(d *model.User) (bool, error) {
	result := r.DB.Where("email = ?", d.Email).First(&d)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}

func (r *UserRepository) VerifyUserAccount(d *model.User) (uint, error) {
	var user model.User

	// Fetch the User record from the database
	if err := r.DB.Where("uuid = ?", d.UUID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, err
		}
		return 0, nil
	}

	user.Verified = d.Verified
	user.VerifiedAt = d.VerifiedAt

	htime := time.Now().UTC()
	user.UpdatedAt = htime

	if err := r.DB.Save(&user).Error; err != nil {
		fmt.Printf("Error updating user: %v\n", err)
		return 0, err
	}

	return user.ID, nil
}

func (r *UserRepository) Login(d *model.User) (model.UserResponse, error) {
	var user model.User

	// Fetch the user record from the database based on the provided email
	if err := r.DB.Where("email = ?", d.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.UserResponse{}, fmt.Errorf("user not found")
		}
		return model.UserResponse{}, fmt.Errorf("error querying database: %w", err)
	}

	userResponse := r.createUserResponse(user)

	return userResponse, nil
}

func (r *UserRepository) FindUserById(d *model.User) (model.UserResponse, error) {

	var user model.User
	if err := r.DB.Where("uuid = ?", d.UUID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.UserResponse{}, nil
		}

		return model.UserResponse{}, err
	}

	userResponse := r.createUserResponse(user)

	return userResponse, nil
}

func (r *UserRepository) FindUserByEmail(d *model.User) (model.UserResponse, error) {
	var user model.User
	if err := r.DB.Where("email = ?", d.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.UserResponse{}, nil
		}

		return model.UserResponse{}, err
	}

	userResponse := r.createUserResponse(user)

	return userResponse, nil
}

func (r *UserRepository) ResetPassword(d *model.User) error {

	var user model.User

	// Fetch the User record from the database
	if err := r.DB.Where("uuid = ?", d.UUID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return err
		}
		return nil
	}

	//update the password

	user.Password = d.Password

	if err := r.DB.Save(&user).Error; err != nil {
		fmt.Printf("Error updating user: %v\n", err)
		return err
	}

	return nil
}
func (r *UserRepository) FindAllUsers() ([]model.UserResponse, error) {

	return nil, nil
}

func (r *UserRepository) ChangeUserPassword(d *model.User) error {
	var user model.User
	if err := r.DB.Where("uuid = ?", d.UUID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}

	user.Password = d.Password

	if err := r.DB.Save(&user).Error; err != nil {
		fmt.Printf("Error updating user password: %v\n", err)
		return err
	}

	return nil
}

func (r *UserRepository) UpdateUserRecords(d *model.User) error {
	var user model.User

	// Fetch the user record from the database based on the UUID
	if err := r.DB.Where("uuid = ?", d.UUID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("error querying database: %w", err)
	}

	// Update the fields of the user record
	user.FullName = d.FullName
	user.Email = d.Email
	user.PhoneNumber = d.PhoneNumber
	user.Company = d.Company

	htime := time.Now().UTC()
	user.UpdatedAt = htime

	// Save the updated user record to the database
	if err := r.DB.Save(&user).Error; err != nil {
		fmt.Printf("Error updating user: %v\n", err)
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *UserRepository) CreateTempEmail(d *model.UserTempEmail) error {
	if err := r.DB.Create(&d).Error; err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil

}

func (r *UserRepository) MarkUserForDeletion(userId string) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user model.User
	if err := tx.Where("uuid = ?", userId).First(&user).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Check if user is already scheduled for deletion
	if user.ScheduledForDeletion {
		tx.Rollback()
		return errors.New("user is already scheduled for deletion")
	}

	scheduledDeletion := time.Now().Add(DeletionGracePeriod)
	updates := map[string]interface{}{
		"scheduled_for_deletion": true,
		"scheduled_deletion_at":  scheduledDeletion,
		"status":                 StatusPendingDeletion,
	}

	if err := tx.Model(&user).Updates(updates).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to mark user for deletion: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// CancelUserDeletion allows users to cancel their scheduled deletion
func (r *UserRepository) CancelUserDeletion(userId string) error {
	result := r.DB.Model(&model.User{}).
		Where("uuid = ? AND scheduled_for_deletion = ?", userId, true).
		Updates(map[string]interface{}{
			"scheduled_for_deletion": false,
			"scheduled_deletion_at":  nil,
			"status":                 StatusActive,
		})

	if result.Error != nil {
		return fmt.Errorf("failed to cancel deletion: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found or not scheduled for deletion")
	}

	return nil
}

// PermanentlyDeleteUser handles the actual deletion of a user
func (r *UserRepository) PermanentlyDeleteUser(userId string) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Find and verify user status
	var user model.User
	if err := tx.Where("uuid = ? AND scheduled_for_deletion = ? AND scheduled_deletion_at <= ?",
		userId, true, time.Now()).First(&user).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found or not scheduled for deletion")
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Get user statistics for archiving
	userStats, err := r.GetUserStats(userId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get user stats for archiving: %w", err)
	}

	// Create archive record
	archiveData := model.UserArchive{
		UserID:           user.UUID,
		Email:            user.Email,
		FullName:         user.FullName,
		Company:          user.Company,
		DeletedAt:        time.Now(),
		AccountCreatedAt: user.CreatedAt,
		VerifiedAt:       user.VerifiedAt,
		LastLoginAt:      user.LastLoginAt,
		DeletionReason:   "user_requested",
		AccountStats: model.AccountStats{
			TotalContacts:      userStats.TotalContacts,
			TotalCampaigns:     userStats.TotalCampaigns,
			TotalTemplates:     userStats.TotalTemplates,
			TotalCampaignsSent: userStats.TotalCampaignsSent,
			TotalGroups:        userStats.TotalGroups,
		},
	}

	if err := tx.Create(&archiveData).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create archive record: %w", err)
	}

	// Delete associated records
	deletions := []struct {
		model interface{}
		name  string
	}{
		{&model.Contact{}, "contacts"},
		{&model.ContactGroup{}, "contact groups"},
		{&model.Template{}, "templates"},
		{&model.Campaign{}, "campaigns"},
	}

	for _, deletion := range deletions {
		if err := tx.Where("user_id = ?", userId).Delete(deletion.model).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete %s: %w", deletion.name, err)
		}
	}

	// Delete the user
	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

type UserStats struct {
	TotalContacts      int64 `json:"total_contacts"`
	TotalCampaigns     int64 `json:"total_campaigns"`
	TotalTemplates     int64 `json:"total_templates"`
	TotalCampaignsSent int64 `json:"total_campaigns_sent"`
	//	TotalSubscriptions int64 `json:"total_subscriptions"`
	TotalGroups int64 `json:"total_groups"`
}

func (r *UserRepository) GetUserStats(userID string) (UserStats, error) {
	var stats UserStats

	// Count total contacts for the user
	if err := r.DB.Model(&model.Contact{}).Where("user_id = ?", userID).Count(&stats.TotalContacts).Error; err != nil {
		return stats, fmt.Errorf("failed to count contacts: %w", err)
	}

	// Count total campaigns for the user
	if err := r.DB.Model(&model.Campaign{}).Where("user_id = ?", userID).Count(&stats.TotalCampaigns).Error; err != nil {
		return stats, fmt.Errorf("failed to count campaigns: %w", err)
	}

	if err := r.DB.Model(&model.Template{}).Where("user_id = ?", userID).Count(&stats.TotalTemplates).Error; err != nil {
		return stats, fmt.Errorf("failed to count templates: %w", err)
	}

	if err := r.DB.Model(&model.Campaign{}).Where("user_id = ? AND status = ?", userID, model.Sent).Count(&stats.TotalCampaignsSent).Error; err != nil {
		return stats, fmt.Errorf("failed to count campaigns: %w", err)
	}

	if err := r.DB.Model(&model.ContactGroup{}).Where("user_id = ? ", userID).Count(&stats.TotalGroups).Error; err != nil {
		return stats, fmt.Errorf("failed to count groups: %w", err)
	}

	return stats, nil
}
