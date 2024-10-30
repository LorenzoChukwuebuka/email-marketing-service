package adminrepository

import (
	"email-marketing-service/api/v1/model"
	adminmodel "email-marketing-service/api/v1/model/admin"
	"email-marketing-service/api/v1/repository"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

// AdminUsersRepository handles database operations related to Admin and User management
type AdminUsersRepository struct {
	DB *gorm.DB
}

// NewAdminUserRepository initializes a new AdminUsersRepository
func NewAdminUsersRepository(db *gorm.DB) *AdminUsersRepository {
	return &AdminUsersRepository{
		DB: db,
	}
}

// ConvertUserToUserResponse converts a User model into a UserResponse struct
func ConvertUserToUserResponse(user model.User) model.UserResponse {
	verifiedAt := ""
	if user.VerifiedAt != nil {
		verifiedAt = user.VerifiedAt.Format(time.RFC3339)
	}
	deletedAt := ""
	if user.DeletedAt.Valid {
		deletedAt = user.DeletedAt.Time.Format(time.RFC3339)
	}

	return model.UserResponse{
		UUID:        user.UUID,
		FullName:    user.FullName,
		Email:       user.Email,
		Company:     user.Company,
		PhoneNumber: user.PhoneNumber,
		Verified:    user.Verified,
		Blocked:     user.Blocked,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
		VerifiedAt:  &verifiedAt,
		DeletedAt:   &deletedAt,
	}
}

// GetAllUsers retrieves all users and converts them to UserResponse
func (r *AdminUsersRepository) GetAllUsers(search string, params repository.PaginationParams) (repository.PaginatedResult, error) {
	var users []model.User
	result := r.DB.Find(&users)

	if search != "" {
		result = result.Where("full_name ILIKE ? OR email ILIKE ? OR company ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if result.Error != nil {
		return repository.PaginatedResult{}, result.Error
	}

	paginatedResult, err := repository.Paginate(result, params, &users)

	if err != nil {
		return repository.PaginatedResult{}, fmt.Errorf("failed to paginate contacts: %w", err)
	}

	// Convert []User to []UserResponse
	var userResponses []model.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, ConvertUserToUserResponse(user))
	}

	paginatedResult.Data = userResponses

	return paginatedResult, nil
}

// BlockUser blocks a user by their UUID and returns the updated UserResponse
func (r *AdminUsersRepository) BlockUser(userId string) (*model.UserResponse, error) {
	var user model.User
	result := r.DB.Where("uuid = ?", userId).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	user.Blocked = true
	if err := r.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	userResponse := ConvertUserToUserResponse(user)
	return &userResponse, nil
}

// UnblockUser unblocks a user by their UUID and returns the updated UserResponse
func (r *AdminUsersRepository) UnblockUser(userId string) (*model.UserResponse, error) {
	var user model.User
	result := r.DB.Where("uuid = ?", userId).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	user.Blocked = false
	if err := r.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	userResponse := ConvertUserToUserResponse(user)
	return &userResponse, nil
}

// GetVerifiedUsers retrieves all verified users and returns a list of UserResponse
func (r *AdminUsersRepository) GetVerifiedUsers(search string, params repository.PaginationParams) (repository.PaginatedResult, error) {
	var users []model.User
	result := r.DB.Where("verified = ?", true).Find(&users)

	if search != "" {
		result = result.Where("full_name ILIKE ? OR email ILIKE ? OR company ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if result.Error != nil {
		return repository.PaginatedResult{}, result.Error
	}

	paginatedResult, err := repository.Paginate(result, params, &users)

	if err != nil {
		return repository.PaginatedResult{}, fmt.Errorf("failed to paginate contacts: %w", err)
	}

	// Convert []User to []UserResponse
	var userResponses []model.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, ConvertUserToUserResponse(user))
	}

	paginatedResult.Data = userResponses

	return paginatedResult, nil
}

func (r *AdminUsersRepository) GetUnVerifiedUsers(search string, params repository.PaginationParams) (repository.PaginatedResult, error) {
	var users []model.User
	result := r.DB.Where("verified = ?", false).Find(&users)

	if search != "" {
		result = result.Where("full_name ILIKE ? OR email ILIKE ? OR company ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if result.Error != nil {
		return repository.PaginatedResult{}, result.Error
	}

	paginatedResult, err := repository.Paginate(result, params, &users)

	if err != nil {
		return repository.PaginatedResult{}, fmt.Errorf("failed to paginate contacts: %w", err)
	}

	// Convert []User to []UserResponse
	var userResponses []model.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, ConvertUserToUserResponse(user))
	}

	paginatedResult.Data = userResponses

	return paginatedResult, nil
}

// GetSingleUser retrieves a single user by UUID and returns the UserResponse
func (r *AdminUsersRepository) GetSingleUser(userId string) (*model.UserResponse, error) {
	var user model.User
	result := r.DB.Where("uuid = ?", userId).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	userResponse := ConvertUserToUserResponse(user)
	return &userResponse, nil
}

// VerifyUser marks a user as verified by their UUID and returns the updated UserResponse
func (r *AdminUsersRepository) VerifyUser(userId string) (*model.UserResponse, error) {
	var user model.User
	result := r.DB.Where("uuid = ?", userId).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	// Check if the user is already verified
	if user.Verified {
		return nil, errors.New("user is already verified")
	}

	// Mark user as verified and set VerifiedAt to current time
	now := time.Now()
	user.Verified = true
	user.VerifiedAt = &now

	if err := r.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	userResponse := ConvertUserToUserResponse(user)
	return &userResponse, nil
}

type AdminUserStats struct {
	TotalContacts      int64 `json:"total_contacts"`
	TotalCampaigns     int64 `json:"total_campaigns"`
	TotalTemplates     int64 `json:"total_templates"`
	TotalCampaignsSent int64 `json:"total_campaigns_sent"`
	//	TotalSubscriptions int64 `json:"total_subscriptions"`
	TotalGroups int64 `json:"total_groups"`
}

func (r *AdminUsersRepository) GetUserStats(userID string) (AdminUserStats, error) {
	var stats AdminUserStats

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

func (r *AdminUsersRepository) AllUsersEmail() ([]model.UserResponse, error) {
	var users []model.User
	result := r.DB.Select("email").Where("verified = ?", true).Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	// Convert []User to []UserResponse
	var userResponses []model.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, ConvertUserToUserResponse(user))
	}

	return userResponses, nil
}

func (r *AdminUsersRepository) SaveMailLog(d adminmodel.AdminMailLog) error {
	if err := r.DB.Create(&d).Error; err != nil {
		return fmt.Errorf("failed to insert key: %w", err)
	}
	return nil
}

func (r *AdminUsersRepository) DeleteUser(userId string) error {
	// Start a transaction to ensure data consistency
	tx := r.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Define cleanup in case of error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// First verify the user exists
	var user model.User
	if err := tx.Where("uuid = ?", userId).First(&user).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Delete associated records in order to maintain referential integrity
	// Delete user's contacts
	if err := tx.Where("user_id = ?", userId).Delete(&model.Contact{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete contacts: %w", err)
	}

	// Delete user's contact groups
	if err := tx.Where("user_id = ?", userId).Delete(&model.ContactGroup{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete contact groups: %w", err)
	}

	// Delete user's templates
	if err := tx.Where("user_id = ?", userId).Delete(&model.Template{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete templates: %w", err)
	}

	// Delete user's campaigns
	if err := tx.Where("user_id = ?", userId).Delete(&model.Campaign{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete campaigns: %w", err)
	}

	// Finally, delete the user
	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *AdminUsersRepository) GetPendingDeletions() ([]model.User, error) {
	var users []model.User
	err := r.DB.Where("scheduled_for_deletion = ? AND scheduled_deletion_at <= ?",
		true, time.Now()).Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get pending deletions: %w", err)
	}
	return users, nil
}
