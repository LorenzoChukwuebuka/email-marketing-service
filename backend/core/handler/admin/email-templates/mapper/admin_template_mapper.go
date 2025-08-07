package mapper

import (
	"database/sql"
	"email-marketing-service/core/handler/admin/email-templates/dto"
	db "email-marketing-service/internal/db/sqlc"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type MappableTemplate interface {
	GetID() uuid.UUID
	GetUserID() uuid.NullUUID
	GetCompanyID() uuid.NullUUID
	GetTemplateName() string
	GetType() string
	GetSenderName() sql.NullString
	GetFromEmail() sql.NullString
	GetSubject() sql.NullString
	GetEmailHtml() sql.NullString
	GetTags() sql.NullString
	GetDescription() sql.NullString
	GetImageUrl() sql.NullString
	GetEditorType() sql.NullString
	GetIsEditable() sql.NullBool
	GetIsPublished() sql.NullBool
	GetIsPublicTemplate() sql.NullBool
	GetIsGalleryTemplate() sql.NullBool
	GetIsActive() sql.NullBool
	GetEmailDesign() interface{}
	GetCreatedAt() sql.NullTime
	GetUpdatedAt() sql.NullTime
	GetDeletedAt() sql.NullTime
	GetUserFullname() sql.NullString
	GetUserEmail() sql.NullString
	GetUserPicture() sql.NullString
	GetCompanyName() sql.NullString
}

func MapTemplateToDTO(dbTemplate db.Template) any {
	return &dto.AdminTemplateDTO{
		TemplateID: dbTemplate.ID.String(),
		TemplateName:      dbTemplate.TemplateName,
		SenderName:        dbTemplate.SenderName.String,
		FromEmail:         dbTemplate.FromEmail.String,
		Subject:           dbTemplate.Subject.String,
		Type:              dbTemplate.Type,
		EmailHtml:         dbTemplate.EmailHtml.String,
		EmailDesign:       dbTemplate.EmailDesign.RawMessage,
		IsEditable:        dbTemplate.IsEditable.Bool,
		IsPublished:       dbTemplate.IsPublished.Bool,
		IsPublicTemplate:  dbTemplate.IsPublicTemplate.Bool,
		IsGalleryTemplate: dbTemplate.IsGalleryTemplate.Bool,
		Tags:              dbTemplate.Tags.String,
		Description:       dbTemplate.Description.String,
		ImageUrl:          dbTemplate.ImageUrl.String,
		IsActive:          dbTemplate.IsActive.Bool,
		EditorType:        dbTemplate.EditorType.String,
	}
}

func MapTemplateResponse[T MappableTemplate](t []T) []dto.AdminTemplateResponse {
	if len(t) == 0 {
		return []dto.AdminTemplateResponse{}
	}

	var response []dto.AdminTemplateResponse

	for _, template := range t {
		templateResponse := dto.AdminTemplateResponse{
			ID:           template.GetID(),
			UserID:       getNullUUID(template.GetUserID()),
			CompanyID:    getNullUUID(template.GetCompanyID()),
			TemplateName: template.GetTemplateName(),
			Type:         template.GetType(),
			// Handle nullable string fields
			SenderName:  getNullStringValue(template.GetSenderName()),
			FromEmail:   getNullStringValue(template.GetFromEmail()),
			Subject:     getNullStringValue(template.GetSubject()),
			EmailHtml:   getNullStringValue(template.GetEmailHtml()),
			Tags:        getNullStringValue(template.GetTags()),
			Description: getNullStringValue(template.GetDescription()),
			ImageUrl:    getNullStringValue(template.GetImageUrl()),
			EditorType:  getNullStringValue(template.GetEditorType()),
			// Handle nullable boolean fields
			IsEditable:        getNullBoolValue(template.GetIsEditable()),
			IsPublished:       getNullBoolValue(template.GetIsPublished()),
			IsPublicTemplate:  getNullBoolValue(template.GetIsPublicTemplate()),
			IsGalleryTemplate: getNullBoolValue(template.GetIsGalleryTemplate()),
			IsActive:          getNullBoolValue(template.GetIsActive()),
			// Handle nullable JSONB field
			EmailDesign: getJSONValue(template.GetEmailDesign()),
			// Handle nullable time fields
			CreatedAt: getNullTimeValue(template.GetCreatedAt()),
			UpdatedAt: getNullTimeValue(template.GetUpdatedAt()),
			DeletedAt: getNullTimeValue(template.GetDeletedAt()),
			// Map user information
			User: dto.TemplateUserResponse{
				UserFullname: getNullStringValue(template.GetUserFullname()),
				UserEmail:    getNullStringValue(template.GetUserEmail()),
				UserPicture:  getNullStringValue(template.GetUserPicture()),
			},
			// Map company information
			Company: dto.TemplateCompanyResponse{
				CompanyName: getNullStringValue(template.GetCompanyName()),
			},
		}

		response = append(response, templateResponse)
	}

	return response
}

// Helper function to extract string value from sql.NullString
func getNullStringValue(nullString sql.NullString) string {
	if nullString.Valid {
		return nullString.String
	}
	return ""
}

// Helper function to extract bool value from sql.NullBool
func getNullBoolValue(nullBool sql.NullBool) bool {
	if nullBool.Valid {
		return nullBool.Bool
	}
	return false
}

// Helper function to extract time value from sql.NullTime
func getNullTimeValue(nullTime sql.NullTime) time.Time {
	if nullTime.Valid {
		return nullTime.Time
	}
	return time.Time{}
}

func getNullUUID(id uuid.NullUUID) uuid.UUID {
	if id.Valid {
		return id.UUID
	}
	return uuid.UUID{}
}

// Helper function to extract JSON value from nullable JSON
// This might need adjustment based on your actual JSONB type
func getJSONValue(nullableJSON interface{}) json.RawMessage {
	// If it's already a json.RawMessage
	if rm, ok := nullableJSON.(json.RawMessage); ok && rm != nil {
		return rm
	}

	// If it's a sql.NullString containing JSON
	if ns, ok := nullableJSON.(sql.NullString); ok && ns.Valid {
		return json.RawMessage(ns.String)
	}

	// If it's a struct with RawMessage field (like in your example)
	if nj, ok := nullableJSON.(struct {
		RawMessage json.RawMessage
		Valid      bool
	}); ok && nj.Valid && nj.RawMessage != nil {
		return nj.RawMessage
	}

	return nil
}

 

func  MapSingleTemplateResponse[T MappableTemplate](t T) dto.AdminTemplateResponse {
	return dto.AdminTemplateResponse{
		ID:                t.GetID(),
		UserID:            getNullUUID(t.GetUserID()),
		CompanyID:         getNullUUID(t.GetCompanyID()),
		TemplateName:      t.GetTemplateName(),
		Type:              t.GetType(),
		SenderName:        getNullStringValue(t.GetSenderName()),
		FromEmail:         getNullStringValue(t.GetFromEmail()),
		Subject:           getNullStringValue(t.GetSubject()),
		EmailHtml:         getNullStringValue(t.GetEmailHtml()),
		Tags:              getNullStringValue(t.GetTags()),
		Description:       getNullStringValue(t.GetDescription()),
		ImageUrl:          getNullStringValue(t.GetImageUrl()),
		EditorType:        getNullStringValue(t.GetEditorType()),
		IsEditable:        getNullBoolValue(t.GetIsEditable()),
		IsPublished:       getNullBoolValue(t.GetIsPublished()),
		IsPublicTemplate:  getNullBoolValue(t.GetIsPublicTemplate()),
		IsGalleryTemplate: getNullBoolValue(t.GetIsGalleryTemplate()),
		IsActive:          getNullBoolValue(t.GetIsActive()),
		EmailDesign:       getJSONValue(t.GetEmailDesign()),
		CreatedAt:         getNullTimeValue(t.GetCreatedAt()),
		UpdatedAt:         getNullTimeValue(t.GetUpdatedAt()),
		DeletedAt:         getNullTimeValue(t.GetDeletedAt()),
		User: dto.TemplateUserResponse{
			UserFullname: getNullStringValue(t.GetUserFullname()),
			UserEmail:    getNullStringValue(t.GetUserEmail()),
			UserPicture:  getNullStringValue(t.GetUserPicture()),
		},
		Company: dto.TemplateCompanyResponse{
			CompanyName: getNullStringValue(t.GetCompanyName()),
		},
	}
}