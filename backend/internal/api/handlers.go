package api

import (
	"encoding/json"
	"net/http"

	"certificate-status-page/internal/certificate"
	"certificate-status-page/internal/config"
)

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
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info)
	}
}

func AllEndpointsHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		results := make(map[string]*certificate.SSLDetails)
		client := &http.Client{}
		for _, endpoint := range cfg.Endpoints {
			info, err := certificate.GetCertificateInfo(client, endpoint)
			if err != nil {
				results[endpoint] = nil
				continue
			}
			results[endpoint] = info
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}
