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

	s.Addr = config.SMTP_PORT
	s.Domain = config.SMTP_SERVER
	s.WriteTimeout = 600 * time.Second
	s.ReadTimeout = 600 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	// Determine the server mode (development/production)
	mode := os.Getenv("SERVER_MODE")

	// Enable TLS only in production mode
	if mode == "production" {
		// cert, err := tls.LoadX509KeyPair("/etc/letsencrypt/live/smtp.crabmailer.com/fullchain.pem", "/etc/letsencrypt/live/smtp.crabmailer.com/privkey.pem")
		// if err != nil {
		// 	log.Fatalf("Failed to load TLS certificate: %v", err)
		// }

		// s.TLSConfig = &tls.Config{
		// 	Certificates: []tls.Certificate{cert},
		// 	MinVersion:   tls.VersionTLS12, // Ensure that TLS 1.2 or higher is used
		// }

		// Secure SMTP (SMTPS) port 465
		go func() {
			log.Println("Starting secure SMTP server on port 465")
			errChan := make(chan error, 1)
			go func() {
				errChan <- s.ListenAndServeTLS()
			}()
		}()
	} else {
		log.Println("Running in development mode: Insecure SMTP server on port", s.Addr)
	}

	// Start the insecure SMTP server (Port 1025 for testing)
	log.Println("Starting SMTP server at", s.Addr)

	errChan := make(chan error, 1)
	go func() {
		errChan <- s.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Println("Shutting down SMTP server...")
		return s.Close()
	case err := <-errChan:
		if err != smtp.ErrServerClosed {
			return utils.TraceError(fmt.Errorf("SMTP server error:%v", err))
		}
		return nil
	}
}
