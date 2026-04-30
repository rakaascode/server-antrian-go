package config

import "os"

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPass      string
	DBName      string
	JWTSecret   string
	ServerPort  string
	FonnteToken string
	RedisAddr   string
}

func Load() *Config {
	return &Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPass:      getEnv("DB_PASS", ""),
		DBName:      getEnv("DB_NAME", "antrian_db"),
		JWTSecret:   getEnv("JWT_SECRET", "rahasia-jwt"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		FonnteToken: getEnv("FONNTE_TOKEN", ""),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
