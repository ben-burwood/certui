package api

import (
	"log"
	"net/http"
	"sync"
	"time"

	"certui/internal/certificate"
	"certui/internal/domain"
)

// EndpointDetails contains the certificate, domain, and WHOIS details for an endpoint.
type EndpointDetails struct {
	Domain domain.DomainDetails
	Whois  *domain.WhoisDetails
	SSL    *certificate.SSLDetails
}

const endpointCacheTTL = time.Hour

var endpointCache sync.Map // map[domain.Domain]endpointCacheEntry

type endpointCacheEntry struct {
	details   *EndpointDetails
	expiresAt time.Time
}

// fetchEndpointDetails fetches certificate, domain, and WHOIS details concurrently for a single endpoint.
// force - cache-buster - bypasses the cache read, writes fresh cache entries
func fetchEndpointDetails(client *http.Client, endpoint domain.Domain, force bool) *EndpointDetails {
	if !force {
		if v, ok := endpointCache.Load(endpoint); ok {
			entry := v.(endpointCacheEntry)
			if time.Now().Before(entry.expiresAt) {
				return entry.details
			}
		}
	}

	var ssl *certificate.SSLDetails
	var domainDetails domain.DomainDetails
	var whoisDetails *domain.WhoisDetails
	var wg sync.WaitGroup

	wg.Add(3)
	go func() {
		defer wg.Done()
		var err error
		ssl, err = certificate.GetCertificateInfo(client, endpoint)
		if err != nil {
			log.Printf("Error getting certificate info for %s: %v", endpoint, err)
		}
	}()
	go func() {
		defer wg.Done()
		domainDetails = domain.GetDomainDetails(endpoint)
	}()
	go func() {
		defer wg.Done()
		var err error
		whoisDetails, err = domain.WhoisForDomain(endpoint)
		if err != nil {
			log.Printf("Error getting WHOIS info for %s: %v", endpoint, err)
		}
	}()
	wg.Wait()

	details := &EndpointDetails{Domain: domainDetails, Whois: whoisDetails, SSL: ssl}
	endpointCache.Store(endpoint, endpointCacheEntry{
		details:   details,
		expiresAt: time.Now().Add(endpointCacheTTL),
	})
	return details
}
