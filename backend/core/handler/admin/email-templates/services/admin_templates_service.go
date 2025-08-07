package services

import (
	"context"
	"database/sql"
	"email-marketing-service/core/handler/admin/email-templates/dto"
	"email-marketing-service/core/handler/admin/email-templates/mapper"
	"email-marketing-service/internal/common"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/helper"
	"errors"
	"fmt"
	"github.com/sqlc-dev/pqtype"
)

type AdminTemplatesService struct {
	store db.Store
}

func NewAdminTemplatesService(store db.Store) *AdminTemplatesService {
	return &AdminTemplatesService{
		store: store,
	}
}

func (s *AdminTemplatesService) CreateTemplate(ctx context.Context, req *dto.AdminTemplateDTO) (any, error) {
	if err := helper.ValidateData(req); err != nil {
		return nil, errors.Join(common.ErrValidatingRequest, err)
	}

	templateExists, err := s.store.CheckTemplateExists(ctx, req.TemplateName)
	if err != nil {
		return nil, err
	}

	if templateExists {
		return nil, common.ErrRecordExists
	}

	template, err := s.store.CreateTemplate(ctx, db.CreateTemplateParams{
		TemplateName: req.TemplateName,
		SenderName:   sql.NullString{String: req.SenderName, Valid: true},
		FromEmail:    sql.NullString{String: req.FromEmail, Valid: true},
		Subject:      sql.NullString{String: req.Subject, Valid: true},
		Type:         req.Type,
		EmailHtml:    sql.NullString{String: req.EmailHtml, Valid: true},
		EmailDesign:  pqtype.NullRawMessage{RawMessage: req.EmailDesign, Valid: true},
		IsEditable:   sql.NullBool{Bool: req.IsEditable, Valid: true},
		IsPublished:  sql.NullBool{Bool: req.IsPublished, Valid: true},
		IsPublicTemplate: sql.NullBool{
			Bool:  req.IsPublicTemplate,
			Valid: true,
		},
		IsGalleryTemplate: sql.NullBool{
			Bool:  req.IsGalleryTemplate,
			Valid: true,
		},
		Tags:        sql.NullString{String: req.Tags, Valid: true},
		Description: sql.NullString{String: req.Description, Valid: true},
		ImageUrl:    sql.NullString{String: req.ImageUrl, Valid: true},
		IsActive:    sql.NullBool{Bool: req.IsActive, Valid: true},
		EditorType:  sql.NullString{String: req.EditorType, Valid: true},
	})

	if err != nil {
		return nil, common.ErrCreatingRecord
	}

	return template, nil
}

func (s *AdminTemplatesService) GetTemplate(ctx context.Context, templateId string) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"template": templateId})
	if err != nil {
		return nil, err
	}
	template, err := s.store.GetTemplateByIDGallery(ctx, _uuid["template"])
	if err != nil {
		return nil, err
	}

	response := mapper.MapSingleTemplateResponse(mapper.TemplateAdapter(template))
	return response, nil
}

func (s *AdminTemplatesService) GetTemplatesByType(ctx context.Context, req *dto.AdminFetchGalleryTemplatesDTO) (any, error) {

	templates, err := s.store.ListTemplatesByTypeGallery(ctx, db.ListTemplatesByTypeGalleryParams{
		TemplateType:   req.Type,
		TemplateSearch: req.Search,
		RowOffset:      int32(req.Offset),
		RowLimit:       int32(req.Limit),
	})

	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		return nil, common.ErrFetchingRecord
	}

	count_templates, err := s.store.CountGalleryTemplates(ctx)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		return nil, common.ErrFetchingRecord
	}

	var adaptedTemplates []mapper.TemplateAdapter
	for _, t := range templates {
		adaptedTemplates = append(adaptedTemplates, mapper.TemplateAdapter(t))
	}
	response := mapper.MapTemplateResponse(adaptedTemplates)

	items := make([]any, len(response))
	for i, v := range response {
		items[i] = v
	}

	data := common.Paginate(int(count_templates), items, req.Offset, req.Limit)
	return data, nil
}

func (s *AdminTemplatesService) UpdateTemplate() {}

func (s *AdminTemplatesService) DeleteTemplate(ctx context.Context, templateId string) error {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"template": templateId})
	if err != nil {
		return err
	}

	if err := s.store.HardDeleteTemplate(ctx, _uuid["template"]); err != nil {
		return common.ErrDeletingRecord
	}

	return nil
}

func (s *AdminTemplatesService) ArchiveOrUnArchiveTemplate() {}

func (s *AdminTemplatesService) PublishOrUnpublishTemplate() {}
