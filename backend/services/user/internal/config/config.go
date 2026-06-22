package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GrpcPort    string
	DatabaseURL string
}

func Load() *Config {

	_ = godotenv.Load()

	config := &Config{
		GrpcPort:    getEnv("GRPC_PORT", ":50051"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}

	if config.DatabaseURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	return config
}

func getEnv(key, defVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defVal
}
