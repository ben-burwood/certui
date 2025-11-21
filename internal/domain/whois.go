package domain

import (
	"time"

	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

type WhoisDetails struct {
	Registrar      string
	NameServers    []string
	ExpirationDate time.Time
}

func NewWhoisDetails(whoisInfo *whoisparser.WhoisInfo) *WhoisDetails {
	return &WhoisDetails{
		Registrar:      whoisInfo.Registrar.Name,
		NameServers:    whoisInfo.Domain.NameServers,
		ExpirationDate: *whoisInfo.Domain.ExpirationDateInTime,
	}
}

// WhoisForDomain retrieves the whois information for a domain.
func WhoisForDomain(domain Domain) (*WhoisDetails, error) {
	whoisRaw, err := whois.Whois(string(domain.stripHttpPrefix()))
	if err != nil {
		return nil, err
	}

	whois, err := whoisparser.Parse(whoisRaw)
	if err != nil {
		return nil, err
	}

	return NewWhoisDetails(&whois), nil
}
