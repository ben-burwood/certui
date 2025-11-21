package main

import (
	"certui/internal/api"
	"certui/internal/config"
	"net/http"
	"os"
)

const (
	CertuiConfigPathEnvVar = "CERTUI_CONFIG_PATH"
)

func main() {
	cfg, err := loadConfiguration()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/endpoint", api.EndpointHandler(cfg))
	mux.HandleFunc("GET /api/endpoints", api.AllEndpointsHandler(cfg))
	mux.HandleFunc("GET /api/endpoints-sse", api.EndpointHandlerSSE(cfg))

	// Serve Static Frontend
	mux.Handle("/", http.FileServer(http.Dir("./frontend/dist")))

	http.ListenAndServe("[::]:8080", api.CORSMiddleware(mux))
}

// loadConfiguration loads the configuration from the path specified in the CERTUI_CONFIG_PATH environment variable
func loadConfiguration() (*config.Config, error) {
	configPath := os.Getenv(CertuiConfigPathEnvVar)
	return config.LoadConfig(configPath)
}
