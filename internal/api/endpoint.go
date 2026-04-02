package api

import (
	"net/http"
	"sync"

	"certui/internal/certificate"
	"certui/internal/domain"
)

// EndpointDetails contains the certificate, domain, and WHOIS details for an endpoint.
type EndpointDetails struct {
	Domain domain.DomainDetails
	Whois  *domain.WhoisDetails
	SSL    *certificate.SSLDetails
}

// fetchEndpointDetails fetches certificate, domain, and WHOIS details concurrently for a single endpoint
func fetchEndpointDetails(client *http.Client, endpoint domain.Domain) *EndpointDetails {
	var ssl *certificate.SSLDetails
	var domainDetails domain.DomainDetails
	var whoisDetails *domain.WhoisDetails
	var wg sync.WaitGroup

	wg.Add(3)
	go func() { defer wg.Done(); ssl, _ = certificate.GetCertificateInfo(client, endpoint) }()
	go func() { defer wg.Done(); domainDetails = domain.GetDomainDetails(endpoint) }()
	go func() { defer wg.Done(); whoisDetails, _ = domain.WhoisForDomain(endpoint) }()
	wg.Wait()

	return &EndpointDetails{Domain: domainDetails, Whois: whoisDetails, SSL: ssl}
}
