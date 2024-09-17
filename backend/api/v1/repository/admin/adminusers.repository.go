package adminrepository

import "gorm.io/gorm"

type AdminUsersRepository struct {
	DB *gorm.DB
}

func NewAdminUserRepository(db *gorm.DB) *AdminUsersRepository {
	return &AdminUsersRepository{
		DB: db,
	}
}

func (r *AdminUsersRepository) GetAllUsers() error {
	return nil
}

func (r *AdminUsersRepository) BlockUser(userId string) error {
	return nil
}

func (r *AdminUsersRepository) GetSingleUser() error {
	return nil
}
