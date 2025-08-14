package cronjobs

import (
	"context"
	"database/sql"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/enums"
	"log"
)

type UpdateExpiredSubscriptionJob struct {
	*BaseJob
}

func NewUpdateExpiredSubscriptionJob(store db.Store, ctx context.Context) *UpdateExpiredSubscriptionJob {
	baseJob := NewBaseJob(
		store,
		ctx,
		"auto_update_expired_subs",
		"AutoUpdateExpiredSubscriptions",
		"Automatically update expired subscriptions",
	)

	return &UpdateExpiredSubscriptionJob{
		BaseJob: baseJob,
	}
}

func (j *UpdateExpiredSubscriptionJob) Run()error {
	//get the expired subscriptions still active
	expired_subs, err := j.store.GetExpiredActiveSubscriptions(j.ctx)
	if err != nil {
		log.Fatalf("error fetching data:%v", err)
	}
	// Check if there are any expired subscriptions
	if len(expired_subs) == 0 {
		log.Println("No expired subscriptions found")
		return err
	}

	log.Printf("Found %d expired subscriptions to update", len(expired_subs))

	// Loop over the expired subscriptions
	for i, subscription := range expired_subs {
		log.Printf("Processing subscription %d/%d: ID=%s", i+1, len(expired_subs), subscription.ID)

		// Update the subscription status to expired
		_, err = j.store.UpdateSubscriptionStatus(j.ctx, db.UpdateSubscriptionStatusParams{
			ID:     subscription.ID,
			Status: sql.NullString{String: string(enums.SubscriptionExpired), Valid: true},
		})
		if err != nil {
			log.Printf("Error updating subscription %s: %v", subscription.ID, err)
			continue // Continue with next subscription instead of failing completely
		}

		log.Printf("Successfully updated subscription %s to inactive", subscription.ID)
	}

	log.Printf("Completed processing %d expired subscriptions", len(expired_subs))

	return  nil

}

func (j *UpdateExpiredSubscriptionJob) Schedule() string {
	return "0 0 0 * * *"
}
