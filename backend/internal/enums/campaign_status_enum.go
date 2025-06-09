package enums

type CampaignStatus string

const (
	Draft     CampaignStatus = "draft"
	Saved     CampaignStatus = "saved"
	Scheduled CampaignStatus = "scheduled"
	Sent      CampaignStatus = "sent"
	Queued CampaignStatus = "queued"
)
