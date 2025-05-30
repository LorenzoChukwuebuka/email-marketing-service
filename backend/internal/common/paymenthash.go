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

// VerifyPaymentHash verifies if the provided hash matches the expected hash for the given parameters
func VerifyPaymentHash(providedHash string, paymentID uuid.UUID, userID uuid.UUID, amount int64, subscriptionId uuid.UUID) (bool, error) {
	// Generate the expected hash using the same parameters
	expectedHash, err := GeneratePaymentHash(paymentID, userID, amount, subscriptionId)
	if err != nil {
		return false, fmt.Errorf("failed to generate expected hash: %w", err)
	}

	// Use hmac.Equal for constant-time comparison to prevent timing attacks
	providedHashBytes, err := hex.DecodeString(providedHash)
	if err != nil {
		return false, fmt.Errorf("invalid hash format: %w", err)
	}

	expectedHashBytes, err := hex.DecodeString(expectedHash)
	if err != nil {
		return false, fmt.Errorf("failed to decode expected hash: %w", err)
	}

	return hmac.Equal(providedHashBytes, expectedHashBytes), nil
}

// Alternative simpler verification if you're confident about hex encoding
func VerifyPaymentHashSimple(providedHash string, paymentID uuid.UUID, userID uuid.UUID, amount int64, subscriptionId uuid.UUID) (bool, error) {
	expectedHash, err := GeneratePaymentHash(paymentID, userID, amount, subscriptionId)
	if err != nil {
		return false, err
	}

	return providedHash == expectedHash, nil
}
