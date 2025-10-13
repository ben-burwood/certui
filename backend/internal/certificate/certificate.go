package certificate

import (
	"crypto/tls"
	"net/http"
	"time"
)

type SSLDetails struct {
	Version           uint16
	HandshakeComplete bool
	DidResume         bool
	CipherSuite       uint16
	PeerCertificates  []CertificateDetails
}

type CertificateDetails struct {
	Subject            string
	Issuer             string
	NotBefore          time.Time
	NotAfter           time.Time
	SignatureAlgorithm string
	PublicKeyAlgorithm string
}

// GetCertificateInfo retrieves SSL/TLS certificate information from the specified Address
// using the provided HTTP Client - Certificate Verification is Skipped for the request.
func GetCertificateInfo(client *http.Client, address string) (*SSLDetails, error) {
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	resp, err := client.Get(address)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	sslInfo := SSLDetails{
		Version:           resp.TLS.Version,
		HandshakeComplete: resp.TLS.HandshakeComplete,
		DidResume:         resp.TLS.DidResume,
		CipherSuite:       resp.TLS.CipherSuite,
	}
	// Retrieve information about the peer certificates
	for _, cert := range resp.TLS.PeerCertificates {
		peerCertificate := CertificateDetails{
			Subject:            cert.Subject.String(),
			Issuer:             cert.Issuer.String(),
			NotBefore:          cert.NotBefore,
			NotAfter:           cert.NotAfter,
			SignatureAlgorithm: cert.SignatureAlgorithm.String(),
			PublicKeyAlgorithm: cert.PublicKeyAlgorithm.String(),
		}
		sslInfo.PeerCertificates = append(sslInfo.PeerCertificates, peerCertificate)
	}

	return &sslInfo, nil
}
