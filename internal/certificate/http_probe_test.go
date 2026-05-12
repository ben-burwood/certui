package certificate

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func startALPNServer(t *testing.T, protos []string) (string, func()) {
	t.Helper()
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.TLS = &tls.Config{NextProtos: protos}
	srv.StartTLS()
	hostPort, err := addressToHostPort(srv.URL)
	if err != nil {
		srv.Close()
		t.Fatalf("addressToHostPort(%q): %v", srv.URL, err)
	}
	return hostPort, srv.Close
}

func TestProbeHTTPVersions_PinnedServer(t *testing.T) {
	hostPort, cleanup := startALPNServer(t, []string{"h2", "http/1.1"})
	defer cleanup()

	results := probeHTTPVersions(hostPort, 3*time.Second)

	if len(results) != len(probedHTTPVersions) {
		t.Fatalf("expected %d results, got %d", len(probedHTTPVersions), len(results))
	}

	for _, r := range results {
		switch r.Version {
		case "HTTP/1.1", "HTTP/2":
			if !r.Supported {
				t.Errorf("expected %s to be supported", r.Version)
			}
		case "HTTP/3":
			if r.Supported {
				t.Errorf("expected HTTP/3 NOT supported against a TLS-only test server")
			}
		default:
			t.Errorf("unexpected version %q", r.Version)
		}
	}
}

func TestProbeHTTPVersions_HTTP11Only(t *testing.T) {
	hostPort, cleanup := startALPNServer(t, []string{"http/1.1"})
	defer cleanup()

	results := probeHTTPVersions(hostPort, 3*time.Second)

	for _, r := range results {
		switch r.Version {
		case "HTTP/1.1":
			if !r.Supported {
				t.Errorf("expected HTTP/1.1 supported")
			}
		default:
			if r.Supported {
				t.Errorf("expected %s NOT supported by HTTP/1.1-only server", r.Version)
			}
		}
	}
}

func TestProbeHTTPVersions_UnreachableHost(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	addr := lis.Addr().String()
	lis.Close()

	results := probeHTTPVersions(addr, 2*time.Second)

	for _, r := range results {
		if r.Supported {
			t.Errorf("%s should not be supported against closed port", r.Version)
		}
	}
}
