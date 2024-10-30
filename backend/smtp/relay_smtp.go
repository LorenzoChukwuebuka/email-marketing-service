package smtp_server

import (
	"fmt"
	"net"
	"net/smtp"
	"os"
	"sort"
	"strings"
)

// RelayService handles email relay functionality
type RelayService struct {
	debug bool
}

// NewRelayService creates a new relay service
func NewRelayService(debug bool) *RelayService {
	return &RelayService{
		debug: debug,
	}
}

// MXRecord represents an MX record with its preference
type MXRecord struct {
	Host       string
	Preference uint16
}

// RelayEmail handles the email relay process
func (s *RelayService) RelayEmail(from string, to []string, data []byte) error {
	// Group recipients by domain
	domainRecipients := make(map[string][]string)
	for _, recipient := range to {
		parts := strings.Split(recipient, "@")
		if len(parts) != 2 {
			continue
		}
		domain := parts[1]
		domainRecipients[domain] = append(domainRecipients[domain], recipient)
	}

	// Process each domain separately
	var errors []error
	for domain, recipients := range domainRecipients {
		err := s.relayToDomain(domain, from, recipients, data)
		if err != nil {
			s.debugLog(fmt.Sprintf("Failed to relay to domain %s: %v", domain, err))
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("relay errors occurred: %v", errors)
	}
	return nil
}

func (s *RelayService) relayToDomain(domain, from string, recipients []string, data []byte) error {
	mxRecords, err := s.lookupMXRecords(domain)
	if err != nil {
		return fmt.Errorf("MX lookup failed for domain %s: %w", domain, err)
	}

	var lastError error
	for _, mx := range mxRecords {
		err = s.sendToMXServer(mx.Host, from, recipients, data)
		if err == nil {
			return nil
		}
		lastError = err
		s.debugLog(fmt.Sprintf("Failed to deliver to MX %s: %v", mx.Host, err))
	}

	return fmt.Errorf("failed to deliver to any MX server: %w", lastError)
}

func (s *RelayService) lookupMXRecords(domain string) ([]MXRecord, error) {
	mxs, err := net.LookupMX(domain)
	if err != nil {
		// Try A record as fallback
		_, err := net.LookupHost(domain)
		if err != nil {
			return nil, fmt.Errorf("no MX records and A record lookup failed: %w", err)
		}
		return []MXRecord{{Host: domain, Preference: 0}}, nil
	}

	records := make([]MXRecord, len(mxs))
	for i, mx := range mxs {
		records[i] = MXRecord{
			Host:       strings.TrimSuffix(mx.Host, "."),
			Preference: mx.Pref,
		}
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Preference < records[j].Preference
	})

	return records, nil
}

func (s *RelayService) sendToMXServer(mxHost, from string, to []string, data []byte) error {
	conn, err := net.Dial("tcp", mxHost+":25")
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer conn.Close()

	c, err := smtp.NewClient(conn, mxHost)
	if err != nil {
		return fmt.Errorf("SMTP client creation failed: %w", err)
	}
	defer c.Close()

	hostname := s.getLocalHostname()
	if err = c.Hello(hostname); err != nil {
		return fmt.Errorf("HELO failed: %w", err)
	}

	if err = c.Mail(from); err != nil {
		return fmt.Errorf("MAIL FROM failed: %w", err)
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			s.debugLog(fmt.Sprintf("RCPT TO failed for %s: %v", addr, err))
			continue
		}
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("DATA command failed: %w", err)
	}

	if _, err = w.Write(data); err != nil {
		return fmt.Errorf("data write failed: %w", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("data close failed: %w", err)
	}

	return c.Quit()
}

func (s *RelayService) getLocalHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "localhost"
	}
	return hostname
}

func (s *RelayService) debugLog(message string) {
	if s.debug {
		fmt.Printf("[RELAY DEBUG] %s\n", message)
	}
}
