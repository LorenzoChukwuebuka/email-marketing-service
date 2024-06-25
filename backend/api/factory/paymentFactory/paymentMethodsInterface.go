package paymentmethodFactory

import (
	"email-marketing-service/api/dto"	 
)

type PaymentInterface interface {
	OpenDeposit(d *dto.BasePaymentModelData) (map[string]interface{}, error)
	ProcessDeposit( d *dto.BaseProcessPaymentModel) (*dto.BasePaymentResponse,error)
	OpenRefund()
	ProcessRefund()
	ChargeCard(amount float64)
	Status() string
}


