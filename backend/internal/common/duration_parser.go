package common

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseDuration parses duration strings like "1 month", "30 days", "1 year", etc.
func ParseDuration(duration string) (time.Duration, error) {
	duration = strings.TrimSpace(strings.ToLower(duration))
	parts := strings.Fields(duration)
	
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid duration format: %s", duration)
	}

	value, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid duration value: %s", parts[0])
	}

	unit := parts[1]
	// Handle plural forms
	if strings.HasSuffix(unit, "s") {
		unit = unit[:len(unit)-1]
	}

	switch unit {
	case "day":
		return time.Duration(value) * 24 * time.Hour, nil
	case "week":
		return time.Duration(value) * 7 * 24 * time.Hour, nil
	case "month":
		return time.Duration(value) * 30 * 24 * time.Hour, nil
	case "year":
		return time.Duration(value) * 365 * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unsupported duration unit: %s", unit)
	}
}

// GetDurationInDays returns the number of days for a given duration string
func GetDurationInDays(duration string) (int, error) {
	d, err := ParseDuration(duration)
	if err != nil {
		return 0, err
	}
	return int(d.Hours() / 24), nil
}

// GetBillingCycle converts duration to billing cycle string
func GetBillingCycle(duration string) string {
	duration = strings.TrimSpace(strings.ToLower(duration))
	
	if strings.Contains(duration, "month") {
		return "monthly"
	} else if strings.Contains(duration, "year") {
		return "yearly"
	} else if strings.Contains(duration, "week") {
		return "weekly"
	} else {
		return "custom"
	}
}