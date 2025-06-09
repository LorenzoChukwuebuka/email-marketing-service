package mapper

import (
	"database/sql"
	"email-marketing-service/core/handler/campaigns/dto"
	db "email-marketing-service/internal/db/sqlc"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"time"
)

// MapScheduledCampaignResponse maps a single scheduled campaign row to CampaignResponseDTO
func MapScheduledCampaignResponse(row db.ListScheduledCampaignsByCompanyIDRow) *dto.CampaignResponseDTO {
	// Convert to the base type by creating a new struct with the same values
	baseRow := db.ListCampaignsByCompanyIDRow{
		ID:                        row.ID,
		CompanyID:                 row.CompanyID,
		Name:                      row.Name,
		Subject:                   row.Subject,
		PreviewText:               row.PreviewText,
		UserID:                    row.UserID,
		SenderFromName:            row.SenderFromName,
		TemplateID:                row.TemplateID,
		SentTemplateID:            row.SentTemplateID,
		RecipientInfo:             row.RecipientInfo,
		IsPublished:               row.IsPublished,
		Status:                    row.Status,
		TrackType:                 row.TrackType,
		IsArchived:                row.IsArchived,
		SentAt:                    row.SentAt,
		Sender:                    row.Sender,
		ScheduledAt:               row.ScheduledAt,
		HasCustomLogo:             row.HasCustomLogo,
		CreatedAt:                 row.CreatedAt,
		UpdatedAt:                 row.UpdatedAt,
		DeletedAt:                 row.DeletedAt,
		UserID_2:                  row.UserID_2,
		UserFullname:              row.UserFullname,
		UserEmail:                 row.UserEmail,
		UserPhonenumber:           row.UserPhonenumber,
		UserPicture:               row.UserPicture,
		UserVerified:              row.UserVerified,
		UserBlocked:               row.UserBlocked,
		UserVerifiedAt:            row.UserVerifiedAt,
		UserStatus:                row.UserStatus,
		UserLastLoginAt:           row.UserLastLoginAt,
		UserCreatedAt:             row.UserCreatedAt,
		UserUpdatedAt:             row.UserUpdatedAt,
		CompanyIDRef:              row.CompanyIDRef,
		CompanyName:               row.CompanyName,
		CompanyCreatedAt:          row.CompanyCreatedAt,
		CompanyUpdatedAt:          row.CompanyUpdatedAt,
		TemplateIDRef:             row.TemplateIDRef,
		TemplateUserID:            row.TemplateUserID,
		TemplateCompanyID:         row.TemplateCompanyID,
		TemplateName:              row.TemplateName,
		TemplateSenderName:        row.TemplateSenderName,
		TemplateFromEmail:         row.TemplateFromEmail,
		TemplateSubject:           row.TemplateSubject,
		TemplateType:              row.TemplateType,
		TemplateEmailHtml:         row.TemplateEmailHtml,
		TemplateEmailDesign:       row.TemplateEmailDesign,
		TemplateIsEditable:        row.TemplateIsEditable,
		TemplateIsPublished:       row.TemplateIsPublished,
		TemplateIsPublicTemplate:  row.TemplateIsPublicTemplate,
		TemplateIsGalleryTemplate: row.TemplateIsGalleryTemplate,
		TemplateTags:              row.TemplateTags,
		TemplateDescription:       row.TemplateDescription,
		TemplateImageUrl:          row.TemplateImageUrl,
		TemplateIsActive:          row.TemplateIsActive,
		TemplateEditorType:        row.TemplateEditorType,
		TemplateCreatedAt:         row.TemplateCreatedAt,
		TemplateUpdatedAt:         row.TemplateUpdatedAt,
		TemplateDeletedAt:         row.TemplateDeletedAt,
	}
	
	return MapCampaignResponse(baseRow)
}

// MapScheduledCampaignResponses maps a slice of scheduled campaign rows to slice of CampaignResponseDTO
func MapScheduledCampaignResponses(rows []db.ListScheduledCampaignsByCompanyIDRow) []dto.CampaignResponseDTO {
	result := make([]dto.CampaignResponseDTO, len(rows))
	for i, row := range rows {
		result[i] = *MapScheduledCampaignResponse(row)
	}
	return result
}

// MapGetCampaignResponse maps a single GetCampaignByID row to CampaignResponseDTO
func MapGetCampaignResponse(row db.GetCampaignByIDRow) *dto.CampaignResponseDTO {
	// Convert to the base type by creating a new struct with the same values
	baseRow := db.ListCampaignsByCompanyIDRow{
		ID:                        row.ID,
		CompanyID:                 row.CompanyID,
		Name:                      row.Name,
		Subject:                   row.Subject,
		PreviewText:               row.PreviewText,
		UserID:                    row.UserID,
		SenderFromName:            row.SenderFromName,
		TemplateID:                row.TemplateID,
		SentTemplateID:            row.SentTemplateID,
		RecipientInfo:             row.RecipientInfo,
		IsPublished:               row.IsPublished,
		Status:                    row.Status,
		TrackType:                 row.TrackType,
		IsArchived:                row.IsArchived,
		SentAt:                    row.SentAt,
		Sender:                    row.Sender,
		ScheduledAt:               row.ScheduledAt,
		HasCustomLogo:             row.HasCustomLogo,
		CreatedAt:                 row.CreatedAt,
		UpdatedAt:                 row.UpdatedAt,
		DeletedAt:                 row.DeletedAt,
		UserID_2:                  row.UserID_2,
		UserFullname:              row.UserFullname,
		UserEmail:                 row.UserEmail,
		UserPhonenumber:           row.UserPhonenumber,
		UserPicture:               row.UserPicture,
		UserVerified:              row.UserVerified,
		UserBlocked:               row.UserBlocked,
		UserVerifiedAt:            row.UserVerifiedAt,
		UserStatus:                row.UserStatus,
		UserLastLoginAt:           row.UserLastLoginAt,
		UserCreatedAt:             row.UserCreatedAt,
		UserUpdatedAt:             row.UserUpdatedAt,
		CompanyIDRef:              row.CompanyIDRef,
		CompanyName:               row.CompanyName,
		CompanyCreatedAt:          row.CompanyCreatedAt,
		CompanyUpdatedAt:          row.CompanyUpdatedAt,
		TemplateIDRef:             row.TemplateIDRef,
		TemplateUserID:            row.TemplateUserID,
		TemplateCompanyID:         row.TemplateCompanyID,
		TemplateName:              row.TemplateName,
		TemplateSenderName:        row.TemplateSenderName,
		TemplateFromEmail:         row.TemplateFromEmail,
		TemplateSubject:           row.TemplateSubject,
		TemplateType:              row.TemplateType,
		TemplateEmailHtml:         row.TemplateEmailHtml,
		TemplateEmailDesign:       row.TemplateEmailDesign,
		TemplateIsEditable:        row.TemplateIsEditable,
		TemplateIsPublished:       row.TemplateIsPublished,
		TemplateIsPublicTemplate:  row.TemplateIsPublicTemplate,
		TemplateIsGalleryTemplate: row.TemplateIsGalleryTemplate,
		TemplateTags:              row.TemplateTags,
		TemplateDescription:       row.TemplateDescription,
		TemplateImageUrl:          row.TemplateImageUrl,
		TemplateIsActive:          row.TemplateIsActive,
		TemplateEditorType:        row.TemplateEditorType,
		TemplateCreatedAt:         row.TemplateCreatedAt,
		TemplateUpdatedAt:         row.TemplateUpdatedAt,
		TemplateDeletedAt:         row.TemplateDeletedAt,
	}
	
	return MapCampaignResponse(baseRow)
}

// Original functions remain the same
func MapCampaignResponse(row db.ListCampaignsByCompanyIDRow) *dto.CampaignResponseDTO {
	return &dto.CampaignResponseDTO{
		ID:             row.ID.String(),
		CompanyID:      row.CompanyID.String(),
		Name:           row.Name,
		Subject:        nullStringToPtr(row.Subject),
		PreviewText:    nullStringToPtr(row.PreviewText),
		UserID:         row.UserID.String(),
		SenderFromName: nullStringToPtr(row.SenderFromName),
		TemplateID:     nullUUIDToStringPtr(row.TemplateID),
		SentTemplateID: nullUUIDToStringPtr(row.SentTemplateID),
		RecipientInfo:  nullStringToPtr(row.RecipientInfo),
		IsPublished:    nullBoolToBool(row.IsPublished),
		Status:         nullStringToPtr(row.Status),
		TrackType:      nullStringToPtr(row.TrackType),
		IsArchived:     nullBoolToBool(row.IsArchived),
		SentAt:         nullTimeToPtr(row.SentAt),
		Sender:         nullStringToPtr(row.Sender),
		ScheduledAt:    nullTimeToPtr(row.ScheduledAt),
		HasCustomLogo:  nullBoolToBool(row.HasCustomLogo),
		CreatedAt:      nullTimeToPtr(row.CreatedAt),
		UpdatedAt:      nullTimeToPtr(row.UpdatedAt),
		DeletedAt:      nullTimeToPtr(row.DeletedAt),
		User:           mapUserResponse(row),
		Company:        mapCompanyResponse(row),
		Template:       mapTemplateResponse(row),
	}
}

func MapCampaignResponses(rows []db.ListCampaignsByCompanyIDRow) []dto.CampaignResponseDTO {
	result := make([]dto.CampaignResponseDTO, len(rows))
	for i, row := range rows {
		result[i] = *MapCampaignResponse(row)
	}
	return result
}

func mapUserResponse(row db.ListCampaignsByCompanyIDRow) dto.UserResponse {
	return dto.UserResponse{
		UserID_2:        row.UserID_2.String(),
		UserFullname:    row.UserFullname,
		UserEmail:       row.UserEmail,
		UserPhonenumber: nullStringToPtr(row.UserPhonenumber),
		UserPicture:     nullStringToPtr(row.UserPicture),
		UserVerified:    row.UserVerified,
		UserBlocked:     row.UserBlocked,
		UserVerifiedAt:  nullTimeToPtr(row.UserVerifiedAt),
		UserStatus:      row.UserStatus,
		UserLastLoginAt: nullTimeToPtr(row.UserLastLoginAt),
		UserCreatedAt:   row.UserCreatedAt,
		UserUpdatedAt:   row.UserUpdatedAt,
	}
}

func mapCompanyResponse(row db.ListCampaignsByCompanyIDRow) dto.CompanyResponse {
	return dto.CompanyResponse{
		CompanyIDRef:     row.CompanyIDRef.String(),
		CompanyName:      nullStringToPtr(row.CompanyName),
		CompanyCreatedAt: row.CompanyCreatedAt,
		CompanyUpdatedAt: row.CompanyUpdatedAt,
	}
}

func mapTemplateResponse(row db.ListCampaignsByCompanyIDRow) *dto.TemplateResponse {
	if !row.TemplateIDRef.Valid || row.TemplateIDRef.UUID == uuid.Nil {
		return nil
	}
	return &dto.TemplateResponse{
		TemplateIDRef:             nullUUIDToStringPtr(row.TemplateIDRef),
		TemplateUserID:            nullUUIDToStringPtr(row.TemplateUserID),
		TemplateCompanyID:         nullUUIDToStringPtr(row.TemplateCompanyID),
		TemplateName:              nullStringToPtr(row.TemplateName),
		TemplateSenderName:        nullStringToPtr(row.TemplateSenderName),
		TemplateFromEmail:         nullStringToPtr(row.TemplateFromEmail),
		TemplateSubject:           nullStringToPtr(row.TemplateSubject),
		TemplateType:              nullStringToPtr(row.TemplateType),
		TemplateEmailHtml:         nullStringToPtr(row.TemplateEmailHtml),
		TemplateEmailDesign:       nullRawMessageToPtr(row.TemplateEmailDesign),
		TemplateIsEditable:        nullBoolToBool(row.TemplateIsEditable),
		TemplateIsPublished:       nullBoolToBool(row.TemplateIsPublished),
		TemplateIsPublicTemplate:  nullBoolToBool(row.TemplateIsPublicTemplate),
		TemplateIsGalleryTemplate: nullBoolToBool(row.TemplateIsGalleryTemplate),
		TemplateTags:              nullStringToPtr(row.TemplateTags),
		TemplateDescription:       nullStringToPtr(row.TemplateDescription),
		TemplateImageUrl:          nullStringToPtr(row.TemplateImageUrl),
		TemplateIsActive:          nullBoolToBool(row.TemplateIsActive),
		TemplateEditorType:        nullStringToPtr(row.TemplateEditorType),
		TemplateCreatedAt:         nullTimeToPtr(row.TemplateCreatedAt),
		TemplateUpdatedAt:         nullTimeToPtr(row.TemplateUpdatedAt),
		TemplateDeletedAt:         nullTimeToPtr(row.TemplateDeletedAt),
	}
}

// Helper functions
func nullStringToPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

func nullUUIDToStringPtr(nu uuid.NullUUID) *string {
	if !nu.Valid {
		return nil
	}
	if nu.UUID == uuid.Nil {
		return nil
	}
	str := nu.UUID.String()
	return &str
}

func nullTimeToPtr(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}

func nullBoolToBool(nb sql.NullBool) bool {
	if !nb.Valid {
		return false
	}
	return nb.Bool
}

func nullRawMessageToPtr(nrm pqtype.NullRawMessage) *json.RawMessage {
	if !nrm.Valid {
		return nil
	}
	return &nrm.RawMessage
}