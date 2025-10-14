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
		UserID_2:                  row.UserID_2.UUID,
		UserFullname:              row.UserFullname.String,
		UserEmail:                 row.UserEmail.String,
		UserPhonenumber:           row.UserPhonenumber,
		UserPicture:               row.UserPicture,
		UserVerified:              row.UserVerified.Bool,
		UserBlocked:               row.UserBlocked.Bool,
		UserVerifiedAt:            row.UserVerifiedAt,
		UserStatus:                row.UserStatus.String,
		UserLastLoginAt:           row.UserLastLoginAt,
		UserCreatedAt:             row.UserCreatedAt.Time,
		UserUpdatedAt:             row.UserUpdatedAt.Time,
		CompanyIDRef:              row.CompanyIDRef.UUID,
		CompanyName:               row.CompanyName,
		CompanyCreatedAt:          row.CompanyCreatedAt.Time,
		CompanyUpdatedAt:          row.CompanyUpdatedAt.Time,
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
	// Check TemplateID instead of TemplateIDRef
	if !row.TemplateID.Valid || row.TemplateID.UUID == uuid.Nil {
		return nil
	}
	return &dto.TemplateResponse{
		ID:                nullUUIDToStringPtr(row.TemplateIDRef),
		UserID:            nullUUIDToStringPtr(row.TemplateUserID),
		CompanyID:         nullUUIDToStringPtr(row.TemplateCompanyID),
		Name:              nullStringToPtr(row.TemplateName),
		SenderName:        nullStringToPtr(row.TemplateSenderName),
		FromEmail:         nullStringToPtr(row.TemplateFromEmail),
		Subject:           nullStringToPtr(row.TemplateSubject),
		Type:              nullStringToPtr(row.TemplateType),
		EmailHtml:         nullStringToPtr(row.TemplateEmailHtml),
		EmailDesign:       nullRawMessageToPtr(row.TemplateEmailDesign),
		IsEditable:        nullBoolToBool(row.TemplateIsEditable),
		IsPublished:       nullBoolToBool(row.TemplateIsPublished),
		IsPublicTemplate:  nullBoolToBool(row.TemplateIsPublicTemplate),
		IsGalleryTemplate: nullBoolToBool(row.TemplateIsGalleryTemplate),
		Tags:              nullStringToPtr(row.TemplateTags),
		Description:       nullStringToPtr(row.TemplateDescription),
		ImageUrl:          nullStringToPtr(row.TemplateImageUrl),
		IsActive:          nullBoolToBool(row.TemplateIsActive),
		EditorType:        nullStringToPtr(row.TemplateEditorType),
		CreatedAt:         nullTimeToPtr(row.TemplateCreatedAt),
		UpdatedAt:         nullTimeToPtr(row.TemplateUpdatedAt),
		DeletedAt:         nullTimeToPtr(row.TemplateDeletedAt),
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

func nullInt32ToPtr(ni sql.NullInt32) *int32 {
	if !ni.Valid {
		return nil
	}
	return &ni.Int32
}


func MapCampaignEmailResponse(req []db.EmailCampaignResult) []*dto.EmailCampaignResultResponse {
	if len(req) == 0 {
		return nil
	}

	var result []*dto.EmailCampaignResultResponse
	for _, r := range req {
		result = append(result, &dto.EmailCampaignResultResponse{
			ID:              r.ID.String(),
			CompanyID:       r.CompanyID.String(),
			CampaignID:      r.CampaignID.String(),
			RecipientEmail:  r.RecipientEmail,
			RecipientName:   nullStringToPtr(r.RecipientName),
			Version:         nullStringToPtr(r.Version),
			SentAt:          nullTimeToPtr(r.SentAt),
			OpenedAt:        nullTimeToPtr(r.OpenedAt),
			OpenCount:       nullInt32ToPtr(r.OpenCount),
			ClickedAt:       nullTimeToPtr(r.ClickedAt),
			ClickCount:      nullInt32ToPtr(r.ClickCount),
			ConversionAt:    nullTimeToPtr(r.ConversionAt),
			BounceStatus:    nullStringToPtr(r.BounceStatus),
			UnsubscribedAt:  nullTimeToPtr(r.UnsubscribedAt),
			ComplaintStatus: &r.ComplaintStatus.Valid,
			DeviceType:      nullStringToPtr(r.DeviceType),
			Location:        nullStringToPtr(r.Location),
			RetryCount:      nullInt32ToPtr(r.RetryCount),
			Notes:           nullStringToPtr(r.Notes),
			CreatedAt:       nullTimeToPtr(r.CreatedAt),
			UpdatedAt:       nullTimeToPtr(r.UpdatedAt),
			DeletedAt:       nullTimeToPtr(r.DeletedAt),
		})
	}

	return result
}

func MapCampaignGroups(row []db.GetCampaignContactGroupsRow) []*dto.GetCampaignContactGroupsResponse {
	var groups []*dto.GetCampaignContactGroupsResponse
	for _, r := range row {
		groups = append(groups, &dto.GetCampaignContactGroupsResponse{
			ID:          r.ID.String(),
			GroupName:   r.GroupName,
			Description: nullStringToPtr(r.Description),
			CreatedAt:   nullTimeToPtr(r.CreatedAt),
		})
	}
	return groups
}

// MapTemplateFromSeparateQuery maps template data from GetTemplateByIDWithoutType query
func MapTemplateFromSeparateQuery(template db.GetTemplateByIDWithoutTypeRow) *dto.TemplateResponse {
	return &dto.TemplateResponse{
		ID:                stringToPtr(template.ID.String()),
		UserID:            handleNullUUID(template.UserID),
		CompanyID:         handleNullUUID(template.CompanyID), // Fixed: Handle nullable UUID
		Name:              &template.TemplateName,
		SenderName:        nullStringToPtr(template.SenderName),
		FromEmail:         nullStringToPtr(template.FromEmail),
		Subject:           nullStringToPtr(template.Subject),
		Type:              &template.Type,
		EmailHtml:         nullStringToPtr(template.EmailHtml),
		EmailDesign:       nullRawMessageToPtr(template.EmailDesign),
		IsEditable:        nullBoolToBool(template.IsEditable),
		IsPublished:       nullBoolToBool(template.IsPublished),
		IsPublicTemplate:  nullBoolToBool(template.IsPublicTemplate),
		IsGalleryTemplate: nullBoolToBool(template.IsGalleryTemplate),
		Tags:              nullStringToPtr(template.Tags),
		Description:       nullStringToPtr(template.Description),
		ImageUrl:          nullStringToPtr(template.ImageUrl),
		IsActive:          nullBoolToBool(template.IsActive),
		EditorType:        nullStringToPtr(template.EditorType),
		CreatedAt:         nullTimeToPtr(template.CreatedAt),
		UpdatedAt:         nullTimeToPtr(template.UpdatedAt),
		DeletedAt:         nullTimeToPtr(template.DeletedAt),
	}
}

// handleNullUUID handles nullable UUID conversion to string pointer
func handleNullUUID(nu uuid.NullUUID) *string {
	if !nu.Valid {
		return nil
	}
	if nu.UUID == uuid.Nil {
		return nil
	}
	str := nu.UUID.String()
	return &str
}
// Helper function to convert string to pointer
func stringToPtr(s string) *string {
	return &s
}
