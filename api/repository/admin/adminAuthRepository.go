package adminrepository

import (
	adminmodel "email-marketing-service/api/model/admin"
	"fmt"
	"gorm.io/gorm"
)

type AdminRepository struct {
	DB *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{DB: db}
}

func (r *AdminRepository) CreateAdmin(d *adminmodel.Admin) (*adminmodel.Admin, error) {

	if err := r.DB.Create(&d).Error; err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return d, nil
}

func (r *AdminRepository) Login(d *adminmodel.AdminLogin) (*adminmodel.AdminResponse, error) {

	return nil, nil

}

func (r *AdminRepository) ChangePassword() {

}
