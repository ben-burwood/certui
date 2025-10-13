package main

import (
	"certificate-status-page/internal/api"
	"certificate-status-page/internal/config"
	"fmt"
	"net/http"
	"os"
)

const (
	CertificateStatusPageConfigPathEnvVar = "CERTIFICATE_STATUS_PAGE_CONFIG_PATH"
)

func main() {
	cfg, err := loadConfiguration()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /endpoint", api.EndpointHandler(cfg))
	mux.HandleFunc("GET /endpoints", api.AllEndpointsHandler(cfg))

	listenAddr := fmt.Sprintf("[::]:%d", cfg.Web.Port)
	http.ListenAndServe(listenAddr, mux)
}

// loadConfiguration loads the configuration from the path specified in the CERTIFICATE_STATUS_PAGE_CONFIG_PATH environment variable
func loadConfiguration() (*config.Config, error) {
	configPath := os.Getenv(CertificateStatusPageConfigPathEnvVar)
	return config.LoadConfig(configPath)
}
