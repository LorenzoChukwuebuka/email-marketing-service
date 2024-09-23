package adminservice

import (
	"email-marketing-service/api/v1/repository"
	adminrepository "email-marketing-service/api/v1/repository/admin"
)

type AdminSupportService struct {
	AdminSupportRepo      *adminrepository.AdminSupportRepository
	AdminNotificationRepo *adminrepository.AdminNotificationRepository
	UserNotificationRepo  *repository.UserNotificationRepository
}

func NewAdminSupportService(adminsupportRepo *adminrepository.AdminSupportRepository,
	adminNotificationRepo *adminrepository.AdminNotificationRepository,
	usernotificationRepo *repository.UserNotificationRepository) *AdminSupportService {
	return &AdminSupportService{
		AdminSupportRepo:      adminsupportRepo,
		AdminNotificationRepo: adminNotificationRepo,
		UserNotificationRepo:  usernotificationRepo,
	}

}

func (s *AdminSupportService) GetAllTickets(search string, page int, pageSize int) (repository.PaginatedResult, error) {
	params := repository.PaginationParams{Page: page, PageSize: pageSize}
	tickets, err := s.AdminSupportRepo.GetAllTickets(search, params)

	if err != nil {
		return repository.PaginatedResult{}, err
	}
	return tickets, nil
}
