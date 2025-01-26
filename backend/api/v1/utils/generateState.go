package utils

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

// generateRandomState generates a secure random state string.
func GenerateRandomState() string {
	// Define the number of bytes for the state. 16 bytes is typical.
	const stateBytes = 16

	// Create a byte slice to hold the random data.
	b := make([]byte, stateBytes)

	// Fill the byte slice with cryptographically secure random bytes.
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Failed to generate random state: %v", err)
	}

	// Encode the random bytes as a URL-safe base64 string.
	return base64.URLEncoding.EncodeToString(b)
}
