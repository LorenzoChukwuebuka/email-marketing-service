package services

import (
	"fmt"
	"github.com/emersion/go-imap/client"
	"net"
	"time"
)

type EmailVerifier struct{}

func NewEmailVerifier() *EmailVerifier {
	return &EmailVerifier{}
}

func (ev *EmailVerifier) VerifyDomainEmail(domain string) (bool, error) {
	// Step 1: Verify MX records
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return false, fmt.Errorf("failed to lookup MX records: %v", err)
	}
	if len(mxRecords) == 0 {
		return false, fmt.Errorf("no MX records found for domain %s", domain)
	}

	// Step 2: Attempt IMAP connection
	for _, mx := range mxRecords {
		imapServer := mx.Host
		canConnect, err := ev.checkIMAPConnection(imapServer)
		if err != nil {
			fmt.Printf("Error checking IMAP connection for %s: %v\n", imapServer, err)
			continue
		}
		if canConnect {
			return true, nil
		}
	}

	return false, fmt.Errorf("unable to establish IMAP connection for domain %s", domain)
}

func (ev *EmailVerifier) checkIMAPConnection(server string) (bool, error) {
	// Try common IMAP ports
	ports := []string{"993", "143"}

	for _, port := range ports {
		address := fmt.Sprintf("%s:%s", server, port)

		// Set a timeout for the connection attempt
		dialer := net.Dialer{Timeout: 10 * time.Second}

		conn, err := dialer.Dial("tcp", address)
		if err != nil {
			fmt.Printf("Failed to connect to %s: %v\n", address, err)
			continue
		}
		defer conn.Close()

		c, err := client.New(conn)
		if err != nil {
			fmt.Printf("Failed to create IMAP client for %s: %v\n", address, err)
			continue
		}

		// Don't attempt to login, just check if we can connect
		err = c.Logout()
		if err != nil {
			fmt.Printf("Failed to logout from IMAP server %s: %v\n", address, err)
		}

		return true, nil
	}

	return false, fmt.Errorf("unable to establish IMAP connection to server %s", server)
}
