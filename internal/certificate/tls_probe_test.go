package certificate

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAddressToHostPort(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    string
		wantErr bool
	}{
		{"https no port defaults to 443", "https://example.com", "example.com:443", false},
		{"https with explicit port", "https://example.com:8443", "example.com:8443", false},
		{"ipv4 with port", "https://127.0.0.1:54321", "127.0.0.1:54321", false},
		{"empty string errors", "", "", true},
		{"no host errors", "https://", "", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := addressToHostPort(tc.in)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func startPinnedTLSServer(t *testing.T, version uint16) (string, func()) {
	t.Helper()
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.TLS = &tls.Config{MinVersion: version, MaxVersion: version}
	srv.StartTLS()
	hostPort, err := addressToHostPort(srv.URL)
	if err != nil {
		srv.Close()
		t.Fatalf("addressToHostPort(%q): %v", srv.URL, err)
	}
	return hostPort, srv.Close
}

func TestProbeTLSVersions_PinnedServer(t *testing.T) {
	hostPort, cleanup := startPinnedTLSServer(t, tls.VersionTLS12)
	defer cleanup()

	results := probeTLSVersions(hostPort, 3*time.Second)

	if len(results) != len(probedVersions) {
		t.Fatalf("expected %d results, got %d", len(probedVersions), len(results))
	}

	for _, r := range results {
		switch r.Protocol {
		case "TLS 1.2":
			if !r.Supported {
				t.Errorf("expected TLS 1.2 to be supported")
			}
		default:
			if r.Supported {
				t.Errorf("expected %s NOT supported by TLS 1.2-only server", r.Protocol)
			}
		}
	}
}

func TestProbeTLSVersions_UnreachableHost(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	addr := lis.Addr().String()
	lis.Close()

	results := probeTLSVersions(addr, 2*time.Second)

	for _, r := range results {
		if r.Supported {
			t.Errorf("%s should not be supported against closed port", r.Protocol)
		}
	}
}
