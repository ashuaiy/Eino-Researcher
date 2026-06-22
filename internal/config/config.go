package config

import (
	"os"
	"strconv"
)

type Config struct {
	App       AppConfig
	Postgres  PostgresConfig
	Redis     RedisConfig
	LLM       LLMConfig
	Embedding EmbeddingConfig
}

type AppConfig struct {
	Port string
	Env  string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

type RedisConfig struct {
	Addr string
}

type LLMConfig struct {
	BaseURL string
	APIKey  string
	Model   string
}

type EmbeddingConfig struct {
	BaseURL string
	APIKey  string
	Model   string
	Dim     int
}

func Load() Config {
	return Config{
		App: AppConfig{
			Port: getEnv("APP_PORT", "8080"),
			Env:  getEnv("APP_ENV", "development"),
		},
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
			Database: getEnv("POSTGRES_DB", "eino_researcher"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Addr: getEnv("REDIS_ADDR", "localhost:6379"),
		},
		LLM: LLMConfig{
			BaseURL: getEnv("LLM_BASE_URL", "https://api.openai.com/v1"),
			APIKey:  getEnv("LLM_API_KEY", ""),
			Model:   getEnv("LLM_MODEL", "gpt-4o-mini"),
		},
		Embedding: EmbeddingConfig{
			BaseURL: getEnv("EMBEDDING_BASE_URL", "https://api.openai.com/v1"),
			APIKey:  getEnv("EMBEDDING_API_KEY", ""),
			Model:   getEnv("EMBEDDING_MODEL", "text-embedding-3-small"),
			Dim:     getEnvInt("EMBEDDING_DIM", 1536),
		},
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
