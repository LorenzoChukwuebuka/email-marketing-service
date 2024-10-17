package utils

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"net"
	"github.com/toorop/go-dkim"
	"log"
)

func ValidatePrivateKeyComprehensive(privateKeyPEM string) error {
	// Decode PEM block
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return fmt.Errorf("failed to parse PEM block containing private key")
	}

	// Check key type
	if block.Type != "RSA PRIVATE KEY" {
		return fmt.Errorf("unexpected key type: %s", block.Type)
	}

	// Parse the private key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %v", err)
	}

	// Validate the key
	if err := privateKey.Validate(); err != nil {
		return fmt.Errorf("private key validation failed: %v", err)
	}

	// Check key size
	keySize := privateKey.N.BitLen()
	if keySize < 1024 {
		return fmt.Errorf("key size is too small: %d bits", keySize)
	}

	fmt.Printf("Private key validated successfully. Key size: %d bits\n", keySize)
	return nil
}

func VerifyDomainAndSelector(domain, selector string) error {
	// Construct the DNS record name
	recordName := fmt.Sprintf("%s._domainkey.%s", selector, domain)

	// Look up TXT records
	txtRecords, err := net.LookupTXT(recordName)
	if err != nil {
		return fmt.Errorf("failed to lookup TXT records: %v", err)
	}

	// Check if any record contains the expected content
	for _, record := range txtRecords {
		if strings.Contains(record, "v=DKIM1") {
			fmt.Printf("DKIM record found for %s\n", recordName)
			return nil
		}
	}

	return fmt.Errorf("no valid DKIM record found for %s", recordName)
}



func SignEmail(email *[]byte, domain, selector, privateKey string) ([]byte, error) {
	log.Println("Starting SignEmail function")

	// Validate private key
	err := ValidatePrivateKeyComprehensive(privateKey)
	if err != nil {
		return nil, fmt.Errorf("private key validation failed: %v", err)
	}

	// Verify domain and selector
	err = VerifyDomainAndSelector(domain, selector)
	if err != nil {
		return nil, fmt.Errorf("domain and selector verification failed: %v", err)
	}

	// Ensure email is not nil
	if email == nil {
		return nil, fmt.Errorf("email pointer is nil")
	}

	// Log email content length
	log.Printf("Email content length: %d bytes", len(*email))

	// Set DKIM options
	options := dkim.NewSigOptions()
	options.PrivateKey = []byte(privateKey)
	options.Domain = domain
	options.Selector = selector
	options.SignatureExpireIn = 3600
	options.BodyLength = 50
	options.Headers = []string{"from", "to", "subject"}
	options.AddSignatureTimestamp = true
	options.Canonicalization = "relaxed/relaxed"
	options.Algo = "rsa-sha256"
	log.Printf("DKIM options set: Domain=%s, Selector=%s, Algo=%s", domain, selector, options.Algo)

	// Sign the email
	log.Println("Attempting to sign the email")
	err = dkim.Sign(email, options)
	if err != nil {
		log.Printf("Error signing email: %v", err)
		return nil, fmt.Errorf("error signing email: %w", err)
	}
	log.Println("Email signed successfully")

	return *email, nil
}
