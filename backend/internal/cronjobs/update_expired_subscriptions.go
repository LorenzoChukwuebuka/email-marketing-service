package cronjobs

import (
	"context"
	"database/sql"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/enums"
	"log"
)

type UpdateExpiredSubscriptionJob struct {
	store db.Store
	ctx   context.Context
}

func NewUpdateExpiredSubscriptionJob(store db.Store, ctx context.Context) *UpdateExpiredSubscriptionJob {
	return &UpdateExpiredSubscriptionJob{
		store: store,
		ctx:   ctx,
	}
}

func (j *UpdateExpiredSubscriptionJob) Run() {
	//get the expired subscriptions still active
	expired_subs, err := j.store.GetExpiredActiveSubscriptions(j.ctx)
	if err != nil {
		log.Fatalf("error fetching data:%v", err)
	}
	// Check if there are any expired subscriptions
	if len(expired_subs) == 0 {
		log.Println("No expired subscriptions found")
		return
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

}

func (j *UpdateExpiredSubscriptionJob) Schedule() string {
	return "0 0 0 * * *"
}
