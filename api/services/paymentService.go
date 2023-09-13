package services

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/repository"
)

type PaymentService struct {
	PaymentRepo *repository.PaymentRepository
}

func NewPaymentService(paymentRepo *repository.PaymentRepository) *PaymentService {
	return &PaymentService{
		PaymentRepo: paymentRepo,
	}
}

func (s *PaymentService) InitializePayment(d interface{}) (map[string]interface{}, error) {
	return nil, nil
}

func (s *PaymentService) ConfirmPayment(d *model.PaymentModel) error {
	return nil
}

func (s *PaymentService) GetAllPaymentsForAUser(userId int) ([]model.PaymentResponse, error) {
	return nil, nil
}

func (s *PaymentService) GetSinglePaymentForAUser(userId int, paymentId int) (*model.PaymentResponse, error) {
	return nil, nil
}
