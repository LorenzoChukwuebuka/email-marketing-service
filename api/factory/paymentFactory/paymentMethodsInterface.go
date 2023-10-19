package paymentmethodFactory

import "email-marketing-service/api/model"

type PaymentInterface interface {
	InitializePaymentProcess(d *model.InitPaymentModelData) (map[string]interface{}, error)
	Pay(amount float64)
	Charge(amount float64)
	Refund(amount float64)
	Status() string
}
