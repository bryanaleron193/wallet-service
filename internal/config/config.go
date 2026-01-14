package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DB  DatabaseConfig
	App AppConfig
	JWT JWTConfig
}

type DatabaseConfig struct {
	URL          string        `mapstructure:"DATABASE_URL"`
	MaxOpenConns int           `mapstructure:"DB_MAX_OPEN_CONNS"`
	MaxIdleTime  time.Duration `mapstructure:"DB_MAX_IDLE_TIME"`
	MaxLifetime  time.Duration `mapstructure:"DB_MAX_LIFETIME"`
}

type AppConfig struct {
	Port string
	Env  string
}

type JWTConfig struct {
	Secret string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		App: AppConfig{
			Port: getEnvAsString("APP_PORT", "8081"),
			Env:  getEnvAsString("APP_ENV", "local"),
		},
		DB: DatabaseConfig{
			URL:          getEnvAsString("DATABASE_URL", "postgres://admin:digitalwallet@localhost:5432/wallet?sslmode=disable"),
			MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleTime:  getEnvAsDuration("DB_MAX_IDLE_TIME", 5*time.Minute),
			MaxLifetime:  getEnvAsDuration("DB_MAX_LIFETIME", 1*time.Hour),
		},
		JWT: JWTConfig{
			Secret: getEnvAsString("JWT_SECRET", "digital_wallet"),
		},
	}
}

func getEnvAsString(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}

	return defaultVal
}
