package certificate

import (
	"certui/internal/domain"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestGetCertificateInfo(t *testing.T) {
	client := &http.Client{}

	info, err := GetCertificateInfo(client, "https://www.google.com")
	if err != nil {
		t.Fatalf("GetCertificateInfo returned error: %v", err)
	}

	if info == nil {
		t.Fatal("Expected non-nil SSLDetails")
	}

	if !info.HandshakeComplete {
		t.Error("Expected HandshakeComplete to be true")
	}

	if info.CipherSuite == 0 {
		t.Error("Expected non-zero CipherSuite")
	}

	if len(info.TLSProtocols) == 0 {
		t.Error("Expected TLSProtocols to be populated")
	}

	if len(info.PeerCertificates) == 0 {
		t.Error("Expected at least one peer certificate")
	}

	cert := info.PeerCertificates[0]
	if cert.Subject == "" {
		t.Error("Expected non-empty Subject")
	}
	if cert.Issuer == "" {
		t.Error("Expected non-empty Issuer")
	}
	if cert.NotAfter.Before(cert.NotBefore) {
		t.Error("NotAfter should be after NotBefore")
	}
	if cert.SignatureAlgorithm == "" {
		t.Error("Expected non-empty SignatureAlgorithm")
	}
	if cert.PublicKeyAlgorithm == "" {
		t.Error("Expected non-empty PublicKeyAlgorithm")
	}
	if len(cert.DNSNames) == 0 {
		t.Error("Expected at least one DNSName")
	}
	if len(cert.Fingerprint) != 64 {
		t.Errorf("Expected 64-character SHA-256 Fingerprint, got %q", cert.Fingerprint)
	}
}

func TestNewCertificateDetails_Wildcard(t *testing.T) {
	raw := []byte("raw")
	cert := newCertificateDetails(&x509.Certificate{Raw: raw, DNSNames: []string{"example.com", "*.example.com"}})
	if !cert.IsWildcard {
		t.Error("Expected IsWildcard to be true")
	}
	sum := sha256.Sum256(raw)
	if want := hex.EncodeToString(sum[:]); cert.Fingerprint != want {
		t.Errorf("Expected Fingerprint %q, got %q", want, cert.Fingerprint)
	}
}

func TestNewCertificateDetails_NotWildcard(t *testing.T) {
	cert := newCertificateDetails(&x509.Certificate{DNSNames: []string{"example.com", "www.example.com"}})
	if cert.IsWildcard {
		t.Error("Expected IsWildcard to be false")
	}
}

func TestGetCertificateInfo_Error(t *testing.T) {
	client := &http.Client{}
	_, err := GetCertificateInfo(client, "https://invalid.invalid")
	if err == nil {
		t.Fatal("Expected error for invalid address")
	}
	if !errors.Is(err, err) {
		t.Error("Expected a network error")
	}
}

func TestInferValidationMethod_Wildcard(t *testing.T) {
	leaf := &x509.Certificate{DNSNames: []string{"example.com", "*.example.com"}}
	v := inferValidationMethod(leaf, validationRecords{})
	if v.Method != MethodDNS || !v.Confirmed {
		t.Errorf("Expected confirmed DNS-01 for wildcard, got %+v", v)
	}
}

func TestInferValidationMethod_NoEvidence(t *testing.T) {
	leaf := &x509.Certificate{DNSNames: []string{"example.com"}}
	v := inferValidationMethod(leaf, validationRecords{})
	if v.Method != MethodHTTP || v.Confirmed {
		t.Errorf("Expected unconfirmed HTTP-01, got %+v", v)
	}
}

func TestInferValidationMethod_ACMEChallenge(t *testing.T) {
	leaf := &x509.Certificate{DNSNames: []string{"example.com"}}
	for _, challenge := range []domain.ACMEChallenge{{CNAME: "abc.auth.acme-dns.io"}, {HasTXT: true}} {
		v := inferValidationMethod(leaf, validationRecords{challenge: challenge})
		if v.Method != MethodDNS || v.Confirmed {
			t.Errorf("Expected unconfirmed DNS-01 for %+v, got %+v", challenge, v)
		}
	}
}

func TestInferValidationMethod_CAA(t *testing.T) {
	leaf := &x509.Certificate{DNSNames: []string{"example.com"}}
	tests := []struct {
		name      string
		issuers   []domain.CAAIssuer
		method    string
		confirmed bool
	}{
		{"dns only", []domain.CAAIssuer{{CA: "letsencrypt.org", ValidationMethods: []string{"dns-01"}}}, MethodDNS, true},
		{"http only", []domain.CAAIssuer{{CA: "letsencrypt.org", ValidationMethods: []string{"http-01"}}}, MethodHTTP, true},
		{"multiple methods", []domain.CAAIssuer{{CA: "letsencrypt.org", ValidationMethods: []string{"dns-01", "http-01"}}}, MethodHTTP, false},
		{"unrestricted issuer", []domain.CAAIssuer{
			{CA: "letsencrypt.org", ValidationMethods: []string{"dns-01"}},
			{CA: "sectigo.com"},
		}, MethodHTTP, false},
		{"deny issuer ignored", []domain.CAAIssuer{
			{CA: "letsencrypt.org", ValidationMethods: []string{"dns-01"}},
			{CA: ""},
		}, MethodDNS, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			caa := &domain.CAAPolicy{Domain: "example.com", Issuers: tt.issuers}
			v := inferValidationMethod(leaf, validationRecords{caa: caa})
			if v.Method != tt.method || v.Confirmed != tt.confirmed {
				t.Errorf("Expected %s (confirmed=%v), got %+v", tt.method, tt.confirmed, v)
			}
		})
	}
}

func TestInferValidationMethod_CAAError(t *testing.T) {
	leaf := &x509.Certificate{DNSNames: []string{"example.com"}}
	v := inferValidationMethod(leaf, validationRecords{caaErr: errors.New("timeout")})
	if v.Method != MethodHTTP || v.Confirmed || !strings.Contains(v.Reason, "CAA lookup failed") {
		t.Errorf("Expected unconfirmed HTTP-01 noting the CAA failure, got %+v", v)
	}
}

func TestAddressDomainName(t *testing.T) {
	tests := map[string]string{
		"https://www.example.com":     "www.example.com",
		"https://example.com.:8443/x": "example.com",
		"https://192.168.1.10":        "",
		"https://[2001:db8::1]:443":   "",
		"https://localhost":           "",
	}
	for address, want := range tests {
		if got := addressDomainName(address); got != want {
			t.Errorf("addressDomainName(%q) = %q, want %q", address, got, want)
		}
	}
}
