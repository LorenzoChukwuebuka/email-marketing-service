package smtp_server

import (
    "bytes"
    "crypto/tls"
    "fmt"
    "net"
    "net/smtp"
    "os"
    "sort"
    "strings"
    "time"
)

// RelayConfig holds configuration for the relay service
type RelayConfig struct {
    Debug          bool
    DialTimeout    time.Duration
    RetryAttempts  int
    RetryDelay     time.Duration
    PreferredPorts []string
    TLSConfig      *tls.Config
}

// DefaultRelayConfig returns default configuration
func DefaultRelayConfig() *RelayConfig {
    return &RelayConfig{
        Debug:          true,
        DialTimeout:    5 * time.Second,
        RetryAttempts:  3,
        RetryDelay:     5 * time.Second,
        PreferredPorts: []string{"587", "465", "25"},
        TLSConfig:      &tls.Config{InsecureSkipVerify: false},
    }
}

// MXRecord represents an MX record with its preference
type MXRecord struct {
    Host       string
    Preference uint16
}

// RelayService handles email relay functionality
type RelayService struct {
    config *RelayConfig
}

// NewRelayService creates a new relay service with the given configuration
func NewRelayService(config *RelayConfig) *RelayService {
    if config == nil {
        config = DefaultRelayConfig()
    }
    return &RelayService{config: config}
}

// Email represents an email message
type Email struct {
    From    string
    To      []string
    Subject string
    Body    []byte
    Headers map[string]string
}

// RelayEmail handles the email relay process
func (s *RelayService) RelayEmail(email *Email) error {
    if len(email.To) == 0 {
        return fmt.Errorf("no recipients provided")
    }

    // Add headers to the email data
    messageWithHeaders := s.addHeaders(email)

    // Group recipients by domain for efficient delivery
    domainRecipients := make(map[string][]string)
    for _, recipient := range email.To {
        parts := strings.Split(recipient, "@")
        if len(parts) != 2 {
            s.debugLog(fmt.Sprintf("Invalid recipient address: %s", recipient))
            continue
        }
        domain := parts[1]
        domainRecipients[domain] = append(domainRecipients[domain], recipient)
    }

    // Process each domain separately
    var errors []error
    for domain, recipients := range domainRecipients {
        err := s.relayToDomain(domain, email.From, recipients, messageWithHeaders)
        if err != nil {
            s.debugLog(fmt.Sprintf("Failed to relay to domain %s: %v", domain, err))
            errors = append(errors, fmt.Errorf("domain %s: %w", domain, err))
        }
    }

    if len(errors) > 0 {
        return fmt.Errorf("relay errors occurred: %v", errors)
    }
    return nil
}

// addHeaders adds necessary headers to the email
func (s *RelayService) addHeaders(email *Email) []byte {
    var buffer bytes.Buffer

    // Add standard headers
    buffer.WriteString(fmt.Sprintf("From: %s\r\n", email.From))
    buffer.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(email.To, ", ")))
    buffer.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))
    buffer.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))
    buffer.WriteString("MIME-Version: 1.0\r\n")
    buffer.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
    buffer.WriteString(fmt.Sprintf("Message-ID: <%d@%s>\r\n", time.Now().UnixNano(), strings.Split(email.From, "@")[1]))

    // Add custom headers
    for key, value := range email.Headers {
        buffer.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
    }

    buffer.WriteString("\r\n") // Separate headers from body
    buffer.Write(email.Body)

    return buffer.Bytes()
}

// lookupMXRecords fetches MX records for the domain
func (s *RelayService) lookupMXRecords(domain string) ([]MXRecord, error) {
    mxs, err := net.LookupMX(domain)
    if err != nil {
        s.debugLog(fmt.Sprintf("MX lookup failed for %s, trying A record", domain))
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

    // Sort by MX preference
    sort.Slice(records, func(i, j int) bool {
        return records[i].Preference < records[j].Preference
    })

    return records, nil
}

// relayToDomain relays email to the MX server for the target domain
func (s *RelayService) relayToDomain(domain, from string, recipients []string, data []byte) error {
    mxRecords, err := s.lookupMXRecords(domain)
    if err != nil {
        return fmt.Errorf("MX lookup failed for domain %s: %w", domain, err)
    }

    var lastError error
    for _, mx := range mxRecords {
        for attempt := 0; attempt < s.config.RetryAttempts; attempt++ {
            if attempt > 0 {
                backoff := s.config.RetryDelay * time.Duration(attempt)
                s.debugLog(fmt.Sprintf("Retrying delivery to %s after %v", mx.Host, backoff))
                time.Sleep(backoff)
            }

            err = s.sendToMXServer(mx.Host, from, recipients, data)
            if err == nil {
                s.debugLog(fmt.Sprintf("Successfully delivered to MX %s", mx.Host))
                return nil
            }
            lastError = err
            s.debugLog(fmt.Sprintf("Attempt %d: Failed to deliver to MX %s: %v", attempt+1, mx.Host, err))
        }
    }

    return fmt.Errorf("failed to deliver to any MX server after retries: %w", lastError)
}

// sendToMXServer connects and relays the email to the MX server
func (s *RelayService) sendToMXServer(mxHost, from string, to []string, data []byte) error {
    var lastErr error

    for _, port := range s.config.PreferredPorts {
        addr := fmt.Sprintf("%s:%s", mxHost, port)
        s.debugLog(fmt.Sprintf("Attempting connection to %s", addr))

        dialer := &net.Dialer{
            Timeout: s.config.DialTimeout,
        }

        var conn net.Conn
        var err error

        if port == "465" {
            // Direct TLS connection
            tlsConfig := s.config.TLSConfig.Clone()
            tlsConfig.ServerName = mxHost
            conn, err = tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
        } else {
            conn, err = dialer.Dial("tcp", addr)
        }

        if err != nil {
            lastErr = fmt.Errorf("connection failed on port %s: %w", port, err)
            s.debugLog(fmt.Sprintf("Connection failed on port %s: %v", port, err))
            continue
        }

        success, err := s.handleSMTPSession(conn, mxHost, from, to, data, port)
        if success {
            return nil
        }
        lastErr = err
    }

    return lastErr
}

// handleSMTPSession manages the SMTP session for sending an email
func (s *RelayService) handleSMTPSession(conn net.Conn, mxHost, from string, to []string, data []byte, port string) (bool, error) {
    defer conn.Close()

    c, err := smtp.NewClient(conn, mxHost)
    if err != nil {
        return false, fmt.Errorf("SMTP client creation failed: %w", err)
    }
    defer c.Close()

    hostname := s.getLocalHostname()
    if err = c.Hello(hostname); err != nil {
        return false, fmt.Errorf("HELO failed: %w", err)
    }

    if port != "465" {
        if ok, _ := c.Extension("STARTTLS"); ok {
            tlsConfig := s.config.TLSConfig.Clone()
            tlsConfig.ServerName = mxHost
            if err = c.StartTLS(tlsConfig); err != nil {
                s.debugLog(fmt.Sprintf("TLS upgrade failed: %v", err))
            } else {
                s.debugLog("TLS upgrade successful")
            }
        }
    }

    if err = c.Mail(from); err != nil {
        return false, fmt.Errorf("MAIL FROM failed: %w", err)
    }

    successfulRcpts := 0
    for _, addr := range to {
        if err = c.Rcpt(addr); err != nil {
            s.debugLog(fmt.Sprintf("RCPT TO failed for %s: %v", addr, err))
            continue
        }
        successfulRcpts++
    }

    if successfulRcpts == 0 {
        return false, fmt.Errorf("no valid recipients")
    }

    w, err := c.Data()
    if err != nil {
        return false, fmt.Errorf("DATA command failed: %w", err)
    }

    if _, err = w.Write(data); err != nil {
        w.Close()
        return false, fmt.Errorf("data write failed: %w", err)
    }

    if err = w.Close(); err != nil {
        return false, fmt.Errorf("data close failed: %w", err)
    }

    if err = c.Quit(); err != nil {
        return false, fmt.Errorf("QUIT failed: %w", err)
    }

    return true, nil
}

func (s *RelayService) getLocalHostname() string {
    hostname, err := os.Hostname()
    if err != nil {
        return "localhost"
    }
    return hostname
}

func (s *RelayService) debugLog(message string) {
    if s.config.Debug {
        fmt.Printf("[RELAY DEBUG] %s\n", message)
    }
}