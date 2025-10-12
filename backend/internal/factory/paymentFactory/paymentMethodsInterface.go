package paymentmethodFactory

import (
	"email-marketing-service/internal/domain"
	"net/http"
)

type PaymentInterface interface {
	OpenDeposit(d *domain.BasePaymentModelData) (map[string]interface{}, error)
	ProcessDeposit(d *domain.BaseProcessPaymentModel) (*domain.BasePaymentResponse, error)
	OpenRefund()
	ProcessRefund()
	ChargeCard(amount float64)
	Status() string
	WebhookHandler(r *http.Request) (*domain.BasePaymentResponse, error)
}
