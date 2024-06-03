package config

import "os"

type DBConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewDBConfig() DBConfig {
	return DBConfig{
		Host:     GetEnv("DB_HOST", "localhost"),
		User:     GetEnv("DB_USER", "postgres"),
		Password: GetEnv("DB_PASSWORD", ""),
		DBName:   GetEnv("DB_NAME", "library"),
		SSLMode:  GetEnv("DB_SSLMODE", "disable"),
	}
}

func GetEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
