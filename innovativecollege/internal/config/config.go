package config

import (
	"os"
)

type Config struct {
	MongoURI     string
	DatabaseName string
	Port         string
}

func Load() *Config {
	return &Config{
		MongoURI:     getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		DatabaseName: getEnv("DATABASE_NAME", "innovativecollege"),
		Port:         getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

