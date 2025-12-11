package config

import (
	"fmt"
	"os"
)

// Config holds application configuration loaded from environment.
type Config struct {
	MasterDSN   string
	Replica1DSN string
	Replica2DSN string
	Replica3DSN string
	Replica4DSN string
	Replica5DSN string
	APIPort            string
}

// Load reads environment variables and returns Config.
func Load() (Config, error) {
	cfg := Config{
		MasterDSN:   os.Getenv("MASTER_DSN"),
		Replica1DSN: os.Getenv("REPLICA1_DSN"),
		Replica2DSN: os.Getenv("REPLICA2_DSN"),
		Replica3DSN: os.Getenv("REPLICA3_DSN"),
		Replica4DSN: os.Getenv("REPLICA4_DSN"),
		Replica5DSN: os.Getenv("REPLICA5_DSN"),
		APIPort:            getenvDefault("API_PORT", "8080"),
	}

	if cfg.MasterDSN == "" {
		return cfg, fmt.Errorf("MASTER_DSN is required")
	}
	if cfg.Replica1DSN == "" {
		return cfg, fmt.Errorf("REPLICA1_DSN is required")
	}
	if cfg.Replica2DSN == "" {
		return cfg, fmt.Errorf("REPLICA2_DSN is required")
	}
	if cfg.Replica3DSN == "" {
		return cfg, fmt.Errorf("REPLICA3_DSN is required")
	}
	if cfg.Replica4DSN == "" {
		return cfg, fmt.Errorf("REPLICA4_DSN is required")
	}
	if cfg.Replica5DSN == "" {
		return cfg, fmt.Errorf("REPLICA5_DSN is required")
	}

	return cfg, nil
}

// ReplicaDSNs returns the ordered list of replica DSNs.
func (c Config) ReplicaDSNs() []string {
	return []string{
		c.Replica1DSN,
		c.Replica2DSN,
		c.Replica3DSN,
		c.Replica4DSN,
		c.Replica5DSN,
	}
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}



