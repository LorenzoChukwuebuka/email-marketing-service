package dto


type BasePaymentModelData struct {
	UserId        string     `json:"user_id" validate:"required"`
	Email         string  `json:"email"`
	AmountToPay   float64 `json:"amount_to_pay" validate:"required"`
	PaymentMethod string  `json:"payment_method" validate:"required"`
	Duration      string  `json:"duration" validate:"required"`
	PlanId        int     `json:"plan_id" validate:"required"`
}

type BaseProcessPaymentModel struct {
	PaymentMethod string `json:"payment_method" validate:"required"`
	Reference     string `json:"reference" validate:"required"`
}

type BasePaymentResponse struct {
	Amount   float64
	PlanID   int
	UserID   string
	Duration string
	Email    string
	Status   string
}
