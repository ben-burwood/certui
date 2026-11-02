package config

import (
	"certui/internal/domain"
	"errors"
	"reflect"
	"testing"
)

func TestLoadConfigFromEndpointsEnv(t *testing.T) {
	t.Setenv(EndpointsEnvVar, "https://a.com, b.com,,https://c.com ")

	cfg, err := LoadConfig("does-not-exist.yml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []domain.Domain{"https://a.com", "b.com", "https://c.com"}
	if !reflect.DeepEqual(cfg.Endpoints, want) {
		t.Errorf("got %v, want %v", cfg.Endpoints, want)
	}
}

func TestLoadConfigFromEmptyEndpointsEnv(t *testing.T) {
	t.Setenv(EndpointsEnvVar, " , ")

	if _, err := LoadConfig(""); !errors.Is(err, ErrNoEndpointsInConfig) {
		t.Errorf("got %v, want %v", err, ErrNoEndpointsInConfig)
	}
}
