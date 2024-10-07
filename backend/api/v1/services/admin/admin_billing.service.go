package adminservice

import (
	"email-marketing-service/api/v1/model"
	adminrepository "email-marketing-service/api/v1/repository/admin"
	"time"
)

type AdminBillingService struct {
	AdminBillingRepo *adminrepository.AdminBillingRepository
}

func NewAdminBillingService(adminBillingRepo *adminrepository.AdminBillingRepository) *AdminBillingService {
	return &AdminBillingService{
		AdminBillingRepo: adminBillingRepo,
	}
}

func GetTimeRange(period string) (time.Time, time.Time) {
	now := time.Now()
	var startDate time.Time

	switch period {
	case "1day":
		startDate = now.AddDate(0, 0, -1)
	case "1week":
		startDate = now.AddDate(0, 0, -7)
	case "1month":
		startDate = now.AddDate(0, -1, 0)
	case "1year":
		startDate = now.AddDate(-1, 0, 0)
	default:
		startDate = now
	}

	return startDate, now
}

// Get the total billing (payments) for a specific user
func (s *AdminBillingService) GetTotalBillingForUser(userId uint) (float32, error) {
	total, err := s.AdminBillingRepo.GetTotalBillingForUser(userId)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// Get all the billings (payments) for a specific user
func (s *AdminBillingService) GetAllBillingsForUser(userId uint) ([]model.Billing, error) {
	billings, err := s.AdminBillingRepo.GetAllBillingsForUser(userId)
	if err != nil {
		return nil, err
	}
	return billings, nil
}

// Get the total of all billings within a given time range (1 day, 1 week, 1 month, 1 year)
func (s *AdminBillingService) GetTotalBillingsByPeriod(period string) (float32, error) {
	startDate, endDate := GetTimeRange(period)
	total, err := s.AdminBillingRepo.GetTotalBillingsByTimeRange(startDate, endDate)
	if err != nil {
		return 0, err
	}
	return total, nil
}
