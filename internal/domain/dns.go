package domain

import (
	"context"
	"net"
	"os"
	"strings"
	"sync"

	"codeberg.org/miekg/dns"
	"codeberg.org/miekg/dns/dnsconf"
	"codeberg.org/miekg/dns/dnsutil"
)

const (
	// DNSServerEnvVar overrides the DNS server used for CAA and _acme-challenge lookups (host or host:port)
	DNSServerEnvVar   = "CERTUI_DNS_SERVER"
	fallbackDNSServer = "1.1.1.1:53"
	ednsUDPSize       = 1232
)

var dnsServer = sync.OnceValue(func() string {
	if server := os.Getenv(DNSServerEnvVar); server != "" {
		if _, _, err := net.SplitHostPort(server); err != nil {
			server = net.JoinHostPort(strings.Trim(server, "[]"), "53")
		}
		return server
	}
	if cfg, err := dnsconf.FromFile("/etc/resolv.conf"); err == nil && len(cfg.Servers) > 0 {
		return net.JoinHostPort(cfg.Servers[0], cfg.Port)
	}
	return fallbackDNSServer
})

// query sends a recursive query to the configured DNS server, retrying over TCP if the UDP response is truncated
func query(ctx context.Context, name string, qtype uint16) (*dns.Msg, error) {
	client := dns.NewClient()
	m := dns.NewMsg(name, qtype)
	m.UDPSize = ednsUDPSize

	r, _, err := client.Exchange(ctx, m, "udp", dnsServer())
	if err == nil && r.Truncated {
		m = dns.NewMsg(name, qtype)
		r, _, err = client.Exchange(ctx, m, "tcp", dnsServer())
	}
	if err != nil {
		return nil, err
	}
	if r.Rcode != dns.RcodeSuccess && r.Rcode != dns.RcodeNameError {
		return nil, &net.DNSError{Err: "query failed: " + dnsutil.RcodeToString(r.Rcode), Name: name, Server: dnsServer()}
	}
	return r, nil
}

// ACMEChallenge describes the _acme-challenge record for a host
type ACMEChallenge struct {
	CNAME  string // Target of a CNAME delegation (e.g. to acme-dns), "" if not delegated
	HasTXT bool
}

// LookupACMEChallenge queries the TXT record of _acme-challenge.<host>; any CNAME chain is returned in the same answer
func LookupACMEChallenge(ctx context.Context, host string) (ACMEChallenge, error) {
	r, err := query(ctx, "_acme-challenge."+host, dns.TypeTXT)
	if err != nil {
		return ACMEChallenge{}, err
	}
	var challenge ACMEChallenge
	for _, rr := range r.Answer {
		switch rr := rr.(type) {
		case *dns.CNAME:
			if challenge.CNAME == "" {
				challenge.CNAME = strings.TrimSuffix(rr.Target, ".")
			}
		case *dns.TXT:
			challenge.HasTXT = true
		}
	}
	return challenge, nil
}

// CAAIssuer is a single CAA issue/issuewild property (RFC 8659), with RFC 8657 parameters
type CAAIssuer struct {
	CA                string   // Issuer domain (e.g. letsencrypt.org), "" means no CA may issue
	ValidationMethods []string // Permitted ACME challenge types, empty means unrestricted
	AccountURI        string
}

// CAAPolicy is the relevant CAA record set for a host
type CAAPolicy struct {
	Domain          string // Domain the records were found on - the host itself or a parent
	Issuers         []CAAIssuer
	WildcardIssuers []CAAIssuer // Only set when issuewild records are present
}

// LookupCAA finds the relevant CAA record set by climbing from the host towards the root (RFC 8659 Section 3).
// Returns nil if no CAA records exist, meaning any CA may issue.
func LookupCAA(ctx context.Context, host string) (*CAAPolicy, error) {
	labels := strings.Split(strings.TrimSuffix(host, "."), ".")
	for i := range labels {
		name := strings.Join(labels[i:], ".")
		r, err := query(ctx, name, dns.TypeCAA)
		if err != nil {
			return nil, err
		}

		policy := CAAPolicy{Domain: name}
		found := false
		for _, rr := range r.Answer {
			caa, ok := rr.(*dns.CAA)
			if !ok {
				continue
			}
			found = true
			switch strings.ToLower(caa.Tag) {
			case "issue":
				policy.Issuers = append(policy.Issuers, parseCAAIssuer(caa.Value))
			case "issuewild":
				policy.WildcardIssuers = append(policy.WildcardIssuers, parseCAAIssuer(caa.Value))
			}
		}
		if found {
			return &policy, nil
		}
	}
	return nil, nil
}

// parseCAAIssuer parses an issue/issuewild value, e.g. "letsencrypt.org; validationmethods=dns-01; accounturi=..."
func parseCAAIssuer(value string) CAAIssuer {
	parts := strings.Split(value, ";")
	issuer := CAAIssuer{CA: strings.TrimSpace(parts[0])}
	for _, param := range parts[1:] {
		key, val, ok := strings.Cut(strings.TrimSpace(param), "=")
		if !ok {
			continue
		}
		switch strings.ToLower(strings.TrimSpace(key)) {
		case "validationmethods":
			for _, method := range strings.Split(val, ",") {
				if method = strings.ToLower(strings.TrimSpace(method)); method != "" {
					issuer.ValidationMethods = append(issuer.ValidationMethods, method)
				}
			}
		case "accounturi":
			issuer.AccountURI = strings.TrimSpace(val)
		}
	}
	return issuer
}
