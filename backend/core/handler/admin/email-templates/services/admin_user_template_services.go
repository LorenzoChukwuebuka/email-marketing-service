package services

import (
	"context"
	"email-marketing-service/core/handler/admin/email-templates/dto"
	"email-marketing-service/core/handler/admin/email-templates/mapper"
	"email-marketing-service/internal/common"
	db "email-marketing-service/internal/db/sqlc"
	"fmt"
	"github.com/google/uuid"
)

type AdminUserTemplatesService struct {
	store db.Store
}

func NewAdminUserTemplatesService(store db.Store) *AdminUserTemplatesService {
	return &AdminUserTemplatesService{
		store: store,
	}
}

func (s *AdminUserTemplatesService) GetUserTemplates(ctx context.Context, req *dto.AdminFetchTemplateDTO) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"user": req.UserID,
	})
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	templates, err := s.store.ListTemplatesByType(ctx, db.ListTemplatesByTypeParams{
		Type:    req.Type,
		UserID:  uuid.NullUUID{UUID: _uuid["user"], Valid: true},
		Limit:   int32(req.Limit),
		Offset:  int32(req.Offset),
		Column5: req.SearchQuery,
	})
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		return nil, common.ErrFetchingRecord
	}

	count_templates, err := s.store.CountTemplatesByUserID(ctx, uuid.NullUUID{UUID: _uuid["user"], Valid: true})
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		return nil, common.ErrFetchingRecord
	}

	var adaptedTemplates []mapper.ListTemplatesByTypeRowAdapter
	for _, t := range templates {
		adaptedTemplates = append(adaptedTemplates, mapper.ListTemplatesByTypeRowAdapter(t))
	}
	response := mapper.MapTemplateResponse(adaptedTemplates)

	items := make([]any, len(response))
	for i, v := range response {
		items[i] = v
	}

	data := common.Paginate(int(count_templates), items, req.Offset, req.Limit)
	return data, nil

}

func (s *AdminUserTemplatesService) GetUserTemplateById(ctx context.Context, req *dto.AdminFetchTemplateDTO) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"user":     req.UserID,
		"template": req.TemplateId,
	})
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	template, err := s.store.GetTemplateByID(ctx, db.GetTemplateByIDParams{
		UserID:     uuid.NullUUID{UUID: _uuid["user"], Valid: true},
		TemplateID: _uuid["template"],
		Type:       req.Type,
	})
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		return nil, common.ErrFetchingRecord
	}

	response := mapper.MapSingleTemplateResponse(mapper.TemplateByIDRowAdapter(template))
	return response, nil
}

func (s *AdminUserTemplatesService) SuspendTemplate(userId int64, templateId int64) error {
	// Logic to suspend a template for the user
	return nil
}

func (s *AdminUserTemplatesService) ActivateTemplate(userId int64, templateId int64) error {
	// Logic to activate a template for the user
	return nil
}
