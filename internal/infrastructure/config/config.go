package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	XMLRiver XMLRiverConfig
	XMLStock XMLStockConfig
	Kafka    KafkaConfig
	Async    AsyncConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ServerConfig struct {
	Port           string
	TrustedProxies []string
}

type XMLRiverConfig struct {
	UserID  string
	APIKey  string
	BaseURL string
}

type XMLStockConfig struct {
	UserID  string
	APIKey  string
	BaseURL string
}

type KafkaConfig struct {
	Brokers []string
}

type AsyncConfig struct {
	WorkerCount int
	BatchSize   int
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "go_seo"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Port:           getEnv("SERVER_PORT", "8080"),
			TrustedProxies: getEnvAsStringSlice("SERVER_TRUSTED_PROXIES", []string{"127.0.0.1", "::1"}),
		},
		XMLRiver: XMLRiverConfig{
			UserID:  getEnv("XMLRIVER_USER_ID", ""),
			APIKey:  getEnv("XMLRIVER_API_KEY", ""),
			BaseURL: getEnv("XMLRIVER_BASE_URL", "https://xmlriver.com"),
		},
		XMLStock: XMLStockConfig{
			UserID:  getEnv("XMLSTOCK_USER_ID", ""),
			APIKey:  getEnv("XMLSTOCK_API_KEY", ""),
			BaseURL: getEnv("XMLSTOCK_BASE_URL", "https://xmlstock.com"),
		},
		Kafka: KafkaConfig{
			Brokers: getEnvAsStringSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
		},
		Async: AsyncConfig{
			WorkerCount: getEnvAsInt("WORKER_COUNT", 20),
			BatchSize:   getEnvAsInt("BATCH_SIZE", 100),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsStringSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}
