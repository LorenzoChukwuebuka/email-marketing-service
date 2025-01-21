package utils

import (
	"strings"
)

func IsCompliant(content string) bool {
	// Basic restricted keywords
	restrictedWords := []string{"get-rich-quick", "win money", "free meds", "adult content"}

	for _, word := range restrictedWords {
		if strings.Contains(strings.ToLower(content), word) {
			return false // Violation detected
		}
	}
	return true
}
