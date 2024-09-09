package services

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/utils"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"net"
	"strings"
	"text/template"
)

type DNSRecord struct {
	Type       string
	RecordName string
	Value      string
}

type DomainService struct {
	DomainRepo *repository.DomainRepository
}

func NewDomainService(domainRepo *repository.DomainRepository) *DomainService {
	return &DomainService{
		DomainRepo: domainRepo,
	}
}

func (s *DomainService) CreateDomain(d *dto.DomainDTO) (map[string]interface{}, error) {
	if err := utils.ValidateData(d); err != nil {
		return nil, fmt.Errorf("invalid plan data: %w", err)
	}

	domainModel := &model.Domains{
		UUID:   uuid.New().String(),
		UserID: d.UserId,
		Domain: d.Domain,
	}

	isDomainExists, err := s.DomainRepo.CheckIfDomainExists(domainModel)
	if err != nil {
		return nil, err
	}

	if isDomainExists {
		return nil, fmt.Errorf("domain already exists")
	}

	verifyDomain := s.verifyDomain(d.Domain)
	if !verifyDomain {
		return nil, fmt.Errorf("your domain is not verified")
	}

	// Generate DNS records
	txtRecord := s.generateTXTRecord(d.Domain)
	dmarcRecord := s.generateDMARCRecord(d.Domain)
	dkimSelector := s.generateDKIMSelector()
	dkimPublicKey, dkimPrivateKey, err := s.generateDKIMKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to generate DKIM keys: %v", err)
	}

	// Create structured DNS records
	dnsRecords := []DNSRecord{
		{
			Type:       "TXT",
			RecordName: "@ (or leave it empty)",
			Value:      txtRecord,
		},
		{
			Type:       "TXT",
			RecordName: "_dmarc",
			Value:      dmarcRecord,
		},
		{
			Type:       "TXT",
			RecordName: fmt.Sprintf("%s._domainkey", dkimSelector),
			Value:      fmt.Sprintf("v=DKIM1; k=rsa; p=%s", dkimPublicKey),
		},
	}

	// Generate downloadable content
	downloadContent, err := s.generateDownloadableRecords(d.Domain, dnsRecords)
	if err != nil {
		return nil, fmt.Errorf("failed to generate downloadable records: %v", err)
	}

	// Update domain model
	domainModel.TXTRecord = txtRecord
	domainModel.DMARCRecord = dmarcRecord
	domainModel.DKIMSelector = dkimSelector
	domainModel.DKIMPublicKey = dkimPublicKey
	domainModel.DKIMPrivateKey = dkimPrivateKey

	if err := s.DomainRepo.CreateDomain(domainModel); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"domain":               d.Domain,
		"dns_records":          dnsRecords,
		"downloadable_records": downloadContent,
	}, nil
}

func (s *DomainService) verifyDomain(domain string) bool {
	ips, err := net.LookupIP(domain)
	if err != nil {
		fmt.Printf("Error looking up IP for %s: %v\n", domain, err)
		return false
	}
	fmt.Printf("IP addresses for %s: %v\n", domain, ips)

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		fmt.Printf("Warning: Error looking up MX records for %s: %v\n", domain, err)
	} else {
		fmt.Printf("MX records for %s:\n", domain)
		for _, mx := range mxRecords {
			fmt.Printf("  %v (priority: %v)\n", mx.Host, mx.Pref)
		}
	}
	return true
}

func (s *DomainService) generateTXTRecord(domain string) string {

	// Generate a unique hash based on the domain and a secret key
	hash := sha256.Sum256([]byte(domain + config.ENC_KEY))

	// Convert the first 16 bytes of the hash to base64
	verificationCode := base64.URLEncoding.EncodeToString(hash[:16])

	// Remove any non-alphanumeric characters
	verificationCode = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return -1
	}, verificationCode)

	// Truncate to a reasonable length (e.g., 20 characters)
	if len(verificationCode) > 20 {
		verificationCode = verificationCode[:20]
	}

	return fmt.Sprintf("email-verify=%s", verificationCode)
	//return fmt.Sprintf("your-email-service-verification=%s", domain)
}

func (s *DomainService) generateDMARCRecord(domain string) string {
	return fmt.Sprintf("v=DMARC1; p=none; rua=mailto:dmarc-reports@%s", domain)
}

func (s *DomainService) generateDKIMSelector() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("key%x", b)
}

func (s *DomainService) generateDKIMKeys() (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 500)
	if err != nil {
		return "", "", err
	}

	publicKeyDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}

	publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKeyDER)
	privateKeyPEM := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBase64 := base64.StdEncoding.EncodeToString(privateKeyPEM)

	return publicKeyBase64, privateKeyBase64, nil
}

func (s *DomainService) generateDownloadableRecords(domain string, records []DNSRecord) (string, error) {
	tmpl := `Domain: {{.Domain}}

		DNS Records to add:

		{{range .Records}}
		Type:        {{.Type}}
		Record name: {{.RecordName}}
		Value:       {{.Value}}

		{{end}}
		Please add these records to your domain's DNS settings. If you need help, contact your domain registrar or DNS provider.
		`

	data := struct {
		Domain  string
		Records []DNSRecord
	}{
		Domain:  domain,
		Records: records,
	}

	t, err := template.New("dns_records").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *DomainService) InitiateVerification(domainID string) (bool, error) {
	domain, err := s.DomainRepo.GetDomain(domainID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve domain: %v", err)
	}

	if domain.Verified {
		return false, fmt.Errorf("domain has been verified")
	}

	verified, err := s.VerifyDNSRecords(domain)
	if err != nil {
		return false, err
	}

	domainModel := &model.Domains{
		UUID: domain.UUID,
	}

	if verified {
		domainModel.Verified = true
		err = s.DomainRepo.UpdateDomain(domainModel)
		if err != nil {
			return false, fmt.Errorf("failed to update domain status: %v", err)
		}
	}

	return verified, nil
}

func (s *DomainService) VerifyDNSRecords(domain *model.DomainsResponse) (bool, error) {
	txtVerified, err := s.verifyTXTRecord(domain.Domain, domain.TXTRecord)
	if err != nil {
		return false, err
	}

	if !txtVerified {
		return false, fmt.Errorf("TXT record verification failed")
	}

	dmarcVerified, err := s.verifyDMARCRecord(domain.Domain, domain.DMARCRecord)
	if err != nil {
		return false, err
	}

	if !dmarcVerified {
		return false, fmt.Errorf("DMARC record verification failed")
	}

	dkimVerified, err := s.verifyDKIMRecord(domain.Domain, domain.DKIMSelector, domain.DKIMPublicKey)
	if err != nil {
		return false, err
	}

	if !dkimVerified {
		return false, fmt.Errorf("DKIM record verification failed")
	}

	return true, nil
}

func (s *DomainService) verifyTXTRecord(domain, expectedRecord string) (bool, error) {
	records, err := net.LookupTXT(domain)
	if err != nil {
		return false, err
	}

	for _, record := range records {
		if record == expectedRecord {
			return true, nil
		}
	}

	return false, nil
}

func (s *DomainService) verifyDMARCRecord(domain, expectedRecord string) (bool, error) {
	records, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		return false, err
	}

	for _, record := range records {
		if record == expectedRecord {
			return true, nil
		}
	}

	return false, nil
}

func (s *DomainService) verifyDKIMRecord(domain, selector, publicKey string) (bool, error) {
	dkimDomain := fmt.Sprintf("%s._domainkey.%s", selector, domain)
	fmt.Printf("Looking up DKIM record for: %s\n", dkimDomain)

	records, err := net.LookupTXT(dkimDomain)
	if err != nil {
		fmt.Printf("Error looking up DKIM record: %v\n", err)
		return false, err
	}

	fmt.Printf("Found %d DKIM records\n", len(records))

	expectedRecord := fmt.Sprintf("v=DKIM1; k=rsa; p=%s", publicKey)
	fmt.Printf("Expected DKIM record: %s\n", expectedRecord)

	for i, record := range records {
		fmt.Printf("Record %d: %s\n", i+1, record)
		if strings.TrimSpace(record) == expectedRecord {
			fmt.Println("DKIM record match found!")
			return true, nil
		}
	}

	fmt.Println("No matching DKIM record found")
	return false, nil
}

func (s *DomainService) DeleteDomain(id string) error {
	if err := s.DomainRepo.DeleteDomain(id); err != nil {
		return err
	}
	return nil
}

func (s *DomainService) GetDomain(uuid string) (*model.DomainsResponse, error) {
	getDomain, err := s.DomainRepo.GetDomain(uuid)
	if err != nil {
		return nil, err
	}
	return getDomain, nil
}

func (s *DomainService) GetAllDomains(userId string, page int, pageSize int, searchQuery string) (repository.PaginatedResult, error) {
	paginationParams := repository.PaginationParams{Page: page, PageSize: pageSize}
	getAllDoamins, err := s.DomainRepo.GetAllDomains(userId, searchQuery, paginationParams)

	if err != nil {
		return repository.PaginatedResult{}, nil
	}

	return getAllDoamins, nil
}
