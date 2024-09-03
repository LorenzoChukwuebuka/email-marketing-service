package repository

import (
	"email-marketing-service/api/v1/model"
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

func (r *UserRepository) CreateTempEmail(d *model.UserTempEmail) (  error) {
	if err := r.DB.Create(&d).Error; err != nil {
		return   fmt.Errorf("failed to insert user: %w", err)
	}
	return  nil

}
