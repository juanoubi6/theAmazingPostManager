package config

import (
	"bufio"
	"os"
	"strings"
)

type Config struct {
	ENV        string
	PORT       string
	JWT_SECRET string
	CORS       string

	REDIS_ADDR              string
	LAST_COMMENTS_LIST_NAME string
	LAST_POSTS_LIST_NAME    string
	LAST_COMMENTS_LENGTH    string
	LAST_POSTS_LENGTH       string

	RABBITMQ_USER     string
	RABBITMQ_PASSWORD string
	RABBITMQ_HOST     string
	RABBITMQ_PORT     string
	RABBIT_NOTIFICATION_EXCHANGE string

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
		err := readEnv()
		if err != nil {
			panic(err)
		}
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

		REDIS_ADDR: GetEnv("REDIS_ADDR", ":6379"),

		RABBITMQ_HOST:     GetEnv("RABBITMQ_HOST", "localhost"),
		RABBITMQ_PORT:     GetEnv("RABBITMQ_PORT", "5672"),
		RABBITMQ_USER:     GetEnv("RABBITMQ_USER", "guest"),
		RABBITMQ_PASSWORD: GetEnv("RABBITMQ_PASSWORD", "guest"),
		RABBIT_NOTIFICATION_EXCHANGE: GetEnv("RABBIT_NOTIFICATION_EXCHANGE", "notifications_queue"),

		LAST_COMMENTS_LIST_NAME: GetEnv("LAST_COMMENTS_LIST_NAME", "lastComments"),
		LAST_POSTS_LIST_NAME:    GetEnv("LAST_POSTS_LIST_NAME", "lastPosts"),
		LAST_COMMENTS_LENGTH:    GetEnv("LAST_COMMENTS_LENGTH", "10"),
		LAST_POSTS_LENGTH:       GetEnv("LAST_POSTS_LENGTH", "10"),

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

func readEnv() error {
	file, err := os.Open(".env")
	if err != nil {
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		values := strings.Split(scanner.Text(), "=")
		if len(values) == 2 {
			err = os.Setenv(values[0], values[1])
			if err != nil {
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
