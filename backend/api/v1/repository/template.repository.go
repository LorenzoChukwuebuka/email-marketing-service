package repository

import (
	"email-marketing-service/api/v1/model"
	"fmt"
	"gorm.io/gorm"
)

type TemplateRepository struct {
	DB *gorm.DB
}

func NewTemplateRepository(db *gorm.DB) *TemplateRepository {
	return &TemplateRepository{DB: db}
}

func (r *TemplateRepository) CreateAndUpdateTemplate(d *model.Template) error {
	if err := r.DB.Create(&d).Error; err != nil {
		return fmt.Errorf("failed to save template : %w", err)
	}
	return nil
}

func (r *TemplateRepository) CheckMarketingNameExists(d *model.Template) (bool, error) {
	result := r.DB.Where("template_name = ? AND user_id = ?", d.TemplateName, d.UserId).First(&d)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}

func (r *TemplateRepository) GetTransactionalTemplate(userId string, templateId string) (*model.TemplateResponse, error) {
	var template model.Template
	result := r.DB.Where("type = ? AND user_id = ? AND uuid = ?", model.Transactional, userId, templateId).First(&template)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	response := convertToTemplateResponse(&template)
	return response, nil
}

func (r *TemplateRepository) GetMarketingTemplate(userId string, templateId string) (*model.TemplateResponse, error) {
	var template model.Template
	result := r.DB.Where("type = ? AND user_id = ? AND uuid = ?", model.Marketing, userId, templateId).First(&template)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	response := convertToTemplateResponse(&template)
	return response, nil
}

func (r *TemplateRepository) UpdateTemplate(d *model.Template) error {
	return r.DB.Model(&model.Template{}).Where("uuid = ?", d.UUID).Updates(d).Error
}

func (r *TemplateRepository) DeleteTemplate(d *model.Template) error {
	if err := r.DB.Delete(d).Error; err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}
	return nil
}

func (r *TemplateRepository) GetAllTransactionalTemplates(userId string) ([]model.TemplateResponse, error) {
	var templates []model.Template
	if err := r.DB.Where("type = ? AND user_id = ?", model.Transactional, userId).Order("created_at DESC").Find(&templates).Error; err != nil {
		return nil, fmt.Errorf("failed to get transactional templates: %w", err)
	}

	var templateResponses []model.TemplateResponse
	for _, template := range templates {
		templateResponses = append(templateResponses, *convertToTemplateResponse(&template))
	}

	return templateResponses, nil
}

func (r *TemplateRepository) GetAllMarketingTemplates(userId string) ([]model.TemplateResponse, error) {
	var templates []model.Template
	if err := r.DB.Where("type = ? AND user_id = ?", model.Marketing, userId).Order("created_at DESC").Find(&templates).Error; err != nil {
		return nil, fmt.Errorf("failed to get marketing templates: %w", err)
	}

	var templateResponses []model.TemplateResponse
	for _, template := range templates {
		templateResponses = append(templateResponses, *convertToTemplateResponse(&template))
	}

	return templateResponses, nil
}

func (r *TemplateRepository) GetSingleTemplate(templateId string) (*model.TemplateResponse, error) {
	var template model.Template
	result := r.DB.Where("  uuid = ?", templateId).First(&template)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	response := convertToTemplateResponse(&template)
	return response, nil
}

func convertToTemplateResponse(t *model.Template) *model.TemplateResponse {
	var deletedAt *string
	if t.DeletedAt.Valid {
		deletedAtStr := t.DeletedAt.Time.Format("2006-01-02 15:04:05")
		deletedAt = &deletedAtStr
	}

	return &model.TemplateResponse{
		ID:                t.ID,
		UUID:              t.UUID,
		CreatedAt:         t.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         t.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt:         deletedAt,
		UserId:            t.UserId,
		TemplateName:      t.TemplateName,
		SenderName:        t.SenderName,
		FromEmail:         t.FromEmail,
		Subject:           t.Subject,
		Type:              t.Type,
		EmailHtml:         t.EmailHtml,
		EmailDesign:       t.EmailDesign,
		IsEditable:        t.IsEditable,
		IsPublished:       t.IsPublished,
		IsPublicTemplate:  t.IsPublicTemplate,
		IsGalleryTemplate: t.IsGalleryTemplate,
		Tags:              t.Tags,
		Description:       t.Description,
		ImageUrl:          t.ImageUrl,
		IsActive:          t.IsActive,
		EditorType:        t.EditorType,
	}
}
