package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"certui/internal/certificate"
	"certui/internal/config"
	"certui/internal/domain"
)

type EndpointDetails struct {
	Domain domain.DomainDetails
	Whois  *domain.WhoisDetails
	SSL    *certificate.SSLDetails
}

// AllEndpointsHandler handles requests for all Endpoints Details
func AllEndpointsHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		results := make(map[domain.Domain]*EndpointDetails)
		client := &http.Client{}
		for _, endpoint := range cfg.Endpoints {
			endpointDomain := domain.Domain(endpoint)

			ssl, err := certificate.GetCertificateInfo(client, endpointDomain)
			if err != nil {
				ssl = nil
			}

			domainDetails := domain.GetDomainDetails(endpointDomain)

			whoisDetails, err := domain.WhoisForDomain(endpointDomain)
			if err != nil {
				whoisDetails = nil
			}

			response := EndpointDetails{
				Domain: domainDetails,
				Whois:  whoisDetails,
				SSL:    ssl,
			}

			results[endpoint] = &response
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}

// EndpointHandlerSSE handles Server Sent Events (SSE) Endpoint.
func EndpointHandlerSSE(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		client := &http.Client{}
		for _, endpoint := range cfg.Endpoints {
			endpointDomain := domain.Domain(endpoint)
			ssl, _ := certificate.GetCertificateInfo(client, endpointDomain)
			domainDetails := domain.GetDomainDetails(endpointDomain)
			whoisDetails, _ := domain.WhoisForDomain(endpointDomain)

			response := EndpointDetails{
				Domain: domainDetails,
				Whois:  whoisDetails,
				SSL:    ssl,
			}
			wrapped := struct {
				Endpoint domain.Domain   `json:"endpoint"`
				Details  EndpointDetails `json:"details"`
			}{Endpoint: endpoint, Details: response}

			b, _ := json.Marshal(wrapped)
			fmt.Fprintf(w, "data: %s\n\n", b) // Send Single Endpoint Data
			w.(http.Flusher).Flush()
		}
		fmt.Fprintf(w, "event: done\ndata: {}\n\n") // Send Done Event
		w.(http.Flusher).Flush()
	}
}
