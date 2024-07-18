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

func (r *AdminRepository) ChangePassword() {

}
