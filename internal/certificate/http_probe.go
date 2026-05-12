package certificate

import (
	"context"
	"crypto/tls"
	"net"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
)

type HTTPVersionSupport struct {
	Version   string
	Supported bool
}

var probedHTTPVersions = []struct {
	alpn string
	name string
	quic bool
}{
	{"http/1.1", "HTTP/1.1", false},
	{"h2", "HTTP/2", false},
	{"h3", "HTTP/3", true},
}

// Shared so quic-go's UDP buffer-size check logs at most once at process startup
// (per quic-go/wiki/UDP-Buffer-Sizes), and so probes reuse a single UDP socket.
var (
	h3TransportOnce sync.Once
	h3Transport     *quic.Transport
)

func sharedH3Transport() *quic.Transport {
	h3TransportOnce.Do(func() {
		udp, err := net.ListenUDP("udp", &net.UDPAddr{})
		if err != nil {
			return
		}
		h3Transport = &quic.Transport{Conn: udp}
	})
	return h3Transport
}

func probeHTTPVersions(hostPort string, timeout time.Duration) []HTTPVersionSupport {
	results := make([]HTTPVersionSupport, len(probedHTTPVersions))
	var wg sync.WaitGroup
	dialer := &net.Dialer{Timeout: timeout}
	for i, pv := range probedHTTPVersions {
		wg.Add(1)
		go func(i int, alpn, name string, isQUIC bool) {
			defer wg.Done()
			results[i].Version = name
			if isQUIC {
				results[i].Supported = probeH3(hostPort, alpn, timeout)
				return
			}
			results[i].Supported = probeALPN(dialer, hostPort, alpn)
		}(i, pv.alpn, pv.name, pv.quic)
	}
	wg.Wait()
	return results
}

func probeALPN(dialer *net.Dialer, hostPort, alpn string) bool {
	conn, err := tls.DialWithDialer(dialer, "tcp", hostPort, &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{alpn},
	})
	if err != nil {
		return false
	}
	defer conn.Close()
	return conn.ConnectionState().NegotiatedProtocol == alpn
}

func probeH3(hostPort, alpn string, timeout time.Duration) bool {
	tr := sharedH3Transport()
	if tr == nil {
		return false
	}
	addr, err := net.ResolveUDPAddr("udp", hostPort)
	if err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	conn, err := tr.Dial(ctx, addr, &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{alpn},
	}, nil)
	if err != nil {
		return false
	}
	_ = conn.CloseWithError(0, "")
	return true
}
