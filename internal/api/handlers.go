package api

import (
	"encoding/json"
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

// EndpointHandler handles requests for a single Endpoint Details
func EndpointHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpoint := r.URL.Query().Get("endpoint")
		if endpoint == "" {
			http.Error(w, "Missing endpoint parameter", http.StatusBadRequest)
			return
		}
		endpointDomain := domain.Domain(endpoint)

		client := &http.Client{}
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

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
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
