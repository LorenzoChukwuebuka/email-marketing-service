package paymentmethodFactory

import "email-marketing-service/api/model"

type PaymentInterface interface {
	OpenDeposit(d *model.InitPaymentModelData) (map[string]interface{}, error)
	ProcessDeposit(amount float64)
	OpenRefund()
	ProcessRefund()
	ChargeCard(amount float64)
	Status() string
}
