package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	DB     DBConfig
	Redis  RedisConfig
	JWT    JWTConfig
	Resend ResendConfig
	Google GoogleConfig
	Email  EmailConfig
}

type ServerConfig struct {
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	AccessSecret     string
	RefreshSecret    string
	AccessExpireMin  int
	RefreshExpireDay int
}

type ResendConfig struct {
	APIKey string
	From   string
}

type GoogleConfig struct {
	ClientID string
}

type EmailConfig struct {
	CodeExpireMin int
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file, using environment variables")
	}
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		DB: DBConfig{
			Host:     mustEnv("DB_HOST"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     mustEnv("DB_USER"),
			Password: mustEnv("DB_PASSWORD"),
			Name:     mustEnv("DB_NAME"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Addr:     mustEnv("REDIS_ADDR"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			AccessSecret:     mustEnv("JWT_ACCESS_SECRET"),
			RefreshSecret:    mustEnv("JWT_REFRESH_SECRET"),
			AccessExpireMin:  getEnvInt("JWT_ACCESS_EXPIRE_MIN", 15),
			RefreshExpireDay: getEnvInt("JWT_REFRESH_EXPIRE_DAY", 7),
		},
		Resend: ResendConfig{
			APIKey: mustEnv("RESEND_API_KEY"),
			From:   mustEnv("RESEND_FROM"),
		},
		Google: GoogleConfig{
			ClientID: mustEnv("GOOGLE_CLIENT_ID"),
		},
		Email: EmailConfig{
			CodeExpireMin: getEnvInt("EMAIL_CODE_EXPIRE_MIN", 5),
		},
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return v
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return i
}
