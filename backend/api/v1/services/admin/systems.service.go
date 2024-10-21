package adminservice

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	adminmodel "email-marketing-service/api/v1/model/admin"
	adminrepository "email-marketing-service/api/v1/repository/admin"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SystemsService struct {
	SystemsRepo *adminrepository.SystemRepository
}

func NewSystemsService(systemRepo *adminrepository.SystemRepository) *SystemsService {
	return &SystemsService{
		SystemsRepo: systemRepo,
	}
}

func (s *SystemsService) GenerateAndSaveSMTPCredentials(domain string) (*adminmodel.SystemsSMTPSetting, error) {
	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key pair: %w", err)
	}

	// Encode private key to PEM
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	privateKeyString := string(pem.EncodeToMemory(privateKeyPEM))

	// Encode public key
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}

	// Base64 encode the public key and format it for DKIM
	publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKeyBytes)
	formattedPublicKey := formatPublicKeyForDKIM(publicKeyBase64)

	// Generate DKIM selector with domain and timestamp
	selector := fmt.Sprintf("dkim%d", privateKey.N.Int64()%1000000)

	// Generate DNS records
	txtRecord := fmt.Sprintf("v=DKIM1; k=rsa; p=%s", formattedPublicKey)
	spfRecord := "v=spf1 mx -all"
	dmarcRecord := fmt.Sprintf("v=DMARC1; p=none; rua=mailto:postmaster@%s", domain)
	mxRecord := fmt.Sprintf("%s. 10 mail.%s.", domain, domain)

	// Create SystemsSMTPSetting
	smtpSetting := &adminmodel.SystemsSMTPSetting{
		Domain:         domain,
		TXTRecord:      txtRecord,
		DMARCRecord:    dmarcRecord,
		DKIMSelector:   selector,
		DKIMPublicKey:  publicKeyBase64,
		DKIMPrivateKey: privateKeyString,
		SPFRecord:      spfRecord,
		MXRecord:       mxRecord,
		Verified:       false,
	}

	// Save to database
	if err := s.SystemsRepo.CreateSMTPSettings(smtpSetting); err != nil {
		return nil, fmt.Errorf("failed to save SMTP settings to database: %w", err)
	}

	// Save to file
	if err := saveSMTPSettingsToFile(smtpSetting); err != nil {
		return nil, fmt.Errorf("failed to save SMTP settings to file: %w", err)
	}

	return smtpSetting, nil
}

// Function to format the public key for DKIM by splitting it into 253-character chunks
func formatPublicKeyForDKIM(publicKey string) string {
	// Remove any newlines
	publicKey = strings.ReplaceAll(publicKey, "\n", "")

	// Split the key into chunks of 253 characters (DNS TXT record limit)
	var chunks []string
	for len(publicKey) > 253 {
		chunks = append(chunks, publicKey[:253])
		publicKey = publicKey[253:]
	}
	chunks = append(chunks, publicKey)

	// Join the chunks with double quotes and spaces, ensuring no backslashes
	return "\"" + strings.Join(chunks, "\" \"") + "\""
}

func saveSMTPSettingsToFile(smtpSetting *adminmodel.SystemsSMTPSetting) error {
	// Create a map to store the settings
	settingsMap := map[string]interface{}{
		"Domain":         smtpSetting.Domain,
		"TXTRecord":      smtpSetting.TXTRecord,
		"DMARCRecord":    smtpSetting.DMARCRecord,
		"DKIMSelector":   smtpSetting.DKIMSelector,
		"DKIMPublicKey":  smtpSetting.DKIMPublicKey,
		"DKIMPrivateKey": smtpSetting.DKIMPrivateKey,
		"SPFRecord":      smtpSetting.SPFRecord,
		"MXRecord":       smtpSetting.MXRecord,
	}

	// Convert the map to JSON
	jsonData, err := json.MarshalIndent(settingsMap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal SMTP settings to JSON: %w", err)
	}

	// Create the directory if it doesn't exist
	dir := "./smtp_settings"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create the file path
	filePath := filepath.Join(dir, fmt.Sprintf("%s_smtp_settings.json", smtpSetting.Domain))

	// Write the JSON data to the file
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write SMTP settings to file: %w", err)
	}

	return nil
}

func (s *SystemsService) GetDNSRecords(domain string) (map[string]string, error) {
	smtpSetting, err := s.SystemsRepo.GetSMTPSettings(domain)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve SMTP settings: %w", err)
	}

	records := map[string]string{
		"TXT (DKIM)":  fmt.Sprintf("%s._domainkey.%s TXT %s", smtpSetting.DKIMSelector, domain, smtpSetting.TXTRecord),
		"TXT (SPF)":   fmt.Sprintf("%s TXT %s", domain, smtpSetting.SPFRecord),
		"TXT (DMARC)": fmt.Sprintf("_dmarc.%s TXT %s", domain, smtpSetting.DMARCRecord),
		"MX":          smtpSetting.MXRecord,
	}

	return records, nil
}

func (s *SystemsService) DeleteDNSRecords(domain string) error {
	// Create the file path
	dir := "./smtp_settings"
	filePath := filepath.Join(dir, fmt.Sprintf("%s_smtp_settings.json", domain))

	// Check if the file exists
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, so we skip without error
			return nil
		}
		// For any other error, return it
		return fmt.Errorf("failed to check file existence: %w", err)
	}

	// Remove the file
	err = os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to delete SMTP settings file: %w", err)
	}

	return s.SystemsRepo.DeleteSettings(domain)
}
