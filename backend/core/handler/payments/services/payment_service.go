package services

import (
	"context"
	"database/sql"
	"email-marketing-service/internal/common"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/domain"
	"email-marketing-service/internal/factory/paymentFactory"
	"email-marketing-service/internal/helper"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/sqlc-dev/pqtype"
)

type PaymentService struct {
	store db.Store
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
		Status: sql.NullString{String: "processing", Valid: true},
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

func (s *PaymentService) VerifyPayment() {

}
