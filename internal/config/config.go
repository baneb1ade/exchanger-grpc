package config

import (
	"github.com/joho/godotenv"
	"os"
	"time"
)

type Config struct {
	GRPCConfig
	StorageConfig
}

type GRPCConfig struct {
	Port    string
	Timeout time.Duration
}

type StorageConfig struct {
	StorageUser     string
	StoragePass     string
	StorageHost     string
	StoragePort     string
	StorageDatabase string
}

func getEnvWithDefault(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func MustLoad(cfgPath string) *Config {
	if cfgPath != "" {
		err := godotenv.Load(cfgPath)
		if err != nil {
			panic(err)
		}
	}

	cfg := Config{
		GRPCConfig: GRPCConfig{
			Port:    getEnvWithDefault("GRPC_PORT", "44044"),
			Timeout: 1 * time.Hour,
		},
		StorageConfig: StorageConfig{
			StorageUser:     getEnvWithDefault("STORAGE_USER", "postgres"),
			StoragePass:     getEnvWithDefault("STORAGE_PASSWORD", "postgres"),
			StorageHost:     getEnvWithDefault("STORAGE_HOST", "localhost"),
			StoragePort:     getEnvWithDefault("STORAGE_PORT", "5432"),
			StorageDatabase: getEnvWithDefault("STORAGE_DB", "postgres"),
		},
	}
	return &cfg
}
