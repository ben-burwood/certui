package domain

import (
	"testing"
	"time"
)

func TestWhoisForDomain_Real(t *testing.T) {
	domain := Domain("google.com")
	result, err := WhoisForDomain(domain)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Registrar == "" {
		t.Error("expected Registrar to be non-empty")
	}

	if len(result.NameServers) == 0 {
		t.Error("expected at least one name server")
	}

	// ExpirationDate may change, but should be in the future
	if result.ExpirationDate.Before(time.Now()) {
		t.Errorf("expected expiration date in the future, got %v", result.ExpirationDate)
	}
}

func TestWhoisForDomain_InvalidDomain(t *testing.T) {
	domain := Domain("notarealdomain.tld")
	_, err := WhoisForDomain(domain)
	if err == nil {
		t.Error("expected error for invalid domain, got nil")
	}
}
