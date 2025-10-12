package domain

type BasePaymentModelData struct {
	CompanyId       string  `json:"company_id" validate:"required"`
	UserId          string  `json:"user_id" validate:"required"`
	Email           string  `json:"email"`
	AmountToPay     float64 `json:"amount_to_pay" validate:"required"`
	PaymentMethod   string  `json:"payment_method" validate:"required"`
	Duration        string  `json:"duration" validate:"required"`
	PlanId          string  `json:"plan_id" validate:"required"`
	PaymentIntentID string  `json:"payment_intent_id"`
}

type BaseProcessPaymentModel struct {
	PaymentMethod string `json:"payment_method" validate:"required"`
	Reference     string `json:"reference" validate:"required"`
}

type BasePaymentResponse struct {
	Amount          float64 `json:"amount"`
	PlanID          string  `json:"plan_id"`
	UserID          string  `json:"user_id"`
	Duration        string  `json:"duration"`
	Email           string  `json:"email"`
	Status          string  `json:"status"`
	CompanyID       string  `json:"company_id"`
	PaymentIntentID string  `json:"payment_intent_id"`
	Reference       string  `json:"reference"`
}
