package config

import (
	"log"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v2"
)

type Config struct {
	DatabaseURL      string `yaml:"database_url"`
	DatabaseName     string `yaml:"database_name"`
	StreamCollection string `yaml:"stream_collection"`
	BatchCollection  string `yaml:"batch_collection"`
}

func LoadConfig(path string) *Config {
	cleanPath := filepath.Clean(path)

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	if envDBURL := os.Getenv("DB_URL"); envDBURL != "" {
		cfg.DatabaseURL = envDBURL
	}
	if envDBName := os.Getenv("DB_NAME"); envDBName != "" {
		cfg.DatabaseName = envDBName
	}
	if envStreamColl := os.Getenv("STREAM_COLLECTION"); envStreamColl != "" {
		cfg.StreamCollection = envStreamColl
	}
	if envBatchColl := os.Getenv("BATCH_COLLECTION"); envBatchColl != "" {
		cfg.BatchCollection = envBatchColl
	}

	return &cfg
}
