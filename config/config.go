package config

import (
	"fmt"
	"os"
)

type Config struct {
	GoogleMapsKey     string
	ProjectID         string
	CredetialFilePath string
}

func load() (Config, error) {
	cfg := Config{
		GoogleMapsKey:     os.Getenv("GOOGLE_MAP_KEY"),
		ProjectID:         os.Getenv("FIREBASE_PROJECT_ID"),
		CredetialFilePath: os.Getenv("GOOGLE_APPLICATION_CREDENTIAL"),
	}
	if cfg.GoogleMapsKey == "" {
		return cfg, fmt.Errorf("GOOGLE_MAP_KEY empty")
	}
	return cfg, nil
}
