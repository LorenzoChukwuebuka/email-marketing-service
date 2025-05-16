package utils

import (
	adminmodel "email-marketing-service/api/v1/model/admin"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	smtpDirMutex sync.Mutex
	smtpDirPath  string
)

// initSMTPDirectory ensures the SMTP settings directory exists with proper permissions
func initSMTPDirectory() error {
	smtpDirMutex.Lock()
	defer smtpDirMutex.Unlock()

	if os.Getenv("SERVER_MODE") == "production" {
		smtpDirPath = "/app/backend/smtp_settings"
	} else {
		smtpDirPath = "./smtp_settings"
	}

	// Create directory with full permissions
	if err := os.MkdirAll(smtpDirPath, 0777); err != nil {
		return fmt.Errorf("failed to create SMTP settings directory: %w", err)
	}

	// Ensure directory has correct permissions even if it already existed
	if err := os.Chmod(smtpDirPath, 0777); err != nil {
		return fmt.Errorf("failed to set SMTP directory permissions: %w", err)
	}

	return nil
}

// GetSMTPFilePath returns the full path for an SMTP settings file
func GetSMTPFilePath(domain string) string {
	return filepath.Join(smtpDirPath, fmt.Sprintf("%s_smtp_settings.json", domain))
}

// SaveSMTPSettingsToFile saves SMTP settings to a file with proper permissions
func SaveSMTPSettingsToFile(smtpSetting *adminmodel.SystemsSMTPSetting) error {
	if err := initSMTPDirectory(); err != nil {
		return err
	}

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

	jsonData, err := json.MarshalIndent(settingsMap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal SMTP settings to JSON: %w", err)
	}

	filePath := GetSMTPFilePath(smtpSetting.Domain)

	// Write file with full permissions
	if err := os.WriteFile(filePath, jsonData, 0666); err != nil {
		return fmt.Errorf("failed to write SMTP settings to file: %w", err)
	}

	// Ensure file has correct permissions after writing
	if err := os.Chmod(filePath, 0666); err != nil {
		return fmt.Errorf("failed to set SMTP file permissions: %w", err)
	}

	return nil
}

// ReadSMTPSettingsFromFile reads SMTP settings from file with proper error handling
func ReadSMTPSettingsFromFile(domain string) (*adminmodel.SystemsSMTPSetting, error) {
	if err := initSMTPDirectory(); err != nil {
		return nil, err
	}

	filePath := GetSMTPFilePath(domain)

	// Check if file exists before trying to read
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("SMTP settings file does not exist for domain %s: %w", domain, err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SMTP settings file: %w", err)
	}

	var settingsMap map[string]interface{}
	if err := json.Unmarshal(data, &settingsMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal SMTP settings: %w", err)
	}

	smtpSetting := &adminmodel.SystemsSMTPSetting{
		Domain:         settingsMap["Domain"].(string),
		TXTRecord:      settingsMap["TXTRecord"].(string),
		DMARCRecord:    settingsMap["DMARCRecord"].(string),
		DKIMSelector:   settingsMap["DKIMSelector"].(string),
		DKIMPublicKey:  settingsMap["DKIMPublicKey"].(string),
		DKIMPrivateKey: settingsMap["DKIMPrivateKey"].(string),
		SPFRecord:      settingsMap["SPFRecord"].(string),
		MXRecord:       settingsMap["MXRecord"].(string),
	}

	return smtpSetting, nil
}
