package smtp_server

import (
	"context"
	//"crypto/tls"
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/utils"
	"fmt"
	"log"
	//"os"
	"github.com/emersion/go-smtp"
	"gorm.io/gorm"
	"time"
)

// StartSMTPServer starts the SMTP server and listens for shutdown signals
func StartSMTPServer(ctx context.Context, db *gorm.DB) error {

	smtpKeyRepo := repository.NewSMTPkeyRepository(db)

	be := NewBackend(smtpKeyRepo)
	s := smtp.NewServer(be)

	config := utils.LoadEnv()

	// mode := os.Getenv("SERVER_MODE")

	// var address string

	// if mode == "" {
	// 	address = config.SMTP_PORT
	// } else {
	// 	address = config.SMTP_PORT
	// }

	s.Addr = config.SMTP_PORT
	s.Domain = config.SMTP_SERVER
	s.WriteTimeout = 600 * time.Second
	s.ReadTimeout = 600 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	// cert, err := tls.LoadX509KeyPair("/etc/letsencrypt/live/yourdomain.com/fullchain.pem", "/etc/letsencrypt/live/yourdomain.com/privkey.pem")
	// if err != nil {
	// 	log.Fatalf("Failed to load TLS certificate: %v", err)
	// }
	// s.TLSConfig = &tls.Config{
	// 	Certificates: []tls.Certificate{cert},
	// }

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
