package services

import (
	"context"
	"database/sql"
	"email-marketing-service/core/handler/payments/dto"
	"email-marketing-service/core/handler/payments/mapper"
	"email-marketing-service/internal/common"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/domain"
	"email-marketing-service/internal/enums"
	"email-marketing-service/internal/factory/paymentFactory"
	"email-marketing-service/internal/helper"
	worker "email-marketing-service/internal/workers"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sqlc-dev/pqtype"
)

type PaymentService struct {
	store db.Store
	wkr   *worker.Worker
}

func NewPaymentService(store db.Store) *PaymentService {
	return &PaymentService{store: store}
}

func (s *PaymentService) InitiateNewTransaction(ctx context.Context, req domain.BasePaymentModelData) (any, error) {
	if err := helper.ValidateData(req); err != nil {
		return nil, errors.Join(common.ErrValidatingRequest, err)
	}

	_uuid, err := common.ParseUUIDMap(map[string]string{
		"user":    req.UserId,
		"company": req.CompanyId,
	})

	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	//create a new payment intent
	amount := decimal.NewFromFloat(req.AmountToPay)

	paymentIntent, err := s.store.CreatePaymentIntent(ctx, db.CreatePaymentIntentParams{
		UserID:             _uuid["user"],
		CompanyID:          _uuid["company"],
		PaymentIntentID:    req.PaymentMethod,
		Amount:             amount,
		Status:             sql.NullString{String: "pending", Valid: true},
		ReceiptEmail:       sql.NullString{String: req.Email, Valid: true},
		PaymentMethodTypes: []string{"card", "ussd", "bank transfer"},
	})

	if err != nil {
		return nil, err
	}

	req.PaymentIntentID = paymentIntent.ID.String()

	//call the payment processor
	paymentService, err := paymentmethodFactory.PaymentFactory(req.PaymentMethod)

	if err != nil {
		return nil, errors.Join(common.ErrPaymentMethodNotSupported, err)
	}

	result, err := paymentService.OpenDeposit(&req)
	if err != nil {
		return nil, nil
	}

	encoded, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	var paymentMethodID string

	if dataMap, ok := result["data"].(map[string]any); ok {
		if ref, ok := dataMap["reference"].(string); ok {
			paymentMethodID = ref
		} else {
			return nil, fmt.Errorf("reference is not a string")
		}
	} else {
		return nil, fmt.Errorf("data is not a map[string]interface{}")
	}

	_, err = s.store.UpdatePaymentIntent(ctx, db.UpdatePaymentIntentParams{
		ID:     paymentIntent.ID,
		Status: sql.NullString{String: string(enums.PaymentIntentProcessing), Valid: true},
		Metadata: pqtype.NullRawMessage{
			RawMessage: encoded,
			Valid:      true,
		},
		PaymentIntentID: sql.NullString{String: paymentMethodID, Valid: true},
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PaymentService) HandleWebhook(ctx context.Context, r *http.Request, paymentMethod string) (any, error) {

	httpReq := common.RequestFromCtx(ctx)
	var ip net.IP
	if httpReq != nil {
		ip = s.wkr.GetClientIP(httpReq)
	}

	paymentService, err := paymentmethodFactory.PaymentFactory(paymentMethod)

	if err != nil {
		return nil, errors.Join(common.ErrPaymentMethodNotSupported, err)
	}

	data, err := paymentService.WebhookHandler(r)
	if err != nil {

		data := worker.AuditCreatePayload{
			UserID:      uuid.New(),
			ResourceID:  nil,
			Method:      r.Method,
			Endpoint:    r.URL.Path,
			IP:          ip,
			Success:     false,
			RequestBody: r,
		}

		if _, err := s.wkr.EnqueueTask(ctx, worker.TaskAuditLogFailedLogin, data); err != nil {
			// Just log the error, don't fail the request
			log.Printf("Failed to enqueue failed payment intent audit: %v", err)
		}
		return nil, nil
	}

	//call the verify function and pass the payment method and reference with context to it

	//_, err = s.VerifyPayment(ctx, paymentMethod, data.Reference)

	return data, nil
}

func (s *PaymentService) VerifyPayment(ctx context.Context, paymentMethod, reference string) (any, error) {
	//call the payment processor
	paymentService, err := paymentmethodFactory.PaymentFactory(paymentMethod)

	if err != nil {
		return nil, errors.Join(common.ErrPaymentMethodNotSupported, err)
	}

	params := &domain.BaseProcessPaymentModel{
		PaymentMethod: paymentMethod,
		Reference:     reference,
	}

	data, err := paymentService.ProcessDeposit(params)
	if err != nil {
		// Try to update payment intent with processing error if we have reference
		if reference != "" {
			s.updatePaymentIntentErrorWithReference(ctx, reference, fmt.Sprintf("Payment processing failed: %v", err))
		}
		return nil, err
	}

	_uuid, err := common.ParseUUIDMap(map[string]string{
		"payment_intent": data.PaymentIntentID,
	})
	if err != nil {
		// Try to update payment intent with UUID parsing error
		if data != nil && data.PaymentIntentID != "" {
			s.updatePaymentIntentErrorWithReference(ctx, reference, fmt.Sprintf("Invalid UUID in payment response: %v", err))
		}
		return nil, common.ErrInvalidUUID
	}

	// //check if payment intent already exists in payment
	intentExists, err := s.store.CheckPaymentIntentExists(ctx, sql.NullString{String: data.PaymentIntentID, Valid: true})

	if err != nil {
		s.updatePaymentIntentErrorWithReference(ctx, reference, fmt.Sprintf("error fetching intent: %v", err))
		return nil, err
	}

	if intentExists {
		s.updatePaymentIntentErrorWithReference(ctx, reference, "duplicate intent")
		return nil, nil
	}

	//start a db transaction
	err = s.store.ExecTx(ctx, func(q *db.Queries) error {
		// Create subscription and get subscription ID
		subscriptionID, err := s.createsubscription(ctx, *q, *data, reference, paymentMethod)
		if err != nil {
			return err
		}

		// Update the payment intent with success details
		now := time.Now().UTC()
		amount := decimal.NewFromFloat(data.Amount)

		_, err = q.UpdatePaymentIntent(ctx, db.UpdatePaymentIntentParams{
			ID:             _uuid["payment_intent"],
			SubscriptionID: uuid.NullUUID{UUID: subscriptionID, Valid: true},
			Amount:         amount,
			Status:         sql.NullString{String: string(enums.PaymentIntentSuccessful), Valid: true},
			SucceededAt:    sql.NullTime{Time: now, Valid: true},
		})

		if err != nil {
			return fmt.Errorf("error updating payment intent status: %w", err)
		}

		return nil
	})

	if err != nil {
		s.updatePaymentIntentErrorWithReference(ctx, reference, fmt.Sprintf("transaction failed: %v", err))
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	return data, nil
}

func (s *PaymentService) updatePaymentIntentErrorWithReference(ctx context.Context, paymentIntentID, errorMessage string) {
	err := s.store.ExecTx(ctx, func(q *db.Queries) error {
		// First try to get the payment intent to see if it exists
		paymentIntent, err := q.GetPaymentIntentByPaymentIntentID(ctx, paymentIntentID)
		if err != nil {
			log.Printf("Payment intent not found with ID %s: %v", paymentIntentID, err)
			return err
		}

		log.Printf("Found payment intent: %s, updating with error: %s", paymentIntent.ID, errorMessage)

		_, err = q.UpdatePaymentIntentError(ctx, db.UpdatePaymentIntentErrorParams{
			PaymentIntentID:  paymentIntent.PaymentIntentID,
			Status:           sql.NullString{String: string(enums.PaymentIntentFailed), Valid: true},
			LastPaymentError: pqtype.NullRawMessage{RawMessage: json.RawMessage(fmt.Sprintf(`"%s"`, errorMessage)), Valid: true},
		})

		if err != nil {
			log.Printf("Error executing UpdatePaymentIntentError: %v", err)
			return err
		}

		log.Println("Successfully updated payment intent")
		return nil
	})

	if err != nil {
		log.Printf("Transaction failed when updating payment intent error: %v", err)
	}
}

func (s *PaymentService) createsubscription(ctx context.Context, q db.Queries, data domain.BasePaymentResponse, reference, paymentMethod string) (uuid.UUID, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"user":           data.UserID,
		"company":        data.CompanyID,
		"payment_intent": data.PaymentIntentID,
		"plan":           data.PlanID,
	})
	if err != nil {
		return uuid.Nil, common.ErrInvalidUUID
	}

	//get the plan to make sure it actually exists
	_, err = q.GetPlanByID(ctx, _uuid["plan"])
	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.Nil, fmt.Errorf("no plan with this id found: %w", err)
		}
		return uuid.Nil, errors.Join(common.ErrFetchingRecord, err)
	}

	//check if user have any current running subscription
	current_running_subscription, err := q.GetCurrentRunningSubscription(ctx, _uuid["company"])
	if err != nil {
		return uuid.Nil, fmt.Errorf("error fetching current subscription: %w", err)
	}

	//cancel the current subscription and delete the email usage for this subscription
	_, err = q.UpdateSubscriptionStatus(ctx, db.UpdateSubscriptionStatusParams{
		ID:     current_running_subscription.SubscriptionID,
		Status: sql.NullString{String: string(enums.SubscriptionExpired), Valid: true},
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("error updating subscription status: %w", err)
	}

	//delete their email usage
	err = q.DeleteEmailUsageByCompanyIDAndSubscriptionID(ctx, db.DeleteEmailUsageByCompanyIDAndSubscriptionIDParams{
		CompanyID:      _uuid["company"],
		SubscriptionID: current_running_subscription.SubscriptionID,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("error deleting email usage: %w", err)
	}

	//get the last or current payment
	currentPayment, err := q.GetLastPaymentByCompanyID(ctx, _uuid["company"])
	if err != nil {
		return uuid.Nil, fmt.Errorf("error fetching last/current payment: %w", err)
	}

	//get the previous payment hash and compare it with the record
	verify_hash, err := common.VerifyPaymentHash(currentPayment.IntegrityHash.String, currentPayment.ID, currentPayment.UserID, currentPayment.Amount.BigInt().Int64(), currentPayment.SubscriptionID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error verifying hash: %w", err)
	}

	if !verify_hash {
		return uuid.Nil, fmt.Errorf("hash could not be verified: %w", err)
	}

	//go ahead and create the subscription plan
	// Parse the duration
	durationDays, err := common.GetDurationInDays(data.Duration)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing duration: %w", err)
	}
	now := time.Now().UTC()

	amount := decimal.NewFromFloat(data.Amount)
	// Create new subscription
	subscription, err := q.CreateSubscription(ctx, db.CreateSubscriptionParams{
		CompanyID:       _uuid["company"],
		PlanID:          _uuid["plan"],
		Amount:          amount,
		BillingCycle:    sql.NullString{String: data.Duration, Valid: true},
		TrialStartsAt:   sql.NullTime{}, // Set if needed
		TrialEndsAt:     sql.NullTime{}, // Set if needed
		StartsAt:        sql.NullTime{Time: now, Valid: true},
		EndsAt:          sql.NullTime{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true},
		Status:          sql.NullString{String: string(enums.SubscriptionActive), Valid: true},
		NextBillingDate: sql.NullTime{Time: time.Now().Add(31 * 24 * time.Hour), Valid: true},
		AutoRenew:       sql.NullBool{Bool: true, Valid: true}, // Set based on your business logic
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("error creating subscription: %w", err)
	}

	// Create payment record
	payment, err := q.CreatePayment(ctx, db.CreatePaymentParams{
		CompanyID:            _uuid["company"],
		UserID:               _uuid["user"],
		SubscriptionID:       subscription.ID,
		Amount:               amount,
		PaymentID:            sql.NullString{String: data.PaymentIntentID, Valid: true},
		Currency:             sql.NullString{String: "NGN", Valid: true}, // Or from data
		PaymentMethod:        sql.NullString{String: paymentMethod, Valid: true},
		Status:               sql.NullString{String: "completed", Valid: true},
		Notes:                sql.NullString{String: "Payment successful", Valid: true},
		TransactionReference: sql.NullString{String: reference, Valid: true},
		PaymentDate:          sql.NullTime{Time: now, Valid: true},
		BillingPeriodStart:   sql.NullTime{Time: now, Valid: true},
		BillingPeriodEnd:     sql.NullTime{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true},
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("error creating payment: %w", err)
	}

	// Generate and update payment hash
	paymentHash, err := common.GeneratePaymentHash(payment.ID, _uuid["user"], payment.Amount.BigInt().Int64(), subscription.ID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error generating payment hash: %w", err)
	}

	err = q.UpdatePaymentHash(ctx, db.UpdatePaymentHashParams{
		IntegrityHash: sql.NullString{String: paymentHash, Valid: true},
		ID:            payment.ID,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("error updating payment hash: %w", err)
	}

	// Create daily email usage records
	mailingLimits, err := q.GetMailingLimitByPlanID(ctx, _uuid["plan"])
	if err != nil {
		return uuid.Nil, fmt.Errorf("error getting mailing limit: %w", err)
	}

	for i := 0; i < durationDays; i++ {
		periodStart := now.AddDate(0, 0, i).Truncate(24 * time.Hour)
		periodEnd := periodStart.Add(24 * time.Hour).Add(-time.Second)

		_, err := q.CreateDailyEmailUsage(ctx, db.CreateDailyEmailUsageParams{
			CompanyID:        _uuid["company"],
			SubscriptionID:   subscription.ID,
			EmailsSent:       sql.NullInt32{Int32: 0, Valid: true},
			EmailsLimit:      mailingLimits.DailyLimit.Int32,
			UsagePeriodStart: periodStart,
			UsagePeriodEnd:   periodEnd,
			RemainingEmails:  sql.NullInt32{Int32: mailingLimits.DailyLimit.Int32, Valid: true},
		})
		if err != nil {
			return uuid.Nil, fmt.Errorf("error creating daily email usage for %s: %w", periodStart.Format("2006-01-02"), err)
		}
	}

	return subscription.ID, nil
}

func (s *PaymentService) GetCompanyCurrentRunningSubscription() {}

func (s *PaymentService) GetAllPaymentsForACompany(ctx context.Context, d *dto.FetchPayment) (any, error) {
	_uuid, err := common.ParseUUIDMap(map[string]string{
		"company": d.CompanyID,
		"user":    d.UserId})

	if err != nil {
		return nil, common.ErrInvalidUUID
	}
	payments, err := s.store.GetPaymentsByCompanyAndUser(ctx, db.GetPaymentsByCompanyAndUserParams{
		CompanyID: _uuid["company"],
		UserID:    _uuid["user"],
		Limit:     int32(d.Limit),
		Offset:    int32(d.Offset),
	})

	if err != nil {
		return nil, common.ErrFetchingRecord
	}
	payment_count, err := s.store.GetPaymentCounts(ctx, _uuid["company"])
	if err != nil {
		return nil, common.ErrFetchingRecord
	}
	response := mapper.MapPaymentResponses(payments)

	items := make([]any, len(response))
	for i, v := range response {
		items[i] = v
	}
	data := common.Paginate(int(payment_count), items, d.Offset, d.Limit)
	return data, nil
}
