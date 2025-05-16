package adminservice

import (
	"email-marketing-service/api/v1/dto"
	adminmodel "email-marketing-service/api/v1/model/admin"
	adminrepository "email-marketing-service/api/v1/repository/admin"
	"email-marketing-service/api/v1/utils"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrGeneratingToken = errors.New("error generating JWT token")

type AdminService struct {
	AdminRepo *adminrepository.AdminRepository
}

func NewAdminService(adminRepo *adminrepository.AdminRepository) *AdminService {
	return &AdminService{AdminRepo: adminRepo}
}

func (s *AdminService) CreateAdmin(d *dto.Admin) (*adminmodel.Admin, error) {

	if err := utils.ValidateData(d); err != nil {
		return nil, err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(d.Password), bcrypt.DefaultCost)

	adminModel := &adminmodel.Admin{
		UUID:       uuid.New().String(),
		FirstName:  d.FirstName,
		MiddleName: d.MiddleName,
		LastName:   d.LastName,
		Email:      d.Email,
		Type:       "admin",
		Password:   string(password),
	}

	adminUser, err := s.AdminRepo.CreateAdmin(adminModel)

	if err != nil {
		return nil, err
	}

	return adminUser, nil
}

func (s *AdminService) AdminLogin(d *dto.AdminLogin) (map[string]interface{}, error) {
	if err := utils.ValidateData(d); err != nil {
		return nil, err
	}

	adminDetails, err := s.AdminRepo.Login(d.Email)

	if err != nil {
		return nil, fmt.Errorf("invalid email:%w", err)
	}

	//compare password
	if err = bcrypt.CompareHashAndPassword([]byte(adminDetails.Password), []byte(d.Password)); err != nil {
		return nil, fmt.Errorf("passwords do not match:%w", err)
	}

	token, err := utils.GenerateAdminAccessToken(adminDetails.UUID, adminDetails.UUID, adminDetails.Type, adminDetails.Email)

	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateAdminRefreshToken(adminDetails.UUID, adminDetails.UUID, adminDetails.Type, adminDetails.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGeneratingToken, err)
	}

	successMap := map[string]interface{}{
		"status":        "login successful",
		"token":         token,
		"details":       adminDetails,
		"refresh_token": refreshToken,
		"type":          adminDetails.Type,
	}

	return successMap, nil

}

func (s *AdminService) ChangePassword() error {
	return nil
}
