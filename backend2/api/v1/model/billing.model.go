package model

import (
	"gorm.io/gorm"
	"time"
)
//TODO: Restructure this model to be better
type Billing struct {
	gorm.Model
	UUID          string    `json:"uuid"`
	UserId        uint      `json:"user_id"`
	AmountPaid    float32   `json:"amount_paid" validated:"required"`
	PlanId        uint      `json:"plan_id" validated:"required" gorm:"foreignKey:ID"`
	Email         string    `json:"email"`
	Duration      string    `json:"duration"`
	ExpiryDate    time.Time `json:"expiry_date"`
	Reference     string    `json:"reference"`
	TransactionId string    `json:"transaction_id"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	User          *User     `json:"user,omitempty" gorm:"foreignKey:UserId;references:ID"`
	Plan          *Plan     `json:"plan,omitempty" gorm:"foreignKey:PlanId;references:ID"`
}

type BillingResponse struct {
	Id            int           `json:"-"`
	UUID          string        `json:"uuid"`
	UserId        uint          `json:"user_id"`
	AmountPaid    float32       `json:"amount_paid"`
	PlanId        int           `json:"plan_id"`
	Duration      string        `json:"duration"`
	ExpiryDate    string        `json:"expiry_date"`
	Reference     string        `json:"reference"`
	TransactionId string        `json:"transaction_id"`
	PaymentMethod string        `json:"payment_method"`
	Status        string        `json:"status"`
	CreatedAt     string        `json:"created_at"`
	UpdatedAt     string        `json:"updated_at"`
	DeletedAt     *string       `json:"deleted_at"`
	User          *UserResponse `json:"user,omitempty"`
	Plan          *PlanResponse `json:"plan,omitempty"`
}

//TODO:  Add multitenancy and robustness using this model
//TODO : Do what you need to do
// type Subscription struct {
// 	gorm.Model
// 	UUID          string     `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4()"`
// 	CompanyID     uint       `json:"company_id"`
// 	PlanID        uint       `json:"plan_id"`
// 	Amount        float64    `json:"amount" gorm:"type:decimal(10,2);default:0.00"`
// 	BillingCycle  string     `json:"billing_cycle" gorm:"default:monthly"` // 'monthly', 'yearly'
// 	TrialStartsAt *time.Time `json:"trial_starts_at"`
// 	TrialEndsAt   *time.Time `json:"trial_ends_at"`
// 	StartsAt      *time.Time `json:"starts_at"`
// 	EndsAt        *time.Time `json:"ends_at"`
// 	Status        string     `json:"status" gorm:"default:inactive"` // 'active', 'inactive', 'cancelled', 'past_due'
	
// 	// Relationships
// 	Company  *Company  `json:"company" gorm:"foreignKey:CompanyID"`
// 	Plan     *Plan     `json:"plan" gorm:"foreignKey:PlanID"`
// 	Payments []Payment `json:"payments,omitempty" gorm:"foreignKey:SubscriptionID"`
	
// 	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
// }

// type Payment struct {
// 	gorm.Model
// 	SubscriptionID uint    `json:"subscription_id"`
// 	PaymentID      string  `json:"payment_id"` // ID from payment provider
// 	Amount         float64 `json:"amount" gorm:"type:decimal(10,2)"`
// 	Currency       string  `json:"currency"`
// 	PaymentMethod  string  `json:"payment_method"`
// 	Status         string  `json:"status"`
// 	Notes          string  `json:"notes" gorm:"type:text"`
	
// 	// Relationship
// 	Subscription *Subscription `json:"subscription" gorm:"foreignKey:SubscriptionID"`
// }

// type Plan struct {
// 	gorm.Model
// 	UUID         string        `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();index"`
// 	Name         string        `json:"name" gorm:"size:255"`
// 	Description  string        `json:"description" gorm:"type:text"`
// 	Price        float64       `json:"price" gorm:"type:decimal(10,2)"`
// 	BillingCycle string        `json:"billing_cycle" gorm:"default:monthly"` // 'monthly', 'yearly'
// 	Status       string        `json:"status" gorm:"size:255;default:active"`
	
// 	// Features and limits
// 	Features     []PlanFeature `json:"features" gorm:"foreignKey:PlanID"`
// 	MailingLimit MailingLimit  `json:"mailing_limit" gorm:"foreignKey:PlanID"`
// }

// type PlanFeature struct {
// 	gorm.Model
// 	PlanID      uint   `json:"plan_id"`
// 	Name        string `json:"name"`
// 	Description string `json:"description"`
// 	Value       string `json:"value"`
// }

// type MailingLimit struct {
// 	gorm.Model
// 	PlanID              uint `json:"plan_id"`
// 	DailyLimit          int  `json:"daily_limit"`
// 	MonthlyLimit        int  `json:"monthly_limit"`
// 	MaxRecipientsPerMail int  `json:"max_recipients_per_mail"`
// }

// type Company struct {
// 	gorm.Model
// 	UUID           string         `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4()"`
// 	Name           string         `json:"name"`
// 	Email          string         `json:"email" gorm:"unique"`
// 	// Other company fields
// 	Subscriptions []Subscription `json:"subscriptions,omitempty" gorm:"foreignKey:CompanyID"`
// }
