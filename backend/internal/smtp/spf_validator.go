package smtp_server

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// Result represents possible outcomes of SPF validation
type Result int

const (
	Pass Result = iota
	Fail
	SoftFail
	Neutral
	TempError
	PermError
)

// Config holds SPF validator configuration
type Config struct {
	// CacheDuration specifies how long to cache SPF results
	CacheDuration time.Duration
	// MaxLookups specifies maximum DNS lookups per check (RFC 7208 recommends 10)
	MaxLookups int
	// EnableCache determines if results should be cached
	EnableCache bool
}

// DefaultConfig returns recommended SPF configuration
func DefaultConfig() *Config {
	return &Config{
		CacheDuration: 1 * time.Hour,
		MaxLookups:    10,
		EnableCache:   true,
	}
}

// Validator handles SPF validation
type Validator struct {
	// Configuration
	config *Config

	// Cache for SPF results
	cache map[string]cacheEntry
	mutex sync.RWMutex

	// DNS lookup count for current validation
	lookupCount int
}

type cacheEntry struct {
	result    Result
	timestamp time.Time
}

// New creates a new SPF validator with given configuration
func New(config *Config) *Validator {
	if config == nil {
		config = DefaultConfig()
	}

	return &Validator{
		config: config,
		cache:  make(map[string]cacheEntry),
	}
}

// CheckHost performs SPF validation for an IP address against a domain
// Returns: Result and error if any occurred
func (v *Validator) CheckHost(ip string, domain string, sender string) (Result, error) {
	// Reset lookup count for new check
	v.lookupCount = 0

	// Check cache first if enabled
	if v.config.EnableCache {
		if result, found := v.checkCache(ip, domain); found {
			return result, nil
		}
	}

	// Basic validation
	if domain == "" {
		return PermError, fmt.Errorf("empty domain")
	}

	// Parse IP address
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return PermError, fmt.Errorf("invalid IP address: %s", ip)
	}

	// Get SPF record
	record, err := v.getSPFRecord(domain)
	if err != nil {
		return TempError, fmt.Errorf("DNS lookup error: %w", err)
	}

	// No SPF record found
	if record == "" {
		return Neutral, nil
	}

	// Evaluate SPF record
	result := v.evaluateRecord(record, ipAddr, domain)

	// Cache result if enabled
	if v.config.EnableCache {
		v.cacheResult(ip, domain, result)
	}

	return result, nil
}

// getSPFRecord retrieves the SPF record for a domain
func (v *Validator) getSPFRecord(domain string) (string, error) {
	// Increment lookup count
	v.lookupCount++
	if v.lookupCount > v.config.MaxLookups {
		return "", fmt.Errorf("exceeded maximum DNS lookups")
	}

	// Lookup TXT records
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		return "", err
	}

	// Find SPF record
	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			return record, nil
		}
	}

	return "", nil
}

// evaluateRecord processes an SPF record and returns the result
func (v *Validator) evaluateRecord(record string, ip net.IP, domain string) Result {
	debugLog(fmt.Sprintf("Evaluating SPF record: %s for IP: %s and domain: %s", record, ip, domain))

	terms := strings.Fields(record)
	for _, term := range terms[1:] {
		qualifier := "+"
		mechanism := term

		if strings.ContainsAny(term[0:1], "+-?~") {
			qualifier = term[0:1]
			mechanism = term[1:]
		}

		debugLog(fmt.Sprintf("Checking mechanism: %s with qualifier: %s", mechanism, qualifier))
		if v.checkMechanism(mechanism, ip, domain) {
			debugLog(fmt.Sprintf("Mechanism matched! Returning result for qualifier: %s", qualifier))
			switch qualifier {
			case "+":
				return Pass
			case "-":
				return Fail
			case "~":
				return SoftFail
			case "?":
				return Neutral
			}
		}
	}

	debugLog("No mechanisms matched, returning Neutral")
	return Neutral
}

// checkMechanism evaluates a single SPF mechanism
func (v *Validator) checkMechanism(mechanism string, ip net.IP, domain string) bool {
	// Split mechanism into type and value
	parts := strings.SplitN(mechanism, ":", 2)
	mechType := parts[0]
	value := ""
	if len(parts) > 1 {
		value = parts[1]
	}

	switch mechType {
	case "all":
		return true

	case "ip4", "ip6":
		return v.checkIP(value, ip)

	case "a":
		return v.checkA(value, domain, ip)

	case "mx":
		return v.checkMX(value, domain, ip)

		// TODO: Implement additional mechanisms
		// - ptr
		// - exists
		// - include
		// - redirect
	}

	return false
}

// checkIP validates if an IP matches a CIDR range
func (v *Validator) checkIP(cidr string, ip net.IP) bool {
	debugLog(fmt.Sprintf("Checking IP %s against CIDR %s", ip, cidr))
	// Add /32 if no prefix specified
	if !strings.Contains(cidr, "/") {
		if ip.To4() != nil {
			cidr += "/32"
		} else {
			cidr += "/128"
		}
	}

	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		debugLog(fmt.Sprintf("Error parsing CIDR: %v", err))
		return false
	}
	result := network.Contains(ip)
	debugLog(fmt.Sprintf("IP check result: %v", result))
	return result
}

// checkA validates if an IP matches domain's A/AAAA records
func (v *Validator) checkA(value, domain string, ip net.IP) bool {
	// Use domain if no value specified
	if value == "" {
		value = domain
	}

	// Increment lookup count
	v.lookupCount++
	if v.lookupCount > v.config.MaxLookups {
		return false
	}

	// Lookup A/AAAA records
	ips, err := net.LookupIP(value)
	if err != nil {
		return false
	}

	// Check if IP matches any record
	for _, recordIP := range ips {
		if ip.Equal(recordIP) {
			return true
		}
	}

	return false
}

// checkMX validates if an IP matches domain's MX records
func (v *Validator) checkMX(value, domain string, ip net.IP) bool {
	// Use domain if no value specified
	if value == "" {
		value = domain
	}

	// Increment lookup count
	v.lookupCount++
	if v.lookupCount > v.config.MaxLookups {
		return false
	}

	// Lookup MX records
	mxRecords, err := net.LookupMX(value)
	if err != nil {
		return false
	}

	// Check each MX record
	for _, mx := range mxRecords {
		// Increment lookup count for each MX lookup
		v.lookupCount++
		if v.lookupCount > v.config.MaxLookups {
			return false
		}

		// Lookup A/AAAA records for MX
		ips, err := net.LookupIP(mx.Host)
		if err != nil {
			continue
		}

		// Check if IP matches any MX record
		for _, recordIP := range ips {
			if ip.Equal(recordIP) {
				return true
			}
		}
	}

	return false
}

// Cache management methods
func (v *Validator) checkCache(ip, domain string) (Result, bool) {
	v.mutex.RLock()
	defer v.mutex.RUnlock()

	key := fmt.Sprintf("%s:%s", domain, ip)
	if entry, exists := v.cache[key]; exists {
		if time.Since(entry.timestamp) < v.config.CacheDuration {
			return entry.result, true
		}
	}
	return Neutral, false
}

func (v *Validator) cacheResult(ip, domain string, result Result) {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	key := fmt.Sprintf("%s:%s", domain, ip)
	v.cache[key] = cacheEntry{
		result:    result,
		timestamp: time.Now(),
	}
}
