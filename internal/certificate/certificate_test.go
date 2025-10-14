package certificate

import (
	"errors"
	"net/http"
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

	if info.Version == 0 {
		t.Error("Expected non-zero TLS version")
	}

	if info.CipherSuite == 0 {
		t.Error("Expected non-zero CipherSuite")
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
