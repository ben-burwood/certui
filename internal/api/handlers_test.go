package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"certui/internal/certificate"
	"certui/internal/config"
)

func TestEndpointHandler_Success(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	req := httptest.NewRequest("GET", "/ssl?endpoint="+ts.URL, nil)
	rw := httptest.NewRecorder()

	dummyCfg := &config.Config{}
	handler := AllEndpointsHandler(dummyCfg)
	handler(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rw.Code)
	}

	var details certificate.SSLDetails
	err := json.NewDecoder(rw.Body).Decode(&details)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
}
