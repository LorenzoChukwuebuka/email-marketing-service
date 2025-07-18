package dto

type AdminFetchCampaignDTO struct {
	Search    string `json:"search"`
	Offset    int    `json:"offset"`
	Limit     int    `json:"limit"`
	CompanyID string `json:"company_id"`
	UserID    string `json:"user_id"`
	CampaignID string `json:"campaign_id"`
}
