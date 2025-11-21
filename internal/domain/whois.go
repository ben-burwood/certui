package domain

import (
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

func WhoisForDomain(domain Domain) (*whoisparser.WhoisInfo, error) {
	whoisRaw, err := whois.Whois(string(domain.stripHttpPrefix()))
	if err != nil {
		return nil, err
	}

	whois, err := whoisparser.Parse(whoisRaw)
	if err != nil {
		return nil, err
	}

	return &whois, nil
}
