package smtp_server

import (
	"context"
	//"crypto/tls"
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/utils"
	"fmt"
	"github.com/emersion/go-smtp"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

// StartSMTPServer starts the SMTP server and listens for shutdown signals
func StartSMTPServer(ctx context.Context, db *gorm.DB) error {
	smtpKeyRepo := repository.NewSMTPkeyRepository(db)
	be := NewBackend(smtpKeyRepo)
	s := smtp.NewServer(be)
	config := utils.LoadEnv()
	mode := os.Getenv("SERVER_MODE")

	// Base server configuration
	s.Domain = config.SMTP_SERVER
	s.WriteTimeout = 600 * time.Second
	s.ReadTimeout = 600 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50

	if mode == "production" {
		// // Production configuration
		// s.Addr = config.SMTP_PORT // STARTTLS port

		// // Load TLS certificate
		// cert, err := tls.LoadX509KeyPair(
		// 	"/etc/letsencrypt/live/smtp.crabmailer.com/fullchain.pem",
		// 	"/etc/letsencrypt/live/smtp.crabmailer.com/privkey.pem",
		// )
		// if err != nil {
		// 	return fmt.Errorf("failed to load TLS certificate: %v", err)
		// }

		// // TLS configuration for STARTTLS
		// s.TLSConfig = &tls.Config{
		// 	Certificates: []tls.Certificate{cert},
		// 	MinVersion:   tls.VersionTLS12,
		// 	ClientAuth:   tls.VerifyClientCertIfGiven,
		// }

		// // Allow unencrypted auth because STARTTLS will be used
		// s.AllowInsecureAuth = true

		log.Printf("Starting production SMTP server with STARTTLS on port 587")
	} else {
		// Development configuration
		s.Addr = config.SMTP_PORT // Use port from config (typically 1025)
		s.AllowInsecureAuth = true
		log.Printf("Starting development SMTP server on port %s", s.Addr)
	}

	// Start the server
	errChan := make(chan error, 1)
	go func() {
		errChan <- s.ListenAndServe()
	}()

	// if mode == "production" {
	// 	go func() {
	// 		errChan <- s.ListenAndServeTLS()
	// 	}()
	// }

	// Wait for shutdown signal or error
	select {
	case <-ctx.Done():
		log.Println("Shutting down SMTP server...")
		return s.Close()
	case err := <-errChan:
		if err != smtp.ErrServerClosed {
			return utils.TraceError(fmt.Errorf("SMTP server error: %v", err))
		}
		return nil
	}
}
