package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	UserSvcURL string
}

func Load() *Config {

	_ = godotenv.Load()

	config := &Config{
		Port:       getEnv("PORT", "8080"),
		UserSvcURL: getEnv("USER_SVC_URL", ""),
	}

	if config.UserSvcURL == "" {
		log.Fatal("user service url not set")
	}

	return config
}

func getEnv(key, defVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defVal
}
