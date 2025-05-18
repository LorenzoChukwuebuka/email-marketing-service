package adminrepository

import (
	"email-marketing-service/api/v1/model"
	"errors"
	"gorm.io/gorm"
	"time"
)

type AdminTemplateRepository struct {
	DB *gorm.DB
}

// ConvertTemplateToTemplateResponse converts a Template model into a TemplateResponse struct
func ConvertTemplateToTemplateResponse(template model.Template) model.TemplateResponse {
	deletedAt := ""
	if template.DeletedAt.Valid {
		deletedAt = template.DeletedAt.Time.Format(time.RFC3339)
	}
	return model.TemplateResponse{
		UUID:              template.UUID,
		CreatedAt:         template.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         template.UpdatedAt.Format(time.RFC3339),
		DeletedAt:         &deletedAt,
		UserId:            template.UserId,
		TemplateName:      template.TemplateName,
		SenderName:        template.SenderName,
		FromEmail:         template.FromEmail,
		Subject:           template.Subject,
		Type:              template.Type,
		EmailHtml:         template.EmailHtml,
		EmailDesign:       template.EmailDesign,
		IsEditable:        template.IsEditable,
		IsPublished:       template.IsPublished,
		IsPublicTemplate:  template.IsPublicTemplate,
		IsGalleryTemplate: template.IsGalleryTemplate,
		Tags:              template.Tags,
		Description:       template.Description,
		ImageUrl:          template.ImageUrl,
		IsActive:          template.IsActive,
		EditorType:        template.EditorType,
	}
}

// CreateTemplate adds a new template to the database
func (r *AdminTemplateRepository) CreateTemplate(template *model.Template) (*model.TemplateResponse, error) {
	result := r.DB.Create(template)
	if result.Error != nil {
		return nil, result.Error
	}
	templateResponse := ConvertTemplateToTemplateResponse(*template)
	return &templateResponse, nil
}

// GetAllAdminCreatedTemplates retrieves all templates created by admin users
func (r *AdminTemplateRepository) GetAllAdminCreatedTemplates() ([]model.TemplateResponse, error) {
	var templates []model.Template
	result := r.DB.Where("is_public_template = ?", true).Find(&templates)
	if result.Error != nil {
		return nil, result.Error
	}

	var templateResponses []model.TemplateResponse
	for _, template := range templates {
		templateResponses = append(templateResponses, ConvertTemplateToTemplateResponse(template))
	}
	return templateResponses, nil
}

// DeleteAdminCreatedTemplate deletes a specific template created by an admin
func (r *AdminTemplateRepository) DeleteAdminCreatedTemplate(templateId string) error {
	result := r.DB.Where("uuid = ?", templateId).Delete(&model.Template{})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("template not found")
		}
		return result.Error
	}
	return nil
}

// UpdateAdminCreatedTemplate updates an existing template created by an admin
func (r *AdminTemplateRepository) UpdateAdminCreatedTemplate(templateId string, updatedTemplate model.Template) (*model.TemplateResponse, error) {
	var template model.Template
	result := r.DB.Where("uuid = ? AND is_public_template = ?", templateId, true).First(&template)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("template not found")
		}
		return nil, result.Error
	}

	// Update template fields with nil pointer checks
	if updatedTemplate.TemplateName != "" {
		template.TemplateName = updatedTemplate.TemplateName
	}

	if updatedTemplate.SenderName != nil {
		template.SenderName = updatedTemplate.SenderName
	}

	if updatedTemplate.FromEmail != nil {
		template.FromEmail = updatedTemplate.FromEmail
	}

	if updatedTemplate.Subject != nil {
		template.Subject = updatedTemplate.Subject
	}

	if updatedTemplate.EmailHtml != "" {
		template.EmailHtml = updatedTemplate.EmailHtml
	}

	template.IsEditable = updatedTemplate.IsEditable
	template.IsPublished = updatedTemplate.IsPublished
	template.IsPublicTemplate = updatedTemplate.IsPublicTemplate

	if updatedTemplate.Tags != "" {
		template.Tags = updatedTemplate.Tags
	}

	if updatedTemplate.Description != nil {
		template.Description = updatedTemplate.Description
	}

	if updatedTemplate.ImageUrl != nil {
		template.ImageUrl = updatedTemplate.ImageUrl
	}

	template.IsActive = updatedTemplate.IsActive

	if updatedTemplate.EditorType != nil {
		template.EditorType = updatedTemplate.EditorType
	}

	// Save changes
	if err := r.DB.Save(&template).Error; err != nil {
		return nil, err
	}

	templateResponse := ConvertTemplateToTemplateResponse(template)
	return &templateResponse, nil
}

// GetSingleAdminCreatedTemplate retrieves a single template by UUID
func (r *AdminTemplateRepository) GetSingleAdminCreatedTemplate(templateId string) (*model.TemplateResponse, error) {
	var template model.Template
	result := r.DB.Where("uuid = ?", templateId).First(&template)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("template not found")
		}
		return nil, result.Error
	}

	templateResponse := ConvertTemplateToTemplateResponse(template)
	return &templateResponse, nil
}

// GetAllTemplatesCreatedByAUser retrieves all templates created by a specific user
func (r *AdminTemplateRepository) GetAllTemplatesCreatedByAUser(userId string) ([]model.TemplateResponse, error) {
	var templates []model.Template
	result := r.DB.Where("user_id = ?", userId).Find(&templates)
	if result.Error != nil {
		return nil, result.Error
	}

	var templateResponses []model.TemplateResponse
	for _, template := range templates {
		templateResponses = append(templateResponses, ConvertTemplateToTemplateResponse(template))
	}
	return templateResponses, nil
}

// GetASpecificUserTemplate retrieves a specific template created by a user
func (r *AdminTemplateRepository) GetASpecificUserTemplate(userId string, templateId string) (*model.TemplateResponse, error) {
	var template model.Template
	result := r.DB.Where("uuid = ? AND user_id = ?", templateId, userId).First(&template)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("template not found")
		}
		return nil, result.Error
	}

	templateResponse := ConvertTemplateToTemplateResponse(template)
	return &templateResponse, nil
}
