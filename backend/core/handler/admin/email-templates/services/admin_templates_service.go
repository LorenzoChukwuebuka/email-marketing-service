package services

import (
	"context"
	"database/sql"
	"email-marketing-service/core/handler/admin/email-templates/dto"
	"email-marketing-service/internal/common"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/helper"
	"errors"
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

func (s *AdminTemplatesService) GetTemplate(ctx context.Context, req *dto.AdminFetchGalleryTemplatesDTO) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"template": req.TemplateId})
	if err != nil {
		return nil, err
	}
	template, err := s.store.GetTemplateByID(ctx, db.GetTemplateByIDParams{
		TemplateID: _uuid["template"],
	})
	if err != nil {
		return nil, err
	}
	return template, nil
}
func (s *AdminTemplatesService) GetTemplates() {}

func (s *AdminTemplatesService) UpdateTemplate() {}
func (s *AdminTemplatesService) DeleteTemplate() {}

func (s *AdminTemplatesService) ArchiveTemplate() {}
