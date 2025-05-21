package services

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"email-marketing-service/core/handler/admin/systems/dto"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/helper"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Service struct {
	store db.Store
}

func NewAdminSystemsService(store db.Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GenerateAndSaveSMTPCredentials(ctx context.Context, domain string) (*dto.SystemsResponse, error) {
	//check if domain already exists
	_, err := s.store.GetSMTPSettingByDomain(ctx, sql.NullString{String: domain, Valid: true})

	// If it's not a "not found" error, it's a real error we should return
	if !errors.Is(err, sql.ErrNoRows) && !strings.Contains(err.Error(), "no rows") {
		return nil, fmt.Errorf("error checking domain existence: %w", err)
	}

	if err == nil {
		return nil, fmt.Errorf("domain already exists")
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key pair: %w", err)
	}

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	privateKeyString := string(pem.EncodeToMemory(privateKeyPEM))

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}

	publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKeyBytes)
	formattedPublicKey := formatPublicKeyForDKIM(publicKeyBase64)

	selector := fmt.Sprintf("dkim%d", privateKey.N.Int64()%1000000)
	txtRecord := fmt.Sprintf("v=DKIM1; k=rsa; p=%s", formattedPublicKey)
	spfRecord := "v=spf1 mx -all"
	dmarcRecord := fmt.Sprintf("v=DMARC1; p=none; rua=mailto:postmaster@%s", domain)
	mxRecord := fmt.Sprintf("%s. 10 %s.", domain, domain)

	smtpSetting := &dto.SystemsResponse{
		Domain:         domain,
		TXTRecord:      txtRecord,
		DMARCRecord:    dmarcRecord,
		DKIMSelector:   selector,
		DKIMPublicKey:  publicKeyBase64,
		DKIMPrivateKey: privateKeyString,
		SPFRecord:      spfRecord,
		MXRecord:       mxRecord,
		Verified:       true,
	}

	_, err = s.store.CreateSystemsSMTPSettings(ctx, db.CreateSystemsSMTPSettingsParams{
		TxtRecord:      sql.NullString{String: txtRecord, Valid: true},
		DmarcRecord:    sql.NullString{String: dmarcRecord, Valid: true},
		DkimSelector:   sql.NullString{String: selector, Valid: true},
		DkimPublicKey:  sql.NullString{String: publicKeyBase64, Valid: true},
		DkimPrivateKey: sql.NullString{String: privateKeyString, Valid: true},
		SpfRecord:      sql.NullString{String: spfRecord, Valid: true},
		Verified:       sql.NullBool{Bool: true, Valid: true},
		MxRecord:       sql.NullString{String: mxRecord, Valid: true},
		Domain:         sql.NullString{String: domain, Valid: true},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to save SMTP settings to database: %w", err)
	}
	// In GenerateAndSaveSMTPCredentials
	if err := helper.SaveSMTPSettingsToFile(smtpSetting); err != nil {
		return nil, fmt.Errorf("failed to save SMTP settings to file: %w", err)
	}

	return smtpSetting, nil
}

func formatPublicKeyForDKIM(publicKey string) string {
	// Remove any newlines
	publicKey = strings.ReplaceAll(publicKey, "\n", "")

	// Split the key into chunks of 255 characters
	var chunks []string
	for len(publicKey) > 255 {
		chunks = append(chunks, publicKey[:255])
		publicKey = publicKey[255:]
	}
	chunks = append(chunks, publicKey) // Add the last remaining part

	// Join the chunks with double quotes and spaces
	return fmt.Sprintf("\"%s\"", strings.Join(chunks, "\" \""))
}

func (s *Service) GetDNSRecords(ctx context.Context, domain string) (map[string]string, error) {
	smtpSetting, err := s.store.GetSMTPSettingByDomain(ctx, sql.NullString{String: domain, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve SMTP settings: %w", err)
	}

	records := map[string]string{
		"TXT (DKIM)":  fmt.Sprintf("%s._domainkey.%s TXT %s", smtpSetting.DkimSelector.String, domain, smtpSetting.TxtRecord.String),
		"TXT (SPF)":   fmt.Sprintf("%s TXT %s", domain, smtpSetting.SpfRecord.String),
		"TXT (DMARC)": fmt.Sprintf("_dmarc.%s TXT %s", domain, smtpSetting.DmarcRecord.String),
		"MX":          smtpSetting.MxRecord.String,
	}

	return records, nil
}

func (s *Service) DeleteDNSRecords(ctx context.Context, domain string) error {
	// Add validation for domain parameter
	if domain == "" {
		return fmt.Errorf("domain cannot be empty")
	}

	// Log or debug the domain value
	log.Printf("Attempting to delete DNS records for domain: %s", domain)

	// Create the file path
	var dir string
	if os.Getenv("SERVER_MODE") == "production" {
		dir = "/app/backend/smtp_settings"
	} else {
		dir = "./smtp_settings"
	}

	filePath := filepath.Join(dir, fmt.Sprintf("%s_smtp_settings.json", domain))
	// Check if the file exists
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Even if file doesn't exist, we should still try to delete DB record
			return s.store.DeleteSystemsSMTPSetting(ctx, sql.NullString{String: domain, Valid: true})
		}
		return fmt.Errorf("failed to check file existence: %w", err)
	}

	// Remove the file
	err = os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to delete SMTP settings file: %w", err)
	}

	log.Printf("Successfully deleted file for domain: %s, proceeding with DB deletion", domain)

	// Call repository method and capture error
	if err := s.store.DeleteSystemsSMTPSetting(ctx, sql.NullString{String: domain, Valid: true}); err != nil {

		return fmt.Errorf("failed to delete settings in database: %w", err)
	}

	return nil
}
