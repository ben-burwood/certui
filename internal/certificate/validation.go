package certificate

import (
	"context"
	"crypto/x509"
	"net"
	"net/url"
	"strings"
	"time"
)

// ValidationMethod is the inferred ACME challenge type used to issue a certificate.
// Certificates do not record how they were validated, so this is a best-effort inference.
type ValidationMethod struct {
	Method    string // "DNS-01" or "HTTP-01"
	Confirmed bool   // true only when the method is certain (e.g. wildcard certificates)
	Reason    string
}

const (
	MethodDNS  = "DNS-01"
	MethodHTTP = "HTTP-01"
)

// inferValidationMethod infers the issuance challenge type from the leaf certificate and DNS.
//   - Wildcard SANs can only be issued via DNS-01.
//   - An _acme-challenge CNAME (delegation) or TXT record indicates DNS-01.
//   - Otherwise HTTP-01 is assumed.
func inferValidationMethod(leaf *x509.Certificate, address string, timeout time.Duration) ValidationMethod {
	if name := wildcardName(leaf); name != "" {
		return ValidationMethod{Method: MethodDNS, Confirmed: true, Reason: "Wildcard certificate (" + name + ") requires DNS-01"}
	}

	if host := addressHostname(address); host != "" {
		if reason, ok := hasACMEChallengeRecord(host, timeout); ok {
			return ValidationMethod{Method: MethodDNS, Reason: reason}
		}
	}

	return ValidationMethod{Method: MethodHTTP, Reason: "No wildcard or _acme-challenge record found"}
}

// hasACMEChallengeRecord checks for a delegated (CNAME) or present (TXT) _acme-challenge record.
func hasACMEChallengeRecord(host string, timeout time.Duration) (string, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	name := "_acme-challenge." + host
	if cname, err := net.DefaultResolver.LookupCNAME(ctx, name); err == nil && !strings.EqualFold(strings.TrimSuffix(cname, "."), name) {
		return name + " is delegated via CNAME to " + strings.TrimSuffix(cname, "."), true
	}
	if txts, err := net.DefaultResolver.LookupTXT(ctx, name); err == nil && len(txts) > 0 {
		return name + " TXT record found", true
	}
	return "", false
}

func addressHostname(address string) string {
	u, err := url.Parse(address)
	if err != nil {
		return ""
	}
	return u.Hostname()
}
