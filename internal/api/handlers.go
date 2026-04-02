package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"certui/internal/config"
	"certui/internal/domain"
)

// AllEndpointsHandler handles requests for all Endpoints Details
func AllEndpointsHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		results := make(map[domain.Domain]*EndpointDetails)
		var mu sync.Mutex
		var wg sync.WaitGroup

		client := &http.Client{}
		for _, endpoint := range cfg.Endpoints {
			wg.Add(1)
			go func(ep domain.Domain) {
				defer wg.Done()
				details := fetchEndpointDetails(client, ep)
				mu.Lock()
				results[ep] = details
				mu.Unlock()
			}(endpoint)
		}
		wg.Wait()

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

		var mu sync.Mutex
		var wg sync.WaitGroup

		client := &http.Client{Timeout: 10 * time.Second}
		for _, endpoint := range cfg.Endpoints {
			wg.Add(1)
			go func(ep domain.Domain) {
				defer wg.Done()
				details := fetchEndpointDetails(client, ep)
				wrapped := struct {
					Endpoint domain.Domain   `json:"endpoint"`
					Details  EndpointDetails `json:"details"`
				}{Endpoint: ep, Details: *details}

				b, _ := json.Marshal(wrapped)
				mu.Lock()
				fmt.Fprintf(w, "data: %s\n\n", b) // Send Single Endpoint Data
				w.(http.Flusher).Flush()
				mu.Unlock()
			}(endpoint)
		}
		wg.Wait()

		fmt.Fprintf(w, "event: done\ndata: {}\n\n") // Send Done Event
		w.(http.Flusher).Flush()
	}
}
