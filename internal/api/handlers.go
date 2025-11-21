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
		info, err := certificate.GetCertificateInfo(client, endpointDomain)
		if err != nil {
			info = nil
		}

		domainDetails := domain.GetDomainDetails(endpointDomain)

		response := EndpointDetails{
			Domain: domainDetails,
			SSL:    info,
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

			info, err := certificate.GetCertificateInfo(client, endpointDomain)
			if err != nil {
				info = nil
			}

			domainDetails := domain.GetDomainDetails(endpointDomain)

			response := EndpointDetails{
				Domain: domainDetails,
				SSL:    info,
			}

			results[endpoint] = &response
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}
