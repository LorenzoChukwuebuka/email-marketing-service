// package smtp_server

// import (
// 	"context"
// 	"email-marketing-service/internal/common"
// 	"email-marketing-service/internal/config"
// 	db "email-marketing-service/internal/db/sqlc"
// 	//"crypto/tls"
// 	"fmt"
// 	"log"
// 	"os"
// 	"time"

// 	"github.com/emersion/go-smtp"
// )

// var (
// 	cfg = config.LoadEnv()
// )

// // StartSMTPServer starts the SMTP server and listens for shutdown signals
// func StartSMTPServer(ctx context.Context, store db.Store) error {
// 	be := NewBackend(store, ctx)
// 	s := smtp.NewServer(be)
// 	config := cfg
// 	mode := os.Getenv("SERVER_MODE")

// 	// Base server configuration
// 	s.Domain = config.SMTP_SERVER
// 	s.WriteTimeout = 600 * time.Second
// 	s.ReadTimeout = 600 * time.Second
// 	s.MaxMessageBytes = 1024 * 1024
// 	s.MaxRecipients = 50

// 	if mode == "production" {
// 		// // Production configuration
// 		// s.Addr = config.SMTP_PORT // STARTTLS port

// 		// // Load TLS certificate
// 		// cert, err := tls.LoadX509KeyPair(
// 		// 	"/etc/letsencrypt/live/smtp.crabmailer.com/fullchain.pem",
// 		// 	"/etc/letsencrypt/live/smtp.crabmailer.com/privkey.pem",
// 		// )
// 		// if err != nil {
// 		// 	return fmt.Errorf("failed to load TLS certificate: %v", err)
// 		// }

// 		// // TLS configuration for STARTTLS
// 		// s.TLSConfig = &tls.Config{
// 		// 	Certificates: []tls.Certificate{cert},
// 		// 	MinVersion:   tls.VersionTLS12,
// 		// 	ClientAuth:   tls.VerifyClientCertIfGiven,
// 		// }

// 		// // Allow unencrypted auth because STARTTLS will be used
// 		// s.AllowInsecureAuth = true

// 		log.Printf("Starting production SMTP server with STARTTLS on port 587")
// 	} else {
// 		// Development configuration
// 		s.Addr = config.SMTP_PORT // Use port from config (typically 1025)
// 		s.AllowInsecureAuth = true
// 		log.Printf("Starting development SMTP server on port %s", s.Addr)
// 	}

// 	// Start the server
// 	errChan := make(chan error, 1)
// 	go func() {
// 		errChan <- s.ListenAndServe()
// 	}()

// 	// if mode == "production" {
// 	// 	go func() {
// 	// 		errChan <- s.ListenAndServeTLS()
// 	// 	}()
// 	// }

// 	// Wait for shutdown signal or error
// 	select {
// 	case <-ctx.Done():
// 		log.Println("Shutting down SMTP server...")
// 		return s.Close()
// 	case err := <-errChan:
// 		if err != smtp.ErrServerClosed {
// 			return common.TraceError(fmt.Errorf("SMTP server error: %v", err))
// 		}
// 		return nil
// 	}
// }

package smtp_server

import (
	"context"
	"crypto/tls"
	"email-marketing-service/internal/common"
	"email-marketing-service/internal/config"
	db "email-marketing-service/internal/db/sqlc"
	"fmt"
	"log"
	"time"

	"github.com/emersion/go-smtp"
)

var (
	cfg = config.LoadEnv()
)

// SMTPServerConfig holds the configuration for the SMTP server
type SMTPServerConfig struct {
	Port              string
	Domain            string
	TLSEnabled        bool
	CertFile          string
	KeyFile           string
	AllowInsecureAuth bool
	WriteTimeout      time.Duration
	ReadTimeout       time.Duration
	MaxMessageBytes   int
	MaxRecipients     int
}

// SMTPServerOption is a functional option for configuring the SMTP server
type SMTPServerOption func(*SMTPServerConfig)

// WithPort sets the port for the SMTP server
func WithPort(port string) SMTPServerOption {
	return func(c *SMTPServerConfig) {
		c.Port = port
	}
}

// WithDomain sets the domain for the SMTP server
func WithDomain(domain string) SMTPServerOption {
	return func(c *SMTPServerConfig) {
		c.Domain = domain
	}
}

// WithTLS enables TLS with the provided certificate and key files
func WithTLS(certFile, keyFile string) SMTPServerOption {
	return func(c *SMTPServerConfig) {
		c.TLSEnabled = true
		c.CertFile = certFile
		c.KeyFile = keyFile
	}
}

// WithInsecureAuth allows insecure authentication (useful for development or STARTTLS)
func WithInsecureAuth(allow bool) SMTPServerOption {
	return func(c *SMTPServerConfig) {
		c.AllowInsecureAuth = allow
	}
}

// WithTimeouts sets the read and write timeouts
func WithTimeouts(read, write time.Duration) SMTPServerOption {
	return func(c *SMTPServerConfig) {
		c.ReadTimeout = read
		c.WriteTimeout = write
	}
}

// WithLimits sets the maximum message bytes and recipients
func WithLimits(maxMessageBytes, maxRecipients int) SMTPServerOption {
	return func(c *SMTPServerConfig) {
		c.MaxMessageBytes = maxMessageBytes
		c.MaxRecipients = maxRecipients
	}
}

// newDefaultConfig returns a default SMTP server configuration
func newDefaultConfig() *SMTPServerConfig {
	return &SMTPServerConfig{
		Port:              cfg.SMTP_PORT,
		Domain:            cfg.SMTP_SERVER,
		TLSEnabled:        false,
		AllowInsecureAuth: true,
		WriteTimeout:      600 * time.Second,
		ReadTimeout:       600 * time.Second,
		MaxMessageBytes:   1024 * 1024,
		MaxRecipients:     50,
	}
}

// StartSMTPServer starts the SMTP server with the provided options
func StartSMTPServer(ctx context.Context, store db.Store, opts ...SMTPServerOption) error {
	// Start with default configuration
	config := newDefaultConfig()

	// Apply all options
	for _, opt := range opts {
		opt(config)
	}

	be := NewBackend(store, ctx)
	s := smtp.NewServer(be)

	// Apply configuration to server
	s.Domain = config.Domain
	s.Addr = config.Port
	s.WriteTimeout = config.WriteTimeout
	s.ReadTimeout = config.ReadTimeout
	s.MaxMessageBytes = int64(config.MaxMessageBytes)
	s.MaxRecipients = config.MaxRecipients
	s.AllowInsecureAuth = config.AllowInsecureAuth

	// Configure TLS if enabled
	if config.TLSEnabled {
		cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
		if err != nil {
			return fmt.Errorf("failed to load TLS certificate: %v", err)
		}

		s.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
			ClientAuth:   tls.VerifyClientCertIfGiven,
		}

		log.Printf("Starting SMTP server with TLS on port %s", config.Port)
	} else {
		log.Printf("Starting SMTP server without TLS on port %s", config.Port)
	}

	// Start the server
	errChan := make(chan error, 1)
	go func() {
		if config.TLSEnabled {
			errChan <- s.ListenAndServeTLS()
		} else {
			errChan <- s.ListenAndServe()
		}
	}()

	// Wait for shutdown signal or error
	select {
	case <-ctx.Done():
		log.Println("Shutting down SMTP server...")
		return s.Close()
	case err := <-errChan:
		if err != smtp.ErrServerClosed {
			return common.TraceError(fmt.Errorf("SMTP server error: %v", err))
		}
		return nil
	}
}
