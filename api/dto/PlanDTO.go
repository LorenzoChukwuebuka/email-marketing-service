package dto

type status string

const (
	Active   status = "active"
	Inactive status = "inactive"
)

type Plan struct {
	PlanName            string  `json:"planname" validate:"required"`
	Duration            string  `json:"duration" validate:"required"`
	Price               float32 `json:"price"  `
	NumberOfMailsPerDay string  `json:"number_of_mails_per_day" validate:"required"`
	Details             string  `json:"details" validate:"required"`
	Status              *string `json:"status"`
}

type EditPlan struct {
	UUID                string  `json:"user_id"`
	PlanName            string  `json:"planname"`
	Duration            string  `json:"duration"`
	Price               float32 `json:"price" `
	NumberOfMailsPerDay string  `json:"number_of_mails_per_day"`
	Details             string  `json:"details"`
	Status              *string `json:"status"`
}
