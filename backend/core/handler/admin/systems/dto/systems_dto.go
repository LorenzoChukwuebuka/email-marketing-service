package dto

type SystemsDTO struct {
	Domain string `json:"domain" validate:"required"`
}

type SystemsResponse struct {
	UUID           string `json:"uuid"  `
	TXTRecord      string `json:"txt_record"`
	DMARCRecord    string `json:"dmarc_record"`
	DKIMSelector   string `json:"dkim_selector"`
	DKIMPublicKey  string `json:"dkim_public_key"`
	DKIMPrivateKey string `json:"dkim_private_key"`
	SPFRecord      string `json:"spf_record"`
	Verified       bool   `json:"verified"`
	MXRecord       string `json:"mx_record"`
	Domain         string `json:"domain"`
}
