package worker

import (
	"email-marketing-service/internal/domain"
	"net"
	"github.com/google/uuid"
)

type EmailPayload struct {
	Subject string `json:"subject"`
	Email   string `json:"email"`
	Message string `json:"message"`
	Sender  string `json:"sender"`
}

// payloads for email (this is an example)
type EmailExamplePayload struct {
	Email string
	Name  string
}

type AdminNotificationPayload struct {
	UserId            uuid.UUID
	Link              string
	NotificationTitle string
}

type UserNotificationPayload struct {
	UserId           uuid.UUID
	NotifcationTitle string
	AdditionalField  *string
}

type UserDetailsPayload struct {
	Details interface{} `json:"details"`
}

type Change struct {
	Old any `json:"old"`
	New any `json:"new"`
}

type AuditCreatePayload struct {
	UserID      uuid.UUID   `json:"user_id"`
	Resource    string      `json:"resource"`
	ResourceID  *uuid.UUID  `json:"resource_id,omitempty"`
	Method      string      `json:"method"`
	Endpoint    string      `json:"endpoint"`
	IP          net.IP      `json:"ip"`
	Success     bool        `json:"success"`
	RequestBody interface{} `json:"request_body,omitempty"`
}

type AuditUpdatePayload struct {
	UserID      uuid.UUID              `json:"user_id"`
	Resource    string                 `json:"resource"`
	ResourceID  *uuid.UUID             `json:"resource_id,omitempty"`
	Method      string                 `json:"method"`
	Endpoint    string                 `json:"endpoint"`
	IP          net.IP                 `json:"ip"`
	Success     bool                   `json:"success"`
	RequestBody map[string]interface{} `json:"request_body,omitempty"`
	Changes     map[string]Change      `json:"changes,omitempty"`
}

type AuditDeletePayload struct {
	UserID      uuid.UUID              `json:"user_id"`
	Resource    string                 `json:"resource"`
	ResourceID  *uuid.UUID             `json:"resource_id,omitempty"`
	Method      string                 `json:"method"`
	Endpoint    string                 `json:"endpoint"`
	IP          net.IP                 `json:"ip"`
	Success     bool                   `json:"success"`
	RequestBody map[string]interface{} `json:"request_body,omitempty"`
}

type AuditLoginPayload struct {
	UserID   uuid.UUID `json:"user_id"`
	Method   string    `json:"method"`
	Endpoint string    `json:"endpoint"`
	IP       net.IP    `json:"ip"`
	Username string    `json:"username"`
}

type AuditFailedLoginPayload struct {
	AttemptedUsername string `json:"attempted_username"`
	Method            string `json:"method"`
	Endpoint          string `json:"endpoint"`
	IP                net.IP `json:"ip"`
}

type SendCampaignEmailPayload struct {
	CompanyID  uuid.UUID `json:"company_id"`
	UserID     uuid.UUID `json:"user_id"`
	CampaignID uuid.UUID `json:"campaign_id"`
}

type SendAPISMTPEmailsPayload struct {
	EmailPayload domain.EmailRequest `json:"email_payload"`
	CompanyId    uuid.UUID           `json:"company_id"`
	UserId       uuid.UUID           `json:"user_id"`
}
