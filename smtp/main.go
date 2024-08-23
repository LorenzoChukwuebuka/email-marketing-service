package main

import (
	"crabmail/smtp/server"
	"log"
	"time"

	"github.com/emersion/go-smtp"
)

func main() {
	be := &server.Backend{}
	s := smtp.NewServer(be)

	// Configure server settings
	s.Addr = "localhost:1025"
	s.Domain = "localhost"
	s.WriteTimeout = 600 * time.Second
	s.ReadTimeout = 600 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	log.Println("Starting server at", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
