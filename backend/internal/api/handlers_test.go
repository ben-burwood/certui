package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"certificate-status-page/internal/certificate"
	"certificate-status-page/internal/config"
)

func TestEndpointHandler_Success(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	req := httptest.NewRequest("GET", "/ssl?endpoint="+ts.URL, nil)
	rw := httptest.NewRecorder()

	dummyCfg := &config.Config{}
	handler := EndpointHandler(dummyCfg)
	handler(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rw.Code)
	}

	var details certificate.SSLDetails
	err := json.NewDecoder(rw.Body).Decode(&details)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if details.Version == 0 {
		t.Error("expected non-zero Version in SSLDetails")
	}
}

func TestEndpointHandler_MissingParam(t *testing.T) {
	req := httptest.NewRequest("GET", "/ssl", nil)
	rw := httptest.NewRecorder()

	dummyCfg := &config.Config{}
	handler := EndpointHandler(dummyCfg)
	handler(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rw.Code)
	}
}
