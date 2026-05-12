package certificate

import (
	"certui/internal/domain"
	"crypto/tls"
	"net/http"
	"sync"
	"time"
)

type SSLDetails struct {
	HandshakeComplete bool
	DidResume         bool
	CipherSuite       uint16
	PeerCertificates  []CertificateDetails
	TLSProtocols      []TLSProtocolSupport
	HTTPVersions      []HTTPVersionSupport
}

type CertificateDetails struct {
	Subject            string
	Issuer             string
	NotBefore          time.Time
	NotAfter           time.Time
	SignatureAlgorithm string
	PublicKeyAlgorithm string
}

// IsExpired checks if the first certificate in the SSLDetails is expired
func (s *SSLDetails) IsExpired() bool {
	if len(s.PeerCertificates) == 0 {
		return false
	}
	return time.Now().After(s.PeerCertificates[0].NotAfter)
}

const defaultProbeTimeout = 5 * time.Second

// GetCertificateInfo retrieves SSL/TLS certificate information from the specified Address
// using the provided HTTP Client - Certificate Verification is Skipped for the request.
func GetCertificateInfo(client *http.Client, address domain.Domain) (*SSLDetails, error) {
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	var probeWG sync.WaitGroup
	var probes []TLSProtocolSupport
	var httpProbes []HTTPVersionSupport
	if hostPort, err := addressToHostPort(string(address)); err == nil {
		timeout := probeTimeout(client)
		probeWG.Add(2)
		go func() {
			defer probeWG.Done()
			probes = probeTLSVersions(hostPort, timeout)
		}()
		go func() {
			defer probeWG.Done()
			httpProbes = probeHTTPVersions(hostPort, timeout)
		}()
	}

	resp, err := client.Get(string(address))
	if err != nil {
		probeWG.Wait()
		return nil, err
	}
	defer resp.Body.Close()

	sslInfo := SSLDetails{
		HandshakeComplete: resp.TLS.HandshakeComplete,
		DidResume:         resp.TLS.DidResume,
		CipherSuite:       resp.TLS.CipherSuite,
	}
	for _, cert := range resp.TLS.PeerCertificates {
		sslInfo.PeerCertificates = append(sslInfo.PeerCertificates, CertificateDetails{
			Subject:            cert.Subject.String(),
			Issuer:             cert.Issuer.String(),
			NotBefore:          cert.NotBefore,
			NotAfter:           cert.NotAfter,
			SignatureAlgorithm: cert.SignatureAlgorithm.String(),
			PublicKeyAlgorithm: cert.PublicKeyAlgorithm.String(),
		})
	}

	probeWG.Wait()
	sslInfo.TLSProtocols = probes
	sslInfo.HTTPVersions = httpProbes
	return &sslInfo, nil
}

func probeTimeout(client *http.Client) time.Duration {
	if client.Timeout > 0 {
		return client.Timeout
	}
	return defaultProbeTimeout
}
