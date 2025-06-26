package middleware

import (
	"email-marketing-service/core/handler/permission"
	"email-marketing-service/internal/common"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/helper"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

func RequireFeatureAccess(store db.Store, feature string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, companyId, err := helper.ExtractUserId(r)
			if err != nil {
				helper.ErrorResponse(w, fmt.Errorf("auth error"), nil)
				return
			}

			companyUUID, err := uuid.Parse(companyId)
			if err != nil {
				helper.ErrorResponse(w, common.ErrInvalidUUID, nil)
				return
			}

			sub, err := store.GetCurrentRunningSubscription(r.Context(), companyUUID)
			if err != nil {
				helper.ErrorResponse(w, common.ErrFetchingSubscription, nil)
				return
			}

			if err := permission.CheckFeatureAccess(sub.PlanName, feature); err != nil {
				helper.ErrorResponse(w, err, nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

/** usage
r.Handle("/upload",
    RequireFeatureAccess(store, "upload_csv")(
        http.HandlerFunc(controller.UploadContactViaCSV),
    ),
).Methods("POST")


***/
