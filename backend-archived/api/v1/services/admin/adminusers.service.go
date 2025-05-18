package adminservice

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	adminrepository "email-marketing-service/api/v1/repository/admin"
	"email-marketing-service/api/v1/utils"
	"fmt"
	"sync"
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

func (s *AdminUsers) SendEmailToUsers(d *dto.AdminEmailLogDTO) error {

	if err := utils.ValidateData(d); err != nil {
		return fmt.Errorf("invalid data: %w", err)
	}

	// Fetch all users' emails from the repository
	users, err := s.AdminUserRepo.AllUsersEmail()
	if err != nil {
		return err
	}

	// Initialize a wait group to manage goroutines
	var wg sync.WaitGroup

	// Channel to collect errors from goroutines
	errChan := make(chan error, len(users))

	// Loop through the users to send emails concurrently
	for _, user := range users {
		wg.Add(1) // Increment the wait group counter for each goroutine

		// Launch a goroutine to send emails
		go func(email string) {
			defer wg.Done() // Decrement the counter when goroutine finishes

			// Print user email (for logging/debugging purposes)
			print(email)

			// Call the SendMail utility to send the email
			if err := utils.AsyncSendMail(d.Subject, email, d.Content, "info@crabmailer.com", nil, &wg); err != nil {
				// Send error to the channel if sending fails
				errChan <- err
			}
		}(user.Email)
	}

	// Close the error channel once all goroutines finish
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Collect any errors from the channel
	for err := range errChan {
		if err != nil {
			return err // Return the first error encountered
		}
	}

	// Return nil if no errors occurred
	return nil
}
