package smtp_server

import (
	"bytes"
	"context"
	"database/sql"
	db "email-marketing-service/internal/db/sqlc"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/mail"
	"strings"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

// Debug flag to enable/disable debug logging
const Debug = true

// Backend implements SMTP server methods
type Backend struct {
	store        db.Store
	ctx          context.Context
	SPFValidator Validator
	relayService *RelayService
	rateLimiter  *RateLimiter
}

func NewBackend(store db.Store, ctx context.Context) *Backend {
	return &Backend{
		store:        store,
		ctx:          ctx,
		SPFValidator: *New(DefaultConfig()),
		relayService: NewRelayService(nil),
		rateLimiter:  NewRateLimiter(DefaultRateLimiterConfig()),
	}
}

// NewSession is called when a new SMTP session is initiated
// func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
// 	debugLog("New session started")
// 	return &Session{smtpKeyRepo: bkd.SMTPKeyRepo}, nil
// }

// Session represents an SMTP session
type Session struct {
	from         string
	to           []string
	message      bytes.Buffer
	authState    int
	username     string
	password     string
	store        db.Store
	spfValidator *Validator
	remoteIP     string
	relayService *RelayService
	rateLimiter  *RateLimiter
	ctx          context.Context
}

// Update NewSession
func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	debugLog("New session started")
	remoteIP := strings.Split(c.Conn().RemoteAddr().String(), ":")[0]

	// Check connection rate limit
	if err := bkd.rateLimiter.CheckConnection(remoteIP); err != nil {
		debugLog(fmt.Sprintf("Rate limit exceeded for IP %s: %v", remoteIP, err))
		return nil, err
	}

	return &Session{
		store:        bkd.store,
		ctx:          bkd.ctx,
		spfValidator: &bkd.SPFValidator,
		remoteIP:     remoteIP,
		relayService: bkd.relayService,
		rateLimiter:  bkd.rateLimiter,
	}, nil
}

// AuthMechanisms returns the list of supported authentication mechanisms
func (s *Session) AuthMechanisms() []string {
	debugLog("Requested auth mechanisms")
	return []string{"LOGIN", "PLAIN"}
}

// Auth handles the authentication process
func (s *Session) Auth(mech string) (sasl.Server, error) {
	debugLog(fmt.Sprintf("Auth requested with mechanism: %s", mech))
	switch mech {
	// case "LOGIN":
	// 	s.authState = 1
	// 	return sasl.NewLoginAuthServer(s.handleLoginAuth), nil
	case "PLAIN":
		return sasl.NewPlainServer(s.authenticateUser), nil
	default:
		return nil, errors.New("unsupported authentication mechanism")
	}
}

// authenticateUser verifies the provided credentials
func (s *Session) authenticateUser(identity, username, password string) error {
	debugLog(fmt.Sprintf("Authenticating user: %s", username))

	auth, err := s.store.CheckSMTPMasterKeyExists(s.ctx, db.CheckSMTPMasterKeyExistsParams{SmtpLogin: username, Password: password})

	if err != nil {
		debugLog("First authentication attempt failed, retrying with mock function")

		// Retry with the mock function
		auth, err = s.store.CheckSMTPKeyExists(s.ctx, db.CheckSMTPKeyExistsParams{KeyName: username, Password: password})

		if err != nil {
			return errors.New("454 4.7.0 Temporary authentication failure")
		}
	}

	if !auth {
		return errors.New("535 5.7.8 Authentication credentials invalid")
	}

	debugLog("Authentication successful")
	return nil
}

func (s *Session) handleLoginAuth(username, password string) error {
	return s.authenticateUser("", username, password)
}

// Mail handles the MAIL FROM command
func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	debugLog(fmt.Sprintf("Mail from: %s using IP: %s", from, s.remoteIP))
	if from == "" {
		return errors.New("501 5.1.1 Syntax error in parameters or arguments")
	}

	// Extract domain from sender address
	parts := strings.Split(from, "@")
	if len(parts) != 2 {
		return errors.New("501 5.1.7 Invalid sender address")
	}
	domain := parts[1]

	// Perform SPF validation
	result, err := s.spfValidator.CheckHost(s.remoteIP, domain, from)
	if err != nil {
		debugLog(fmt.Sprintf("SPF error for %s: %v", from, err))
		return fmt.Errorf("450 4.4.4 Temporary SPF validation error: %v", err)
	}

	// Handle SPF result
	switch result {
	case Pass:
		debugLog(fmt.Sprintf("SPF passed for %s", from))
	case Fail:
		return errors.New("550 5.7.1 SPF authentication failed")
	case SoftFail:
		debugLog(fmt.Sprintf("SPF soft fail for %s", from))
		// You might want to accept or reject soft fails based on your policy
	case TempError:
		return errors.New("451 4.4.3 Temporary SPF lookup error")
	case PermError:
		return errors.New("550 5.5.2 SPF record could not be correctly interpreted")
	}

	s.from = from
	return nil
}

// Rcpt handles the RCPT TO command
func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	debugLog(fmt.Sprintf("Rcpt to: %s", to))
	if to == "" {
		return errors.New("501 5.1.1 Syntax error in parameters or arguments") // Return an SMTP error code
	}

	s.to = append(s.to, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	debugLog("Starting to process email data")

	// Rate limiting check with logging
	debugLog(fmt.Sprintf("Checking rate limits for IP: %s, From: %s", s.remoteIP, s.from))
	if err := s.rateLimiter.CheckMessage(s.remoteIP, s.from, s.to); err != nil {
		debugLog(fmt.Sprintf("Rate limit exceeded: %v", err))
		return fmt.Errorf("452 4.5.3 %v", err)
	}
	debugLog("Rate limit check passed")

	// Read message data with logging
	debugLog("Reading message data")
	b, err := io.ReadAll(r)
	if err != nil {
		debugLog(fmt.Sprintf("Error reading message data: %v", err))
		return fmt.Errorf("error reading message data: %w", err)
	}
	debugLog(fmt.Sprintf("Message data read successfully, size: %d bytes", len(b)))

	// Parse the message with logging
	debugLog("Parsing message")
	msg, err := mail.ReadMessage(bytes.NewReader(b))
	if err != nil {
		debugLog(fmt.Sprintf("Error parsing message: %v", err))
		return fmt.Errorf("error parsing message: %w", err)
	}
	debugLog("Message parsed successfully")

	// Extract subject with logging
	subject := msg.Header.Get("Subject")
	debugLog(fmt.Sprintf("Subject extracted: %s", subject))

	// Store in session buffer
	s.message.Write(b)
	debugLog(fmt.Sprintf("Message stored in buffer, size: %d bytes", s.message.Len()))

	// Validation with logging
	debugLog("Validating email")
	if err := isValidEmail(s.from, s.to, s.message.String()); err != nil {
		debugLog(fmt.Sprintf("Email validation failed: %v", err))
		return fmt.Errorf("invalid email: %w", err)
	}
	debugLog("Email validation passed")

	// Create email object for relay with logging
	debugLog("Creating email object for relay")
	email := &Email{
		From:    s.from,
		To:      s.to,
		Subject: subject,
		Body:    b,
		Headers: map[string]string{
			"X-Mailer":    "CustomSMTPServer",
			"X-Source-IP": s.remoteIP,
		},
	}
	debugLog(fmt.Sprintf("Email object created - From: %s, To: %v, Subject: %s",
		email.From, email.To, email.Subject))

	// Relay the email with detailed logging
	debugLog("Starting email relay process")
	if err := s.relayService.RelayEmail(email); err != nil {
		debugLog(fmt.Sprintf("Relay failed with detailed error: %+v", err))
		// Return immediately on relay failure
		return fmt.Errorf("550 5.7.0 Relay email failed: %w", err)
	}
	debugLog("Email relay completed successfully")

	// Store for IMAP access with logging
	debugLog("Storing email for IMAP access")
	username := strings.SplitN(s.from, "@", 2)[0]
	const mailbox = "INBOX"

	_, err = s.store.CreateEmailBox(s.ctx, db.CreateEmailBoxParams{UserName: sql.NullString{String: username, Valid: true}, To: sql.NullString{String: strings.Join(s.to, ","), Valid: true}, From: sql.NullString{String: s.from, Valid: true}, Mailbox: sql.NullString{String: mailbox}})
	if err != nil {
		debugLog(fmt.Sprintf("Error storing email: %v", err))
		return fmt.Errorf("error storing email: %w", err)
	}
	debugLog("Email stored successfully")

	// Mark as delivered with logging
	debugLog("Marking email as delivered")

	if err := s.MarkEmailAsDelivered(s.ctx, s.from, s.to); err != nil {
		debugLog(fmt.Sprintf("Error marking email as delivered: %v", err))
		return fmt.Errorf("error logging delivery: %w", err)
	}
	debugLog("Email marked as delivered")

	// Process email bounce handling with logging
	debugLog("Processing potential email bounce")
	if err := s.processEmailBounce(b); err != nil {
		debugLog(fmt.Sprintf("Error processing bounce: %v", err))
		return fmt.Errorf("error processing bounce: %w", err)
	}
	debugLog("Bounce processing completed")

	debugLog("Data processing completed successfully")
	return nil
}

// Reset resets the session state as per SMTP RSET command
func (s *Session) Reset() {
	debugLog("Resetting session")
	s.from = ""
	s.to = nil
	s.message.Reset()
	s.authState = 0
	s.username = ""
	s.password = ""
}

// Logout is called when the session is closed
func (s *Session) Logout() error {
	debugLog("User logged out")
	return nil
}

// Command handles custom commands, particularly for LOGIN auth
func (s *Session) Command(cmd string, args []string) (string, error) {
	debugLog(fmt.Sprintf("Command received: %s, args: %v", cmd, args))
	switch s.authState {
	case 1: // Waiting for username
		if len(args) == 0 {
			return "500 5.5.1 Syntax error, command unrecognized", nil
		}
		decoded, err := base64.StdEncoding.DecodeString(args[0])
		if err != nil {
			return "454 4.7.0 Invalid response", nil
		}
		s.username = string(decoded)
		s.authState = 2
		return "334 UGFzc3dvcmQ=", nil // "Password" in Base64
	case 2: // Waiting for password
		if len(args) == 0 {
			return "500 5.5.1 Syntax error, command unrecognized", nil
		}
		decoded, err := base64.StdEncoding.DecodeString(args[0])
		if err != nil {
			return "454 4.7.0 Invalid response", nil
		}
		s.password = string(decoded)
		s.authState = 0

		auth, err := s.store.CheckSMTPMasterKeyExists(s.ctx, db.CheckSMTPMasterKeyExistsParams{SmtpLogin: s.username, Password: s.password})

		if err != nil {
			debugLog(fmt.Sprintf("Error authenticating user: %v", err))
			return "454 4.7.0 Temporary authentication failure", nil
		}

		if auth {
			debugLog("Authentication successful")
			return "235 2.7.0 Authentication successful", nil
		}

		debugLog("Authentication failed")
		return "535 5.7.8 Authentication credentials invalid", nil
	default:
		return "500 5.5.1 Syntax error, command unrecognized", nil
	}
}

func (s *Session) processEmailBounce(emailContent []byte) error {
	// Parse the email content
	msg, err := mail.ReadMessage(bytes.NewReader(emailContent))
	if err != nil {
		return fmt.Errorf("error reading email message: %w", err)
	}

	// Use mail.Header to access email headers
	header := msg.Header
	bounceReason := header.Get("Diagnostic-Code")

	// Determine if it's a soft or hard bounce
	bounceType := "soft"
	if strings.Contains(strings.ToLower(bounceReason), "permanent") {
		bounceType = "hard"
	}

	// Extract the original recipient email from the headers or message body
	recipientEmail := header.Get("Original-Recipient")

	// Update your database with the bounce information
	log.Printf("Bounce detected for %s: %s (%s)", recipientEmail, bounceReason, bounceType)

	// Call your repository method to update the bounce status
	if err := s.store.UpdateBounceStatus(s.ctx, db.UpdateBounceStatusParams{RecipientEmail: recipientEmail, BounceStatus: sql.NullString{String: bounceType, Valid: true}}); err != nil {
		return fmt.Errorf("error updating bounce status: %w", err)
	}

	return nil
}

// debugLog prints debug information if Debug is true
func debugLog(message string) {
	if Debug {
		log.Printf("[DEBUG] %s", message)
	}
}

func isValidEmail(from string, to []string, message string) error {
	// Check if the 'From' field is valid
	if _, err := mail.ParseAddress(from); err != nil {
		return errors.New("invalid 'From' address")
	}

	// Check if there is at least one valid 'To' recipient
	if len(to) == 0 {
		return errors.New("no 'To' recipients")
	}

	for _, recipient := range to {
		if _, err := mail.ParseAddress(recipient); err != nil {
			return errors.New("invalid 'To' recipient: " + recipient)
		}
	}

	// Ensure the email message contains a subject (this is optional but recommended)
	if !strings.Contains(strings.ToLower(message), "subject:") {
		return errors.New("missing 'Subject' field in email")
	}

	// Check if the message body contains content
	bodyStartIndex := strings.Index(message, "\r\n\r\n")
	if bodyStartIndex == -1 || len(strings.TrimSpace(message[bodyStartIndex:])) == 0 {
		return errors.New("email body is empty")
	}

	return nil
}

func (s *Session) MarkEmailAsDelivered(ctx context.Context, from string, to []string) error {
	// Process each recipient individually
	for _, recipient := range to {
		if err := s.store.MarkEmailAsDelivered(ctx, recipient); err != nil {
			return fmt.Errorf("failed to mark email as delivered for %s: %w", recipient, err)
		}
	}
	return nil
}
