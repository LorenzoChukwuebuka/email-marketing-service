package smtp_server 

import (
	"time"
)


var RLconfig = &RateLimiterConfig{
    ConnectionsPerIP:   50,                 // 50 connections
    ConnectionWindow:   30 * time.Minute,   // per 30 minutes
    MessagesPerIP:     500,                // 500 messages
    MessageWindow:     time.Hour,          // per hour
    RecipientsPerMsg:  100,               // 100 recipients per message
    DefaultMonthQuota: 50000,             // 50,000 messages per month
}


var RCconfig = &RelayConfig{
    Debug:          true,
    DialTimeout:    15 * time.Second,
    RetryAttempts:  5,
    RetryDelay:     3 * time.Second,
    PreferredPorts: []string{"587", "465"},  // Skip port 25 if not needed
}