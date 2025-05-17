package mappers

import (
	"email-marketing-service/core/handler/admin/systems/dto"
	db "email-marketing-service/internal/db/sqlc"
)

func MapSMTPSettingsToDTO(settings db.SystemsSmtpSetting) *dto.SystemsResponse {
	return &dto.SystemsResponse{
		Domain:         settings.Domain.String,
		TXTRecord:      settings.TxtRecord.String,
		DMARCRecord:    settings.DmarcRecord.String,
		DKIMSelector:   settings.DkimSelector.String,
		DKIMPublicKey:  settings.DkimPublicKey.String,
		DKIMPrivateKey: settings.DkimPrivateKey.String,
		SPFRecord:      settings.SpfRecord.String,
		MXRecord:       settings.MxRecord.String,
		Verified:       settings.Verified.Bool,
	}
}
