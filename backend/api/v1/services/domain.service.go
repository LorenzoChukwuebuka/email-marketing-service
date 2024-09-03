package services

import (
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	"fmt"
	"net"
)

type DomainService struct {
	DomainRepo *repository.DomainRepository
}

func NewDomainService(domainRepo *repository.DomainRepository) *DomainService {
	return &DomainService{
		DomainRepo: domainRepo,
	}
}

func (s *DomainService) CreateDomain(d *model.Domains) (map[string]interface{}, error) {

	isDomainExists, err := s.DomainRepo.CheckIfDomainExists(d)

	if err != nil {
		return nil, err
	}

	if isDomainExists {
		return nil, fmt.Errorf("domain already exists")
	}

	//verify domain to make sure it is existing before saving

	verifyDomain := s.verifyDomain(d.Domain)

	if !verifyDomain {
		return nil, fmt.Errorf("your domain is not verified")
	}

	if err := s.DomainRepo.CreateDomain(d); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"domain": d,
	}, nil
}

func (s *DomainService) verifyDomain(domain string) bool {
	// Check if the domain resolves to an IP address
	ips, err := net.LookupIP(domain)
	if err != nil {
		fmt.Printf("Error looking up IP for %s: %v\n", domain, err)
		return false
	}
	fmt.Printf("IP addresses for %s: %v\n", domain, ips)

	// Check for MX records
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		fmt.Printf("Error looking up MX records for %s: %v\n", domain, err)
	} else {
		fmt.Printf("MX records for %s:\n", domain)
		for _, mx := range mxRecords {
			fmt.Printf("  %v (priority: %v)\n", mx.Host, mx.Pref)
		}
	}

	// Check for A records
	aRecords, err := net.LookupHost(domain)
	if err != nil {
		fmt.Printf("Error looking up A records for %s: %v\n", domain, err)
	} else {
		fmt.Printf("A records for %s: %v\n", domain, aRecords)
	}

	// Check for NS records
	nsRecords, err := net.LookupNS(domain)
	if err != nil {
		fmt.Printf("Error looking up NS records for %s: %v\n", domain, err)
	} else {
		fmt.Printf("NS records for %s:\n", domain)
		for _, ns := range nsRecords {
			fmt.Printf("  %v\n", ns.Host)
		}
	}

	// If we've made it this far without returning false, the domain likely exists
	fmt.Printf("The domain %s exists and is properly configured.\n", domain)
	return true
}
