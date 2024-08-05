package dto

type PlanStatus string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

type Plan struct {
	PlanName            string        `json:"planname" validate:"required"`
	Duration            string        `json:"duration" validate:"required"`
	Price               float32       `json:"price"`
	NumberOfMailsPerDay string        `json:"number_of_mails_per_day" validate:"required"`
	Details             string        `json:"details" validate:"required"`
	Status              PlanStatus    `json:"status" validate:"required"`
	Features            []PlanFeature `json:"features"`
}

type EditPlan struct {
	UUID                string        `json:"uuid"`
	PlanName            string        `json:"planname"`
	Duration            string        `json:"duration"`
	Price               float32       `json:"price"`
	NumberOfMailsPerDay string        `json:"number_of_mails_per_day"`
	Details             string        `json:"details"`
	Status              PlanStatus    `json:"status"`
	Features            []PlanFeature `json:"features"`
}

type PlanFeature struct {
	Name        string `json:"name" validate:"required"`
	Identifier  string `json:"identifier" validate:"required"`
	CountLimit  int    `json:"count_limit"`
	SizeLimit   int    `json:"size_limit"`
	IsActive    bool   `json:"is_active"`
	Description string `json:"description"`
}
