package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APPEnv    string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	JWTSecret string
	REDISHost string
	REDISPort string
	TESTMode  string
}

func LoadConfig() Config {
	_ = godotenv.Load()

	config := Config{
		APPEnv:    getEnv("APP_ENV", "development"),
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnv("DB_PORT", "3306"),
		DBUser:    getEnv("DB_USER", "root"),
		DBPass:    getEnv("DB_PASS", ""),
		DBName:    getEnv("DB_NAME", "testdb"),
		JWTSecret: getEnv("JWT_SECRET", "defaultsecret"),
		REDISHost: getEnv("REDIS_HOST", "redis"),
		REDISPort: getEnv("REDIS_PORT", "6379"),
		TESTMode:  getEnv("TEST_MODE", "true"),
	}
	return config
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
