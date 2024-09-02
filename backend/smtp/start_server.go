package smtp_server

import (
	"context"
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/utils"
	"log"
	"time"

	"github.com/emersion/go-smtp"
	"gorm.io/gorm"
)

// StartSMTPServer starts the SMTP server and listens for shutdown signals
func StartSMTPServer(ctx context.Context, db *gorm.DB) error {

	smtpKeyRepo := repository.NewSMTPkeyRepository(db)

	be := NewBackend(smtpKeyRepo)
	s := smtp.NewServer(be)

	config := utils.LoadEnv()

	s.Addr =  config.SMTP_PORT
	s.Domain = config.SMTP_SERVER
	s.WriteTimeout = 600 * time.Second
	s.ReadTimeout = 600 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

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
			return err
		}
		return nil
	}
}
