package imap_server

import (
	"email-marketing-service/api/v1/repository"
	"log"

	"github.com/emersion/go-imap/server"
)

// StartIMAPServer initializes and starts the IMAP server
func StartIMAPServer(smtpKeyRepo *repository.SMTPKeyRepository) error {
	be := NewBackend(smtpKeyRepo)
	s := server.New(be)
	s.Addr = ":1143" // TODO: Make this configurable
	// TODO: Implement TLS support
	// TODO: Add graceful shutdown mechanism
	log.Println("Starting IMAP server at", s.Addr)
	return s.ListenAndServe()
}
