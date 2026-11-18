package domain

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func TestParseCAAIssuer(t *testing.T) {
	tests := []struct {
		value string
		want  CAAIssuer
	}{
		{"letsencrypt.org", CAAIssuer{CA: "letsencrypt.org"}},
		{";", CAAIssuer{CA: ""}},
		{"letsencrypt.org; validationmethods=dns-01", CAAIssuer{CA: "letsencrypt.org", ValidationMethods: []string{"dns-01"}}},
		{
			"letsencrypt.org; accounturi=https://acme-v02.api.letsencrypt.org/acme/acct/1; validationmethods=DNS-01,http-01",
			CAAIssuer{CA: "letsencrypt.org", ValidationMethods: []string{"dns-01", "http-01"}, AccountURI: "https://acme-v02.api.letsencrypt.org/acme/acct/1"},
		},
	}
	for _, tt := range tests {
		if got := parseCAAIssuer(tt.value); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("parseCAAIssuer(%q) = %+v, want %+v", tt.value, got, tt.want)
		}
	}
}

func TestLookupCAA(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// google.com publishes "issue pki.goog"; www.google.com has none so the lookup must climb to the parent
	policy, err := LookupCAA(ctx, "www.google.com")
	if err != nil {
		t.Fatalf("LookupCAA returned error: %v", err)
	}
	if policy == nil {
		t.Fatal("Expected a CAA policy for www.google.com")
	}
	if policy.Domain != "google.com" {
		t.Errorf("Expected policy from google.com, got %q", policy.Domain)
	}
	if len(policy.Issuers) == 0 {
		t.Error("Expected at least one issue record")
	}
}

func TestLookupACMEChallenge_None(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	challenge, err := LookupACMEChallenge(ctx, "invalid.invalid")
	if err != nil {
		t.Fatalf("LookupACMEChallenge returned error: %v", err)
	}
	if challenge.CNAME != "" || challenge.HasTXT {
		t.Errorf("Expected no challenge record, got %+v", challenge)
	}
}
