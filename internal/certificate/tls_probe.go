package certificate

import (
	"crypto/tls"
	"errors"
	"net"
	"net/url"
	"sync"
	"time"
)

type TLSProtocolSupport struct {
	Protocol  string
	Supported bool
}

var probedVersions = []struct {
	v    uint16
	name string
}{
	{tls.VersionTLS10, "TLS 1.0"},
	{tls.VersionTLS11, "TLS 1.1"},
	{tls.VersionTLS12, "TLS 1.2"},
	{tls.VersionTLS13, "TLS 1.3"},
}

// probeTLSVersions runs one handshake per TLS version pinned to that exact version
func probeTLSVersions(hostPort string, timeout time.Duration) []TLSProtocolSupport {
	results := make([]TLSProtocolSupport, len(probedVersions))
	var wg sync.WaitGroup
	dialer := &net.Dialer{Timeout: timeout}
	for i, pv := range probedVersions {
		wg.Add(1)
		go func(i int, version uint16, name string) {
			defer wg.Done()
			results[i].Protocol = name
			conn, err := tls.DialWithDialer(dialer, "tcp", hostPort, &tls.Config{
				InsecureSkipVerify: true,
				MinVersion:         version,
				MaxVersion:         version,
			})
			if err != nil {
				return
			}
			conn.Close()
			results[i].Supported = true
		}(i, pv.v, pv.name)
	}
	wg.Wait()
	return results
}

// addressToHostPort parses a URL and returns host:port, defaulting to :443.
func addressToHostPort(address string) (string, error) {
	u, err := url.Parse(address)
	if err != nil {
		return "", err
	}
	if u.Host == "" {
		return "", errors.New("no host in URL")
	}
	host := u.Host
	if _, _, err := net.SplitHostPort(host); err != nil {
		host = net.JoinHostPort(host, "443")
	}
	return host, nil
}
