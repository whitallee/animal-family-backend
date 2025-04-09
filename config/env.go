package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost string
	Port       string

	DBUser      string
	DBPassword  string
	DBAddress   string
	DBName      string
	JWTExpInSec int64
	JWTSecret   string
}

var Envs = initConfig()

func initConfig() Config {
	err := godotenv.Load("/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return Config{
		PublicHost:  getEnv("PUBLIC_HOST", "http://localhost"),
		Port:        getEnv("PORT", "8080"),
		DBUser:      getEnv("DB_USER", "root"),
		DBPassword:  getEnv("DB_PASSWORD", "passwordNotFound"),
		DBAddress:   fmt.Sprintf("%s:%s", getEnv("DB_HOST", "hostNotFound"), getEnv("DB_PORT", "portNotFound")),
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
