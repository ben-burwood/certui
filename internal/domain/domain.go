package domain

import (
	"net"
)

type Domain string

type HostAddress string

// GetDomainAddress returns the Resolved IP Address of a Domain
func GetDomainAddress(domain Domain) (HostAddress, error) {
	ips, err := net.LookupHost(string(domain))
	if err != nil {
		return "", err
	}
	return HostAddress(ips[0]), nil
}
