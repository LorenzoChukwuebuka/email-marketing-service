package services

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/utils"
	"fmt"
)

type TransactionService struct {
}

func NewTransactionService() *TransactionService {
	return &TransactionService{}
}

func (s *TransactionService) InitiateNewTransaction(d *model.InitPaymentModelData) (map[string]interface{}, error) {
	if err := utils.ValidateData(d); err != nil {
		return nil, err
	}

	//instantiate the payment processor

	newTransaction := Transaction{}

	if err := newTransaction.ChoosePaymentMethod(d.PaymentMethod); err != nil {
		return nil, err
	}

	//choose the method to execute

	result, err := newTransaction.OpenProcessPayment(d)

	if err != nil {
		return nil, fmt.Errorf("error initiating payment: %s", err)
	}

	fmt.Print(result)

	return nil, nil
}
