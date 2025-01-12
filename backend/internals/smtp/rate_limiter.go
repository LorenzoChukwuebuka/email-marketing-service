package smtp_server

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// RateLimiter handles rate limiting for SMTP connections and sending
type RateLimiter struct {
	// Connection rate limiting
	connections    map[string][]time.Time
	connectionsMax int
	connectWindow  time.Duration

	// Message rate limiting
	messages      map[string][]time.Time
	messagesMax   int
	messageWindow time.Duration

	// Recipients per message limiting
	recipientsMax int

	// Monthly message quotas
	monthlyQuotas    map[string]int
	monthlyQuotasMax map[string]int

	mutex sync.RWMutex
}

// RateLimiterConfig holds configuration for rate limiting
type RateLimiterConfig struct {
	ConnectionsPerIP  int           // Maximum connections per IP in window
	ConnectionWindow  time.Duration // Time window for connection counting
	MessagesPerIP     int           // Maximum messages per IP in window
	MessageWindow     time.Duration // Time window for message counting
	RecipientsPerMsg  int           // Maximum recipients per message
	DefaultMonthQuota int           // Default monthly message quota per sender
}

// DefaultRateLimiterConfig returns recommended rate limiting settings
func DefaultRateLimiterConfig() *RateLimiterConfig {
	return &RateLimiterConfig{
		ConnectionsPerIP:  30,        // 30 connections
		ConnectionWindow:  time.Hour, // per hour
		MessagesPerIP:     300,       // 300 messages
		MessageWindow:     time.Hour, // per hour
		RecipientsPerMsg:  50,        // 50 recipients per message
		DefaultMonthQuota: 10000,     // 10,000 messages per month
	}
}

// NewRateLimiter creates a new rate limiter with given configuration
func NewRateLimiter(config *RateLimiterConfig) *RateLimiter {
	if config == nil {
		config = DefaultRateLimiterConfig()
	}

	return &RateLimiter{
		connections:      make(map[string][]time.Time),
		connectionsMax:   config.ConnectionsPerIP,
		connectWindow:    config.ConnectionWindow,
		messages:         make(map[string][]time.Time),
		messagesMax:      config.MessagesPerIP,
		messageWindow:    config.MessageWindow,
		recipientsMax:    config.RecipientsPerMsg,
		monthlyQuotas:    make(map[string]int),
		monthlyQuotasMax: make(map[string]int),
	}
}

// CheckConnection checks if a new connection is allowed from the IP
func (rl *RateLimiter) CheckConnection(ip string) error {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	window := now.Add(-rl.connectWindow)

	// Clean old entries
	if times, exists := rl.connections[ip]; exists {
		var valid []time.Time
		for _, t := range times {
			if t.After(window) {
				valid = append(valid, t)
			}
		}
		rl.connections[ip] = valid
	}

	// Check limit
	if len(rl.connections[ip]) >= rl.connectionsMax {
		return fmt.Errorf("rate limit exceeded: too many connections from %s", ip)
	}

	// Add new connection
	rl.connections[ip] = append(rl.connections[ip], now)
	return nil
}

// CheckMessage checks if a new message is allowed from the IP
func (rl *RateLimiter) CheckMessage(ip string, from string, recipients []string) error {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	// Check recipients limit
	if len(recipients) > rl.recipientsMax {
		return fmt.Errorf("too many recipients: %d (max %d)", len(recipients), rl.recipientsMax)
	}

	now := time.Now()
	window := now.Add(-rl.messageWindow)

	// Clean old entries
	if times, exists := rl.messages[ip]; exists {
		var valid []time.Time
		for _, t := range times {
			if t.After(window) {
				valid = append(valid, t)
			}
		}
		rl.messages[ip] = valid
	}

	// Check rate limit
	if len(rl.messages[ip]) >= rl.messagesMax {
		return fmt.Errorf("rate limit exceeded: too many messages from %s", ip)
	}

	// Check monthly quota
	month := now.Format("2006-01")
	quotaKey := fmt.Sprintf("%s:%s", from, month)

	if quota, exists := rl.monthlyQuotasMax[from]; exists {
		if rl.monthlyQuotas[quotaKey] >= quota {
			return fmt.Errorf("monthly quota exceeded for %s", from)
		}
	} else {
		// Use default quota if no custom quota is set
		if rl.monthlyQuotas[quotaKey] >= DefaultRateLimiterConfig().DefaultMonthQuota {
			return fmt.Errorf("monthly quota exceeded for %s", from)
		}
	}

	// Add new message
	rl.messages[ip] = append(rl.messages[ip], now)
	rl.monthlyQuotas[quotaKey]++
	return nil
}

// SetCustomQuota sets a custom monthly quota for a sender
func (rl *RateLimiter) SetCustomQuota(sender string, quota int) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	rl.monthlyQuotasMax[sender] = quota
}

// GetQuotaUsage returns the current month's quota usage for a sender
func (rl *RateLimiter) GetQuotaUsage(sender string) int {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()

	month := time.Now().Format("2006-01")
	quotaKey := fmt.Sprintf("%s:%s", sender, month)
	return rl.monthlyQuotas[quotaKey]
}

// Cleanup removes old entries from the rate limiter
func (rl *RateLimiter) Cleanup() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	connWindow := now.Add(-rl.connectWindow)
	msgWindow := now.Add(-rl.messageWindow)

	// Clean old connection entries
	for ip, times := range rl.connections {
		var valid []time.Time
		for _, t := range times {
			if t.After(connWindow) {
				valid = append(valid, t)
			}
		}
		if len(valid) == 0 {
			delete(rl.connections, ip)
		} else {
			rl.connections[ip] = valid
		}
	}

	// Clean old message entries
	for ip, times := range rl.messages {
		var valid []time.Time
		for _, t := range times {
			if t.After(msgWindow) {
				valid = append(valid, t)
			}
		}
		if len(valid) == 0 {
			delete(rl.messages, ip)
		} else {
			rl.messages[ip] = valid
		}
	}

	// Clean old monthly quotas (keep only current month)
	currentMonth := now.Format("2006-01")
	for key := range rl.monthlyQuotas {
		if !strings.Contains(key, currentMonth) {
			delete(rl.monthlyQuotas, key)
		}
	}
}
