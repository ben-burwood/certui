package api

import (
	"encoding/json"
	"net/http"

	"certui/internal/certificate"
	"certui/internal/config"
)

type SSLDetailsWithExpired struct {
	certificate.SSLDetails
	IsExpired bool
}

func EndpointHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpoint := r.URL.Query().Get("endpoint")
		if endpoint == "" {
			http.Error(w, "Missing endpoint parameter", http.StatusBadRequest)
			return
		}

		client := &http.Client{}
		info, err := certificate.GetCertificateInfo(client, endpoint)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SSLDetailsWithExpired{
			SSLDetails: *info,
			IsExpired:  info.IsExpired(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func AllEndpointsHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		results := make(map[string]*SSLDetailsWithExpired)
		client := &http.Client{}
		for _, endpoint := range cfg.Endpoints {
			info, err := certificate.GetCertificateInfo(client, endpoint)
			if err != nil {
				results[endpoint] = nil
				continue
			}

			response := SSLDetailsWithExpired{
				SSLDetails: *info,
				IsExpired:  info.IsExpired(),
			}

			results[endpoint] = &response
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}
