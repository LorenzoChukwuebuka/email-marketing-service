package smtp_server

import (
	"bytes"
	"email-marketing-service/api/v1/repository"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"io"
	"log"
	"net/mail"
	"strings"
)

// Debug flag to enable/disable debug logging
const Debug = true

// Backend implements SMTP server methods
type Backend struct {
	SMTPKeyRepo *repository.SMTPKeyRepository
}

func NewBackend(smtpKeyRepo *repository.SMTPKeyRepository) *Backend {
	return &Backend{
		SMTPKeyRepo: smtpKeyRepo,
	}
}

// NewSession is called when a new SMTP session is initiated
func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	debugLog("New session started")
	return &Session{smtpKeyRepo: bkd.SMTPKeyRepo}, nil
}

// Session represents an SMTP session
type Session struct {
	from        string
	to          []string
	message     strings.Builder
	authState   int
	username    string
	password    string
	smtpKeyRepo *repository.SMTPKeyRepository
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
	case "LOGIN":
		s.authState = 1
		return nil, nil
	case "PLAIN":
		return sasl.NewPlainServer(s.authenticateUser), nil
	default:
		return nil, errors.New("unsupported authentication mechanism")
	}
}

// authenticateUser verifies the provided credentials
func (s *Session) authenticateUser(identity, username, password string) error {
	debugLog(fmt.Sprintf("Authenticating user: %s", username))
	auth, err := s.smtpKeyRepo.GetSMTPMasterKeyUserAndPass(username, password)

	if err != nil {
		debugLog("First authentication attempt failed, retrying with mock function")

		// Retry with the mock function
		auth, err = s.smtpKeyRepo.GetSMTPKeyUserAndPass(username, password)

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

// Mail handles the MAIL FROM command
func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	debugLog(fmt.Sprintf("Mail from: %s", from))
	if from == "" {
		return errors.New("501 5.1.1 Syntax error in parameters or arguments")
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

// Data handles the DATA command and receives the email content
func (s *Session) Data(r io.Reader) error {
	debugLog("Receiving message data")
	if b, err := io.ReadAll(r); err != nil {
		return fmt.Errorf("error reading message data: %w", err)
	} else {
		s.message.Write(b)
		debugLog(fmt.Sprintf("Message received:\n%s", s.message.String()))
		log.Printf("Email processed:\nFrom: %s\nTo: %s\nMessage: %s\n", s.from, s.to, s.message.String())

		messageStr := s.message.String()
		// Check if the email is valid before proceeding
		if err := isValidEmail(s.from, s.to, messageStr); err != nil {
			return fmt.Errorf("invalid email: %w", err)
		}

		// Store the email for IMAP access
		username := strings.Split(s.from, "@")[0]
		mailbox := "INBOX"
		err = s.smtpKeyRepo.StoreEmail(username, mailbox, s.from, s.to, b)
		if err != nil {
			return fmt.Errorf("error storing email: %w", err)
		}

		err = s.smtpKeyRepo.MarkEmailAsDelivered(s.from, s.to)
		if err != nil {
			return fmt.Errorf("error logging delivery: %w", err)
		}

		// Handle bounce detection or email processing here

		if err := s.processEmailBounce(b); err != nil {
			return fmt.Errorf("error processing bounce: %w", err)
		}
	}
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

		auth, err := s.smtpKeyRepo.GetSMTPMasterKeyUserAndPass(s.username, s.password)

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
	if err := s.smtpKeyRepo.UpdateBounceStatus(recipientEmail, bounceType); err != nil {
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
