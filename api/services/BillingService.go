package services

import (
	paymentmethodFactory "email-marketing-service/api/factory/paymentFactory"
	"email-marketing-service/api/model"
	"fmt"
	"time"
	"strconv"
	"strings"
)

type BillingService struct{}

func NewBillingService() *BillingService {
	return &BillingService{}
}

func (s *BillingService) ConfirmPayment(paymentmethod string, reference string) (map[string]interface{}, error) {

	paymenservice, err := paymentmethodFactory.PaymentFactory(paymentmethod)

	if err != nil {
		return nil, fmt.Errorf("error instantiating factory: %s", err)
	}

	params := &model.BaseProcessPaymentModel{
		PaymentMethod: paymentmethod,
		Reference:     reference,
	}

	data, err := paymenservice.ProcessDeposit(params)

	if err != nil {
		return nil, err
	}

	fmt.Println(data)

	return nil, nil
}


func ckalculateExpiryDate(duration string) time.Time {
	parts := strings.Split(duration, " ")
	num, _ := strconv.Atoi(parts[0])
	unit := parts[1]

	switch unit {
	case "week":
		return time.Now().AddDate(0, 0, num*7)
	case "day":
		return time.Now().AddDate(0, 0, num)
	case "year":
		return time.Now().AddDate(num, 0, 0)
	case "month":
		return time.Now().AddDate(0, num, 0)
	default:
		return time.Now()
	}
}
