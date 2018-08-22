package config

import "os"

type Config struct {
	ENV        string
	PORT       string
	JWT_SECRET string
	CORS       string

	DB_TYPE     string
	DB_USERNAME string
	DB_PASSWORD string
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string

}

var instance *Config

func GetConfig() *Config {
	if instance == nil {
		config := newConfig()
		instance = &config
	}
	return instance
}

func newConfig() Config {
	return Config{
		ENV:        GetEnv("ENV", "develop"),
		PORT:       GetEnv("PORT", "5003"),
		JWT_SECRET: GetEnv("JWT_SECRET", "j8Ah4kO3"),
		CORS:       GetEnv("CORS", ""),

		DB_TYPE:     GetEnv("DB_TYPE", "mysql"),
		DB_USERNAME: GetEnv("DB_USERNAME", "root"),
		DB_PASSWORD: GetEnv("DB_PASSWORD", "root"),
		DB_HOST:     GetEnv("DB_HOST", "127.0.0.1"),
		DB_PORT:     GetEnv("DB_PORT", "3306"),
		DB_NAME:     GetEnv("DB_NAME", "amazing-code-database"),
	}
}

func GetEnv(key, fallback string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return fallback
}
