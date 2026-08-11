package main

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	RiotAPIKey string
	Region     string
	GameName   string
	Tagline    string
}

func LoadConfig() (*Config, error) {
	loadDotEnv(".env")

	cfg := &Config{
		RiotAPIKey: os.Getenv("RIOT_API_KEY"),
		Region:     getEnvOrDefault("RIOT_REGION", "europe"),
		GameName:   os.Getenv("GAME_NAME"),
		Tagline:    os.Getenv("TAGLINE"),
	}

	if cfg.RiotAPIKey == "" {
		return nil, fmt.Errorf("missing required environment variable: RIOT_API_KEY")
	}
	if cfg.GameName == "" || cfg.Tagline == "" {
		return nil, fmt.Errorf("missing required player target variables: GAME_NAME or TAGLINE")
	}

	return cfg, nil
}

// Lightweight .env file reader using native standard library
func loadDotEnv(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return // File optional; will fall back to standard OS environment variables
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, val, found := strings.Cut(line, "=")
		if found {
			cleanVal := strings.Trim(strings.TrimSpace(val), `"'`)
			os.Setenv(strings.TrimSpace(key), cleanVal)
		}
	}
}

func getEnvOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
