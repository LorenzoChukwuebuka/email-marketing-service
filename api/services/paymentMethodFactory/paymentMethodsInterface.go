package paymentmethodFactory

import "email-marketing-service/api/model"

type PaymentInterface interface {
	Initialize(d *model.PaymentModel)
	Pay(amount float64)
	Charge(amount float64)
	Refund(amount float64)
	Status() string
}
