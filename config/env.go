package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// ENVIRONMENT must be set explicitly to one of these values - there is no
// default. Only EnvDevelopment loads .env from disk via godotenv; EnvTest and
// EnvProduction assume the process environment is already fully populated
// (by the Makefile's shell-sourcing for test, or by ECS for production).
const (
	EnvDevelopment = "development"
	EnvTest        = "test"
	EnvProduction  = "production"
)

func isValidEnvironment(env string) bool {
	switch env {
	case EnvDevelopment, EnvTest, EnvProduction:
		return true
	default:
		return false
	}
}

type Config struct {
	Environment string
	PublicHost  string
	Port        string

	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	JWTExpInSec int64
	JWTSecret   string

	VAPIDPublicKey  string
	VAPIDPrivateKey string
	VAPIDSubject    string

	OpenAIAPIKey   string
	S3AssetsBucket string
	AWSRegion      string
}

var Envs = initConfig()

func initConfig() Config {
	environment := os.Getenv("ENVIRONMENT")
	if !isValidEnvironment(environment) {
		log.Fatalf("ENVIRONMENT must be set explicitly to %q, %q, or %q (got %q)",
			EnvDevelopment, EnvTest, EnvProduction, environment)
	}

	if environment == EnvDevelopment {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	return Config{
		Environment: environment,
		PublicHost:  getEnv("PUBLIC_HOST", "http://localhost"),
		Port:        getEnv("PORT", "8080"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "passwordNotFound"),
		DBName:      getEnv("DB_NAME", "nameNotFound"),
		DBSSLMode:   getEnv("DB_SSL_MODE", defaultDBSSLMode(environment)),
		JWTExpInSec: getEnvAsInt("JWT_EXP_IN_SEC", 3600*24*7),
		JWTSecret:   getEnv("JWT_SECRET", "secretNotFoundAndNotSecretAnymore"),

		VAPIDPublicKey:  getEnv("VAPID_PUBLIC_KEY", ""),
		VAPIDPrivateKey: getEnv("VAPID_PRIVATE_KEY", ""),
		VAPIDSubject:    getEnv("VAPID_SUBJECT", "mailto:noreply@animalfamily.app"),

		OpenAIAPIKey:   getEnv("OPENAI_API_KEY", ""),
		S3AssetsBucket: getEnv("S3_ASSETS_BUCKET", "brindl-assets"),
		AWSRegion:      getEnv("AWS_REGION", "us-east-1"),
	}
}

// defaultDBSSLMode is the DB_SSL_MODE fallback when it isn't set explicitly:
// RDS requires SSL, so production defaults to requiring it; a local Postgres
// typically isn't configured for SSL at all.
func defaultDBSSLMode(environment string) string {
	if environment == EnvProduction {
		return "require"
	}

	return "disable"
}

func (c Config) IsProduction() bool {
	return c.Environment == EnvProduction
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
