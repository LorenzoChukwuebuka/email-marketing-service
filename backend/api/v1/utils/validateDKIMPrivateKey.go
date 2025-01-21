package utils

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func ValidatePrivateKey(privateKey string) error {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return fmt.Errorf("failed to parse PEM block containing the private key")
	}

	if block.Type != "RSA PRIVATE KEY" {
		return fmt.Errorf("not an RSA private key")
	}

	_, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %v", err)
	}

	return nil
}
