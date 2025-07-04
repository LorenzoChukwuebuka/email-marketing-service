package mapper

import (
	"database/sql"
	"email-marketing-service/core/handler/templates/dto"
	db "email-marketing-service/internal/db/sqlc"
	"encoding/json"
	"time"
)


func MapTemplateToDTO(dbTemplate db.Template) any {
	return &dto.TemplateDTO{
		TemplateID:        dbTemplate.ID.String(),
		UserId:           dbTemplate.UserID.String(),
		CompanyID:        dbTemplate.CompanyID.String(),
		TemplateName:     dbTemplate.TemplateName,
		SenderName:       dbTemplate.SenderName.String,
		FromEmail:        dbTemplate.FromEmail.String,
		Subject:          dbTemplate.Subject.String,
		Type:             dbTemplate.Type,
		EmailHtml:        dbTemplate.EmailHtml.String,
		EmailDesign:      dbTemplate.EmailDesign.RawMessage,
		IsEditable:       dbTemplate.IsEditable.Bool,
		IsPublished:      dbTemplate.IsPublished.Bool,
		IsPublicTemplate: dbTemplate.IsPublicTemplate.Bool,
		IsGalleryTemplate: dbTemplate.IsGalleryTemplate.Bool,
		Tags:             dbTemplate.Tags.String,
		Description:      dbTemplate.Description.String,
		ImageUrl:         dbTemplate.ImageUrl.String,
		IsActive:         dbTemplate.IsActive.Bool,
		EditorType:       dbTemplate.EditorType.String,
	}
}


func MapTemplateResponse(templates []db.ListTemplatesByTypeRow) []dto.TemplateResponse {
	if len(templates) == 0 {
		return []dto.TemplateResponse{}
	}

	var response []dto.TemplateResponse

	for _, template := range templates {
		templateResponse := dto.TemplateResponse{
			ID:           template.ID,
			UserID:       template.UserID,
			CompanyID:    template.CompanyID,
			TemplateName: template.TemplateName,
			Type:         template.Type,
			// Handle nullable string fields
			SenderName:  getNullStringValue(template.SenderName),
			FromEmail:   getNullStringValue(template.FromEmail),
			Subject:     getNullStringValue(template.Subject),
			EmailHtml:   getNullStringValue(template.EmailHtml),
			Tags:        getNullStringValue(template.Tags),
			Description: getNullStringValue(template.Description),
			ImageUrl:    getNullStringValue(template.ImageUrl),
			EditorType:  getNullStringValue(template.EditorType),
			// Handle nullable boolean fields
			IsEditable:        getNullBoolValue(template.IsEditable),
			IsPublished:       getNullBoolValue(template.IsPublished),
			IsPublicTemplate:  getNullBoolValue(template.IsPublicTemplate),
			IsGalleryTemplate: getNullBoolValue(template.IsGalleryTemplate),
			IsActive:          getNullBoolValue(template.IsActive),
			// Handle nullable JSONB field
			EmailDesign: getJSONValue(template.EmailDesign),
			// Handle nullable time fields
			CreatedAt: getNullTimeValue(template.CreatedAt),
			UpdatedAt: getNullTimeValue(template.UpdatedAt),
			DeletedAt: getNullTimeValue(template.DeletedAt),
			// Map user information
			User: dto.TemplateUserResponse{
				UserFullname: getNullStringValue(template.UserFullname),
				UserEmail:    getNullStringValue(template.UserEmail),
				UserPicture:  getNullStringValue(template.UserPicture),
			},
			// Map company information
			Company: dto.TemplateCompanyResponse{
				CompanyName: getNullStringValue(template.CompanyName),
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

func MapSingleTemplateResponse(template db.GetTemplateByIDRow) dto.TemplateResponse {
	templateResponse := dto.TemplateResponse{
		ID:           template.ID,
		UserID:       template.UserID,
		CompanyID:    template.CompanyID,
		TemplateName: template.TemplateName,
		Type:         template.Type,
		// Handle nullable string fields
		SenderName:  getNullStringValue(template.SenderName),
		FromEmail:   getNullStringValue(template.FromEmail),
		Subject:     getNullStringValue(template.Subject),
		EmailHtml:   getNullStringValue(template.EmailHtml),
		Tags:        getNullStringValue(template.Tags),
		Description: getNullStringValue(template.Description),
		ImageUrl:    getNullStringValue(template.ImageUrl),
		EditorType:  getNullStringValue(template.EditorType),
		// Handle nullable boolean fields
		IsEditable:        getNullBoolValue(template.IsEditable),
		IsPublished:       getNullBoolValue(template.IsPublished),
		IsPublicTemplate:  getNullBoolValue(template.IsPublicTemplate),
		IsGalleryTemplate: getNullBoolValue(template.IsGalleryTemplate),
		IsActive:          getNullBoolValue(template.IsActive),
		// Handle nullable JSONB field
		EmailDesign: getJSONValue(template.EmailDesign),
		// Handle nullable time fields
		CreatedAt: getNullTimeValue(template.CreatedAt),
		UpdatedAt: getNullTimeValue(template.UpdatedAt),
		DeletedAt: getNullTimeValue(template.DeletedAt),
		// Map user information
		User: dto.TemplateUserResponse{
			UserFullname: getNullStringValue(template.UserFullname),
			UserEmail:    getNullStringValue(template.UserEmail),
			UserPicture:  getNullStringValue(template.UserPicture),
		},
		// Map company information
		Company: dto.TemplateCompanyResponse{
			CompanyName: getNullStringValue(template.CompanyName),
		},
	}

	return templateResponse
}