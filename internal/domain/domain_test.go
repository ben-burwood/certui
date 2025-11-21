package domain

import (
	"testing"
)

func TestGetDomainAddress(t *testing.T) {
	domain := Domain("example.com")
	addr, err := GetDomainAddress(domain)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if addr == "" {
		t.Fatalf("Expected a non-empty address for domain %s", domain)
	}
	t.Logf("Address for %s: %s", domain, addr)
}

func TestGetDomainAddressInvalid(t *testing.T) {
	domain := Domain("invalid")
	addr, err := GetDomainAddress(domain)
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}
	if addr != "" {
		t.Fatalf("Expected an empty address for domain %s", domain)
	}
	t.Logf("Address for %s: %s", domain, addr)
}
