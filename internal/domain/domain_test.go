package domain

import (
	"testing"
)

func TestGetDomainAddress(t *testing.T) {
	domain := Domain("example.com")
	details := GetDomainDetails(domain)
	if details.Domain != domain {
		t.Fatalf("Expected domain %s, got %s", domain, details.Domain)
	}
	if details.HostAddress == "" {
		t.Fatalf("Expected a non-empty address for domain %s", domain)
	}
	if details.Resolves == false {
		t.Fatalf("Expected domain %s to resolve", domain)
	}
}

func TestGetDomainAddressInvalid(t *testing.T) {
	domain := Domain("invalid")
	details := GetDomainDetails(domain)
	if details.Domain != domain {
		t.Fatalf("Expected domain %s, got %s", domain, details.Domain)
	}
	if details.HostAddress != "" {
		t.Fatalf("Expected an empty address for domain %s", domain)
	}
	if details.Resolves == true {
		t.Fatalf("Expected domain %s to not resolve", domain)
	}
}
