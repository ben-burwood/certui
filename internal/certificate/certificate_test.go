package certificate

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"net/http"
	"testing"
	"time"
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
	v := inferValidationMethod(leaf, "https://example.com", time.Second)
	if v.Method != MethodDNS || !v.Confirmed {
		t.Errorf("Expected confirmed DNS-01 for wildcard, got %+v", v)
	}
}

func TestInferValidationMethod_NoEvidence(t *testing.T) {
	leaf := &x509.Certificate{DNSNames: []string{"invalid.invalid"}}
	v := inferValidationMethod(leaf, "https://invalid.invalid", time.Second)
	if v.Method != MethodHTTP || v.Confirmed {
		t.Errorf("Expected unconfirmed HTTP-01, got %+v", v)
	}
}
