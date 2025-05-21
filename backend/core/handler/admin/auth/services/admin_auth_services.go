package services

import (
	"context"
	"database/sql"
	"email-marketing-service/core/handler/admin/auth/dto"
	"email-marketing-service/core/handler/admin/auth/mapper"
	"email-marketing-service/internal/common"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/helper"
	"errors"
	"fmt"
	"strings"
)

type Service struct {
	store db.Store
}

func NewAdminAuthService(store db.Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) CreateAdmin(ctx context.Context, req *dto.AdminRequestDTO) (*dto.AdminRequestDTO, error) {
	if err := helper.ValidateData(req); err != nil {
		return nil, errors.Join(common.ErrValidatingRequest, err)
	}
	_, err := s.store.GetAdminByEmail(ctx, req.Email)
	if !errors.Is(err, sql.ErrNoRows) && !strings.Contains(err.Error(), "no rows") {
		return nil, fmt.Errorf("error checking domain existence: %w", err)
	}

	if err == nil {
		return nil, fmt.Errorf("admin already exists")
	}

	hashedPass, err := common.HashPassword(req.Password)

	if err != nil {
		return nil, common.ErrPasswordHashingFailed
	}

	_, err = s.store.CreateAdmin(ctx, db.CreateAdminParams{
		Firstname:  sql.NullString{String: req.FirstName, Valid: true},
		Middlename: sql.NullString{String: req.MiddleName, Valid: true},
		Lastname:   sql.NullString{String: req.LastName, Valid: true},
		Email:      req.Email,
		Password:   hashedPass,
		Type:       req.Type,
	})

	return req, nil
}

func (s *Service) AdminLogin(ctx context.Context, req *dto.AdminLoginRequest) (*dto.AdminLoginResponse[dto.AdminResponse], error) {
	if err := helper.ValidateData(req); err != nil {
		return nil, errors.Join(common.ErrValidatingRequest, err)
	}

	admin, err := s.store.GetAdminByEmail(ctx, req.Email)
	if err != nil {
		return nil, common.ErrFetchingAdmin
	}

	err = common.CheckPassword(req.Password, admin.Password)
	if err != nil {
		return nil, common.ErrCheckingPasswordHash
	}

	token, err := helper.GenerateAdminAccessToken(admin.ID.String(), admin.ID.String(), admin.Type, admin.Email)

	if err != nil {
		return nil, err
	}

	refreshToken, err := helper.GenerateAdminRefreshToken(admin.ID.String(), admin.ID.String(), admin.Type, admin.Email)
	if err != nil {
		return nil, err
	}

	return &dto.AdminLoginResponse[dto.AdminResponse]{
		Status:       "login successful",
		Token:        token,
		Details:      mapper.MapAdminToResponse(admin),
		RefreshToken: refreshToken,
		Type:         admin.Type,
	}, nil

}
