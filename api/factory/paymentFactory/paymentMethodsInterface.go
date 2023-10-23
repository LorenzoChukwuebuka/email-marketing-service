package paymentmethodFactory

import "email-marketing-service/api/model"

type PaymentInterface interface {
	OpenDeposit(d *model.BasePaymentModelData) (map[string]interface{}, error)
	ProcessDeposit( d *model.BaseProcessPaymentModel) (*model.BasePaymentResponse,error)
	OpenRefund()
	ProcessRefund()
	ChargeCard(amount float64)
	Status() string
}
