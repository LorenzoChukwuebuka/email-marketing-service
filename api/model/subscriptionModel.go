package model

import (
	"time"
)

type SubscriptionModel struct {
	Id        int
	UserId    int
	PlanId    int
	PaymentId int
	StartDate time.Time
	EndDate   time.Time
	Expired   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
