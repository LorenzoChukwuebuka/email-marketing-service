package adminservice

import (
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	adminrepository "email-marketing-service/api/v1/repository/admin"
)

type AdminUsers struct {
	AdminUserRepo *adminrepository.AdminUsersRepository
}

// NewAdminUsersStruct initializes a new AdminUsers service
func NewAdminUsersService(adminUserRepo *adminrepository.AdminUsersRepository) *AdminUsers {
	return &AdminUsers{
		AdminUserRepo: adminUserRepo,
	}
}

// GetAllUsers retrieves all users through the repository and returns a list of UserResponse
func (s *AdminUsers) GetAllUsers(page int, pageSize int, search string) (repository.PaginatedResult, error) {
	params := repository.PaginationParams{Page: page, PageSize: pageSize}
	users, err := s.AdminUserRepo.GetAllUsers(search, params)
	if err != nil {
		return repository.PaginatedResult{}, err
	}
	return users, nil
}

// BlockUser blocks a user by their UUID through the repository and returns the updated UserResponse
func (s *AdminUsers) BlockUser(userId string) (*model.UserResponse, error) {
	userResponse, err := s.AdminUserRepo.BlockUser(userId)
	if err != nil {
		return nil, err
	}
	return userResponse, nil
}

// UnblockUser unblocks a user by their UUID through the repository and returns the updated UserResponse
func (s *AdminUsers) UnblockUser(userId string) (*model.UserResponse, error) {
	userResponse, err := s.AdminUserRepo.UnblockUser(userId)
	if err != nil {
		return nil, err
	}
	return userResponse, nil
}

// GetVerifiedUsers retrieves all verified users through the repository and returns a list of UserResponse
func (s *AdminUsers) GetVerifiedUsers(page int, pageSize int, search string) (repository.PaginatedResult, error) {
	params := repository.PaginationParams{Page: page, PageSize: pageSize}
	users, err := s.AdminUserRepo.GetVerifiedUsers(search, params)
	if err != nil {
		return repository.PaginatedResult{}, err
	}
	return users, nil
}

func (s *AdminUsers) GetUnverifiedUsers(page int, pageSize int, search string) (repository.PaginatedResult, error) {
	params := repository.PaginationParams{Page: page, PageSize: pageSize}
	users, err := s.AdminUserRepo.GetUnVerifiedUsers(search, params)
	if err != nil {
		return repository.PaginatedResult{}, err
	}
	return users, nil
}

// GetSingleUser retrieves a single user by UUID through the repository and returns the UserResponse
func (s *AdminUsers) GetSingleUser(userId string) (*model.UserResponse, error) {
	user, err := s.AdminUserRepo.GetSingleUser(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// VerifyUser verifies a user by their UUID through the repository and returns the updated UserResponse
func (s *AdminUsers) VerifyUser(userId string) (*model.UserResponse, error) {
	user, err := s.AdminUserRepo.VerifyUser(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AdminUsers) GetUserStats(userId string) (*adminrepository.AdminUserStats, error) {
	stats, err := s.AdminUserRepo.GetUserStats(userId)

	if err != nil {
		return nil, err
	}

	return &stats, nil
}
