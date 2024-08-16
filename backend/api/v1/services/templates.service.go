package services

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/utils"
	"fmt"
	"github.com/google/uuid"
)

type TemplateService struct {
	TemplateRepo *repository.TemplateRepository
}

func NewTemplateService(templateRepo *repository.TemplateRepository) *TemplateService {
	return &TemplateService{
		TemplateRepo: templateRepo,
	}
}

func (s *TemplateService) CreateTemplate(d *dto.TemplateDTO) (map[string]interface{}, error) {

	if err := utils.ValidateData(d); err != nil {
		return nil, fmt.Errorf("invalid data: %w", err)
	}

	id := uuid.New().String()

	templateModel := &model.Template{
		UUID:              id,
		UserId:            d.UserId,
		TemplateName:      d.TemplateName,
		SenderName:        d.SenderName,
		FromEmail:         d.FromEmail,
		Subject:           d.Subject,
		Type:              model.TemplateType(d.Type),
		EmailHtml:         d.EmailHtml,
		EmailDesign:       d.EmailDesign,
		IsEditable:        d.IsEditable,
		IsPublished:       d.IsPublished,
		IsPublicTemplate:  d.IsPublicTemplate,
		IsGalleryTemplate: d.IsGalleryTemplate,
		Tags:              d.Tags,
		Description:       d.Description,
		ImageUrl:          d.ImageUrl,
		IsActive:          d.IsActive,
		EditorType:        d.EditorType,
	}

	checkIfTempleExists, err := s.TemplateRepo.CheckMarketingNameExists(templateModel)

	if err != nil {
		return nil, err
	}

	if checkIfTempleExists {
		return nil, fmt.Errorf("template name already exists")
	}

	if err := s.TemplateRepo.CreateAndUpdateTemplate(templateModel); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"templateId":  id,
		"editor-type": d.EditorType,
		"type":        d.Type,
		"message":     "template created successfully",
	}, nil
}

func (s *TemplateService) GetTransactionalTemplate(userId string, templateId string) (*model.TemplateResponse, error) {
	result, err := s.TemplateRepo.GetTransactionalTemplate(userId, templateId)

	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	return result, nil
}

func (s *TemplateService) GetMarketingTemplate(userId string, templateId string) (*model.TemplateResponse, error) {
	result, err := s.TemplateRepo.GetMarketingTemplate(userId, templateId)

	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	return result, nil
}

func (s *TemplateService) GetAllTransactionalTemplates(userId string) ([]model.TemplateResponse, error) {

	result, err := s.TemplateRepo.GetAllTransactionalTemplates(userId)

	if err != nil {
		return []model.TemplateResponse{}, err
	}

	if len(result) < 1 {
		return []model.TemplateResponse{}, nil
	}

	return result, nil
}

func (s *TemplateService) GetAllMarketingTemplates(userId string) ([]model.TemplateResponse, error) {
	result, err := s.TemplateRepo.GetAllMarketingTemplates(userId)

	if err != nil {
		return []model.TemplateResponse{}, err
	}

	if len(result) < 1 {
		return []model.TemplateResponse{}, nil
	}

	return result, nil
}

func (s *TemplateService) UpdateTemplate(d *dto.TemplateDTO, templateId string) error {

	templateModel := &model.Template{
		UUID:              templateId,
		UserId:            d.UserId,
		TemplateName:      d.TemplateName,
		SenderName:        d.SenderName,
		FromEmail:         d.FromEmail,
		Subject:           d.Subject,
		Type:              model.TemplateType(d.Type),
		EmailHtml:         d.EmailHtml,
		EmailDesign:       d.EmailDesign,
		IsEditable:        d.IsEditable,
		IsPublished:       d.IsPublished,
		IsPublicTemplate:  d.IsPublicTemplate,
		IsGalleryTemplate: d.IsGalleryTemplate,
		Tags:              d.Tags,
		Description:       d.Description,
		ImageUrl:          d.ImageUrl,
		IsActive:          d.IsActive,
		EditorType:        d.EditorType,
	}

	if err := s.TemplateRepo.UpdateTemplate(templateModel); err != nil {
		return err
	}
	return nil
}

func (s *TemplateService) DeleteTemplate(userId string, templateId string) error {
	return nil
}


func (s *TemplateService) SendTestMail(userId string, templateId string) error {
	return nil 
}
