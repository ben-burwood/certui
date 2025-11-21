package domain

import (
	"net"
)

type Domain string

type HostAddress string

// getDomainHostAddress returns the Resolved IP Address of a Domain
func getDomainHostAddress(domain Domain) (HostAddress, error) {
	ips, err := net.LookupHost(string(domain))
	if err != nil {
		return "", err
	}
	return HostAddress(ips[0]), nil
}

type DomainDetails struct {
	Domain      Domain
	HostAddress HostAddress
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
