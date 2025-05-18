package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
)

var (
	secretKey   = "123456"
	hashVersion = "1"
)

func GeneratePaymentHash(paymentID uuid.UUID, userID uuid.UUID, amount int64, subscriptionId uuid.UUID) (string, error) {
	h := hmac.New(sha256.New, []byte(secretKey))

	// Concatenate important wallet data
	data := fmt.Sprintf("%s:%s:%d:%s:%s",
		paymentID.String(),
		userID.String(),
		amount,
		subscriptionId,
		hashVersion,
	)

	_, err := h.Write([]byte(data))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
