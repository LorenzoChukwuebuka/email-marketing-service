package helper

import (
	"math/rand"
	"time"
)

func GenerateOTP(length int) string {

	rand.Seed(time.Now().UnixNano())

	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	result := make([]byte, length)

	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}

func GenerateUniqueRandomNumbers(n int) string {
	rand.Seed(time.Now().UnixNano())

	const charset = "0123456789"

	result := make([]byte, n)

	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}
