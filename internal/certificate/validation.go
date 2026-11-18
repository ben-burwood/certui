package certificate

import (
	"certui/internal/domain"
	"context"
	"crypto/x509"
	"log"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"
)

// ValidationMethod is the inferred ACME challenge type used to issue a certificate.
// Certificates do not record how they were validated, so this is a best-effort inference.
type ValidationMethod struct {
	Method    string // "DNS-01" or "HTTP-01"
	Confirmed bool   // true only when the method is certain (wildcard certificate or CAA restriction)
	Reason    string
}

const (
	MethodDNS  = "DNS-01"
	MethodHTTP = "HTTP-01"
)

// validationRecords holds the DNS records used to infer the validation method
type validationRecords struct {
	caa       *domain.CAAPolicy
	caaErr    error
	challenge domain.ACMEChallenge
}

// lookupValidationRecords concurrently fetches the CAA policy and _acme-challenge record for a host
func lookupValidationRecords(ctx context.Context, host string, timeout time.Duration) validationRecords {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var records validationRecords
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if records.caa, records.caaErr = domain.LookupCAA(ctx, host); records.caaErr != nil {
			log.Printf("Error looking up CAA for %s: %v", host, records.caaErr)
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		if records.challenge, err = domain.LookupACMEChallenge(ctx, host); err != nil {
			log.Printf("Error looking up _acme-challenge for %s: %v", host, err)
		}
	}()
	wg.Wait()
	return records
}

// inferValidationMethod infers the issuance challenge type from the leaf certificate and DNS records.
//   - Wildcard SANs can only be issued via DNS-01.
//   - A CAA record restricting validationmethods (RFC 8657) to a single method determines it.
//   - An _acme-challenge CNAME (delegation) or TXT record indicates DNS-01.
//   - Otherwise HTTP-01 is assumed.
func inferValidationMethod(leaf *x509.Certificate, records validationRecords) ValidationMethod {
	if name := wildcardName(leaf); name != "" {
		return ValidationMethod{Method: MethodDNS, Confirmed: true, Reason: "Wildcard certificate (" + name + ") requires DNS-01"}
	}

	if method := caaValidationMethod(records.caa); method != "" {
		return ValidationMethod{Method: method, Confirmed: true, Reason: "CAA on " + records.caa.Domain + " only permits " + method}
	}

	if records.challenge.CNAME != "" {
		return ValidationMethod{Method: MethodDNS, Reason: "_acme-challenge is delegated via CNAME to " + records.challenge.CNAME}
	}
	if records.challenge.HasTXT {
		return ValidationMethod{Method: MethodDNS, Reason: "_acme-challenge TXT record found"}
	}

	if records.caaErr != nil {
		return ValidationMethod{Method: MethodHTTP, Reason: "No wildcard or _acme-challenge record found (CAA lookup failed)"}
	}
	return ValidationMethod{Method: MethodHTTP, Reason: "No wildcard, CAA restriction or _acme-challenge record found"}
}

// caaValidationMethod returns the single ACME challenge type permitted by the CAA issue records, or "" if
// any issuer is unrestricted or more than one challenge type is permitted.
func caaValidationMethod(caa *domain.CAAPolicy) string {
	if caa == nil {
		return ""
	}
	permitted := map[string]bool{}
	for _, issuer := range caa.Issuers {
		if issuer.CA == "" {
			continue // Forbids issuance - doesn't widen the permitted methods
		}
		if len(issuer.ValidationMethods) == 0 {
			return ""
		}
		for _, method := range issuer.ValidationMethods {
			permitted[method] = true
		}
	}
	if len(permitted) != 1 {
		return ""
	}
	for method := range permitted {
		switch method {
		case "dns-01":
			return MethodDNS
		case "http-01":
			return MethodHTTP
		}
	}
	return ""
}

// addressDomainName returns the address's hostname if it is a multi-label domain name, or "" for
// IP literals and single-label hosts, which have no CAA or _acme-challenge records
func addressDomainName(address string) string {
	u, err := url.Parse(address)
	if err != nil {
		return ""
	}
	host := strings.TrimSuffix(u.Hostname(), ".")
	if net.ParseIP(host) != nil || !strings.Contains(host, ".") {
		return ""
	}
	return host
}
