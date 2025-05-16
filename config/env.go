package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost string
	Port       string

	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	JWTExpInSec int64
	JWTSecret   string
}

var Envs = initConfig()

func initConfig() Config {
	if os.Getenv("ENVIRONMENT") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	return Config{
		PublicHost:  getEnv("PUBLIC_HOST", "http://localhost"),
		Port:        getEnv("PORT", "8080"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "passwordNotFound"),
		DBName:      getEnv("DB_NAME", "nameNotFound"),
		JWTExpInSec: getEnvAsInt("JWT_EXP_IN_SEC", 3600*24*7),
		JWTSecret:   getEnv("JWT_SECRET", "secretNotFoundAndNotSecretAnymore"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}
