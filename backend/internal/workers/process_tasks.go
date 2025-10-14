package worker

import (
	"context"
	"database/sql"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/enums"
	"encoding/json"
	"fmt"
	"log"
)

func (w *Worker) ProcessTask(ctx context.Context, task db.Task) error {
	switch TaskType(task.TaskType) {
	case TaskSendWelcomeEmail:
		var payload EmailPayload
		if err := json.Unmarshal(task.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}
		return ProcessWelcomeEmailTask(ctx, payload)

	case TaskAuditLogCreate:
		var payload AuditCreatePayload
		if err := json.Unmarshal(task.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}
		return w.ProcessLogAuditCreate(ctx, payload)
	case TaskAuditLogUpdate:
		var payload AuditUpdatePayload
		if err := json.Unmarshal(task.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}
		return w.ProcessLogAuditUpdate(ctx, payload)

	case TaskAuditLogDelete:
		var payload AuditDeletePayload
		if err := json.Unmarshal(task.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}
		return w.ProcessLogAuditDelete(ctx, payload)

	case TaskAuditLogLogin:
		var payload AuditLoginPayload
		if err := json.Unmarshal(task.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}
		return w.ProcessLogAuditLogin(ctx, payload)

	case TaskAuditLogFailedLogin:
		var payload AuditFailedLoginPayload
		if err := json.Unmarshal(task.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}
		return w.ProcessLogAuditFailedLogin(ctx, payload)

	case TaskSendCampaignEmail:
		var payload SendCampaignEmailPayload

		if err := json.Unmarshal(task.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload:%w", err)
		}

		err := w.ProcessSendCampaignEmailsTask(ctx, payload)

		if err != nil {
			// Log the detailed error
			log.Printf("Campaign processing failed with detailed error: %v", err)

			// Update campaign status to failed
			if updateErr := w.Store.UpdateCampaignStatus(ctx, db.UpdateCampaignStatusParams{
				Status: sql.NullString{String: string(enums.CampaignStatus(enums.Failed)), Valid: true},
				ID:     payload.CampaignID,
				UserID: payload.UserID,
			}); updateErr != nil {
				log.Printf("error occurred while updating status to failed: %v", updateErr)
			}

			if campaignErr, ok := err.(CampaignError); ok {
				w.StoreCampaignError(ctx, payload.CampaignID, campaignErr)
			} else {
				// If it's not a CampaignError, create one
				w.StoreCampaignError(ctx, payload.CampaignID, CampaignError{
					Type:    "UNKNOWN_ERROR",
					Message: "Campaign processing failed",
					Err:     err,
				})
			}
		}

		return err
	case TaskSendAdminNotification:
		var payload AdminNotificationPayload
		if err := json.Unmarshal(task.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}
		return w.ProcessAdminNotificationTask(ctx, payload)
	case TaskSendUserNotification:
		var payload UserNotificationPayload
		if err := json.Unmarshal(task.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}
		return w.ProcessUserNotificationTask(ctx, payload)

	default:
		return fmt.Errorf("unknown task type: %s", task.TaskType)
	}
}
