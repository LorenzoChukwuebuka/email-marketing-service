package adminrepository

import (
	adminmodel "email-marketing-service/api/v1/model/admin"
	"fmt"
	"gorm.io/gorm"
)

type AdminRepository struct {
	DB *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{DB: db}
}

func (r *AdminRepository) createAdminResponse(admin adminmodel.Admin) *adminmodel.AdminResponse {
	return &adminmodel.AdminResponse{
		ID:        admin.ID,
		UUID:      admin.UUID,
		FirstName: admin.FirstName,
		LastName:  admin.LastName,
		Password:  admin.Password,
		Email:     admin.Email,
		Type:      admin.Type,
	}
}

func (r *AdminRepository) CreateAdmin(d *adminmodel.Admin) (*adminmodel.Admin, error) {

	if err := r.DB.Create(&d).Error; err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return d, nil
}

func (r *AdminRepository) Login(email string) (*adminmodel.AdminResponse, error) {
	var admin adminmodel.Admin

	if err := r.DB.Where("email = ?", email).First(&admin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error querying database: %w", err)
	}

	response := r.createAdminResponse(admin)
	return response, nil

}

// FindUserById finds an admin by UUID and returns the admin response
func (r *AdminRepository) FindAdminById(uuid string) (adminmodel.AdminResponse, error) {
	var admin adminmodel.Admin

	if err := r.DB.Where("uuid = ?", uuid).First(&admin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return adminmodel.AdminResponse{}, fmt.Errorf("admin not found")
		}
		return adminmodel.AdminResponse{}, fmt.Errorf("error querying database: %w", err)
	}

	adminResponse := r.createAdminResponse(admin)
	return *adminResponse, nil
}

// ChangePassword updates the admin's password by UUID
func (r *AdminRepository) ChangePassword(uuid string, newPassword string) error {
	var admin adminmodel.Admin

	// Find the admin by UUID
	if err := r.DB.Where("uuid = ?", uuid).First(&admin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("admin not found")
		}
		return fmt.Errorf("error querying database: %w", err)
	}

	// Update the password
	admin.Password = newPassword
	if err := r.DB.Save(&admin).Error; err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}
