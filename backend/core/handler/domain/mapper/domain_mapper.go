package mapper

import (
	"database/sql"
	"email-marketing-service/core/handler/domain/dto"
	db "email-marketing-service/internal/db/sqlc"
	"time"
)

// MapDomain maps a single db.Domain to dto.DomainResponse
func MapDomain(domain db.Domain) *dto.DomainResponse {
	return &dto.DomainResponse{
		ID:             domain.ID.String(),
		UserID:         domain.UserID.String(),
		CompanyID:      domain.CompanyID.String(),
		Domain:         domain.Domain,
		TxtRecord:      nullStringToPtr(domain.TxtRecord),
		DmarcRecord:    nullStringToPtr(domain.DmarcRecord),
		DkimSelector:   nullStringToPtr(domain.DkimSelector),
		DkimPublicKey:  nullStringToPtr(domain.DkimPublicKey),
		DkimPrivateKey: nullStringToPtr(domain.DkimPrivateKey),
		SpfRecord:      nullStringToPtr(domain.SpfRecord),
		Verified:       domain.Verified.Bool,
		MxRecord:       nullStringToPtr(domain.MxRecord),
		CreatedAt:      nullTimeToPtr(domain.CreatedAt),
		UpdatedAt:      nullTimeToPtr(domain.UpdatedAt),
		DeletedAt:      nullTimeToPtr(domain.DeletedAt),
	}

}

// MapDomains maps a slice of db.Domain to a slice of dto.DomainResponse
func MapDomains(domains []db.Domain) []*dto.DomainResponse {
	if len(domains) == 0 {
		return []*dto.DomainResponse{}
	}

	result := make([]*dto.DomainResponse, len(domains))
	for i, domain := range domains {
		result[i] = MapDomain(domain)
	}
	return result
}

// Helper functions to convert sql.Null* types to pointers

func nullStringToPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func nullTimeToPtr(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}
