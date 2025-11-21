package domain

import (
	"net"
	"strings"
)

type Domain string

type Address string

// stripHttpPrefix removes the http/https prefix from a domain string
func (d Domain) stripHttpPrefix() Domain {
	dStr := string(d)
	dStr = strings.TrimPrefix(dStr, "http://")
	dStr = strings.TrimPrefix(dStr, "https://")
	return Domain(dStr)
}

// getDomainHostAddress returns the Resolved IP Address of a Domain
func getDomainHostAddress(domain Domain) (Address, error) {
	domain = domain.stripHttpPrefix()

	ips, err := net.LookupHost(string(domain))
	if err != nil {
		return "", err
	}
	return Address(ips[0]), nil
}

type DomainDetails struct {
	Domain      Domain
	HostAddress Address
	Resolves    bool
}

func GetDomainDetails(domain Domain) DomainDetails {
	hostAddress, err := getDomainHostAddress(domain)
	if err != nil {
		return DomainDetails{
			Domain:      domain,
			HostAddress: "",
			Resolves:    false,
		}
	}
	return DomainDetails{
		Domain:      domain,
		HostAddress: hostAddress,
		Resolves:    true,
	}
}
