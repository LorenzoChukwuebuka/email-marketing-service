package server 

 import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
 

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

// Debug flag to enable/disable debug logging
const Debug = true

// Backend implements SMTP server methods
type Backend struct{}

// NewSession is called when a new SMTP session is initiated
func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	debugLog("New session started")
	return &Session{}, nil
}

// Session represents an SMTP session
type Session struct {
	from      string
	to        []string
	message   strings.Builder
	authState int
	username  string
	password  string
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

	// TODO: Replace with database lookup
	if username != "username" || password != "password" {
		return errors.New("invalid username or password")
	}

	debugLog("Authentication successful")
	return nil
}

// Mail handles the MAIL FROM command
func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	debugLog(fmt.Sprintf("Mail from: %s", from))
	s.from = from
	return nil
}

// Rcpt handles the RCPT TO command
func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	debugLog(fmt.Sprintf("Rcpt to: %s", to))
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
			return "500 Syntax error, command unrecognized", nil
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
			return "500 Syntax error, command unrecognized", nil
		}
		decoded, err := base64.StdEncoding.DecodeString(args[0])
		if err != nil {
			return "454 4.7.0 Invalid response", nil
		}
		s.password = string(decoded)
		s.authState = 0

		// TODO: Replace with database lookup
		if s.username == "username" && s.password == "password" {
			debugLog("Authentication successful")
			return "235 2.7.0 Authentication successful", nil
		}
		debugLog("Authentication failed")
		return "535 5.7.8 Authentication credentials invalid", nil
	default:
		return "", nil
	}
}

// debugLog prints debug information if Debug is true
func debugLog(message string) {
	if Debug {
		log.Printf("[DEBUG] %s", message)
	}
}
