package adminservice

import (
	adminmodel "email-marketing-service/api/model/admin"
	adminrepository "email-marketing-service/api/repository/admin"
	"email-marketing-service/api/utils"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type AdminService struct {
	AdminRepo *adminrepository.AdminRepository
}

func NewAdminService(adminRepo *adminrepository.AdminRepository) *AdminService {
	return &AdminService{AdminRepo: adminRepo}
}

func (s *AdminService) AdminLogin(d *adminmodel.AdminLogin) (map[string]interface{}, error) {
	if err := utils.ValidateData(d); err != nil {
		return nil, err
	}

	adminDetails, err := s.AdminRepo.Login(d)

	if err != nil {
		return nil, fmt.Errorf("invalid email:%w", err)
	}

	//compare password
	if err = bcrypt.CompareHashAndPassword(adminDetails.Password, []byte(d.Password)); err != nil {
		return nil, fmt.Errorf("passwords do not match:%w", err)
	}

	token, err := utils.JWTEncode(adminDetails.ID, adminDetails.Type, adminDetails.Email)

	if err != nil {
		return nil, err
	}

	successMap := map[string]interface{}{
		"status":  "login successful",
		"token":   token,
		"details": adminDetails,
	}

	return successMap, nil

}

func (s *AdminService) ChangePassword() error {
	return nil
}
