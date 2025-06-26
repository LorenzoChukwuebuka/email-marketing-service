package cronjobs

import (
	"context"
	db "email-marketing-service/internal/db/sqlc"
	"log"
)

type DeleteScheduledUsers struct {
	ctx   context.Context
	store db.Store
}

func NewDeleteScheduledUsers(ctx context.Context, store db.Store) *DeleteScheduledUsers {
	return &DeleteScheduledUsers{
		ctx:   ctx,
		store: store,
	}
}

func (j *DeleteScheduledUsers) Run() {
	log.Println("Starting delete scheduled users job...")

	// Delete users scheduled for deletion after 30 days
	deletedUsers, err := j.store.DeleteScheduledUsers(j.ctx)
	if err != nil {
		log.Printf("Error deleting scheduled users: %v", err)
		return
	}

	if len(deletedUsers) == 0 {
		log.Println("No users found for scheduled deletion")
		return
	}

	log.Printf("Successfully deleted %d scheduled users:", len(deletedUsers))

	// Log details of deleted users and send notification emails
	for _, user := range deletedUsers {
		var scheduledAtInfo string
		if user.ScheduledDeletionAt.Valid {
			scheduledAtInfo = user.ScheduledDeletionAt.Time.Format("2006-01-02 15:04:05")
		} else {
			scheduledAtInfo = "unknown"
		}

		log.Printf("- User: %s (ID: %s, Email: %s) - Scheduled for deletion at: %s",
			user.Fullname,
			user.ID.String(),
			user.Email,
			scheduledAtInfo,
		)

		// TODO: Send email notification to user about account deletion
		// This is CRITICAL - users should be notified when their account is permanently deleted
		// This is often required by privacy laws (GDPR, CCPA, etc.)
		//
		// Email details from user:
		// - Full Name: user.Fullname
		// - Email: user.Email
		// - Scheduled Date: user.ScheduledDeletionAt.Time
		// - Deletion Date: user.DeletedAt.Time (now)
		//
		// Example email implementation:
		// err := j.sendAccountDeletedEmail(user.Fullname, user.Email, user.ScheduledDeletionAt.Time)
		// if err != nil {
		//     log.Printf("Failed to send deletion email for user %s: %v", user.Email, err)
		// }
		//
		// TODO: Consider sending notification to company admin as well
		// if the user was part of a company account
	}

	log.Println("Delete scheduled users job completed successfully")
}

func (j *DeleteScheduledUsers) Schedule() string {
	return "0 0 2 * * *" // Daily at 2:00 AM
}

// sendAccountDeletedEmail sends an email notification when a user account is permanently deleted
// func (j *DeleteScheduledUsers) sendAccountDeletedEmail(fullName, email string, scheduledAt time.Time) error {
//     emailSubject := "Your account has been permanently deleted"
//
//     emailBody := fmt.Sprintf(`
//         Dear %s,
//
//         Your account has been permanently deleted as requested.
//
//         Account Details:
//         - Email: %s
//         - Deletion requested on: %s
//         - Permanently deleted on: %s
//
//         All your personal data has been removed from our systems in accordance with our privacy policy.
//
//         If you believe this was done in error, please contact our support team immediately.
//         Note that after permanent deletion, account recovery may not be possible.
//
//         Thank you for using our service.
//
//         Best regards,
//         Support Team
//     `, fullName, email, scheduledAt.Format("January 2, 2006"), time.Now().Format("January 2, 2006"))
//
//     // Send email using your email service
//     // return j.emailService.SendEmail(email, emailSubject, emailBody)
//
//     return nil
// }

// sendCompanyAdminNotification sends notification to company admin about user deletion
// func (j *DeleteScheduledUsers) sendCompanyAdminNotification(deletedUser db.User, companyAdmin db.User) error {
//     emailSubject := fmt.Sprintf("User %s has been permanently deleted", deletedUser.Email)
//
//     emailBody := fmt.Sprintf(`
//         Dear %s,
//
//         A user from your company account has been permanently deleted.
//
//         Deleted User Details:
//         - Name: %s
//         - Email: %s
//         - Company: [Company Name]
//         - Deletion Date: %s
//
//         This user requested account deletion and has now been permanently removed from the system.
//
//         Best regards,
//         Support Team
//     `, companyAdmin.Fullname, deletedUser.Fullname, deletedUser.Email, time.Now().Format("January 2, 2006"))
//
//     // Send email using your email service
//     // return j.emailService.SendEmail(companyAdmin.Email, emailSubject, emailBody)
//
//     return nil
// }
