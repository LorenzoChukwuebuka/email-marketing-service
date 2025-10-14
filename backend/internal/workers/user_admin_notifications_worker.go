package worker

import (
	"context"
	"database/sql"
	db "email-marketing-service/internal/db/sqlc"
	"fmt"
	"log/slog"
)

func (w *Worker) ProcessUserNotificationTask(ctx context.Context, payload UserNotificationPayload) error {
	slog.Info("sending out notifications to user")
	var AdditionalField string

	if payload.AdditionalField != nil {
		AdditionalField = *payload.AdditionalField
	}

	_, err := w.Store.CreateUserNotification(ctx, db.CreateUserNotificationParams{
		UserID:          payload.UserId,
		Title:           payload.NotifcationTitle,
		AdditionalField: AdditionalField,
	})
	if err != nil {
		return fmt.Errorf("failed to create user notification: %w", err)
	}

	return nil
}

func (w *Worker) ProcessAdminNotificationTask(ctx context.Context, payload AdminNotificationPayload) error {
	slog.Info("sending out notifications to admin")
	_, err := w.Store.CreateAdminNotification(ctx, db.CreateAdminNotificationParams{
		UserID: payload.UserId,
		Title:  payload.NotificationTitle,
		Link:   sql.NullString{String: payload.Link, Valid: payload.Link != ""},
	})
	if err != nil {
		return fmt.Errorf("failed to create admin notification: %w", err)
	}
	return nil
}
