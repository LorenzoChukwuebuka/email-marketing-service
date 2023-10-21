package services

import (
	paymentmethodFactory "email-marketing-service/api/factory/paymentFactory"
	"fmt"
)

type BillingService struct{}

func NewBillingService() *BillingService {
	return &BillingService{}
}

func (s *BillingService) ConfirmPayment(paymentmethod string, reference string) (string, error) {

	paymenservice, err := paymentmethodFactory.PaymentFactory(paymentmethod)

	if err != nil {
		return "", fmt.Errorf("error instantiating factory: %s", err)
	}

	paymenservice.ProcessDeposit(60.90)

	return "", nil
}
