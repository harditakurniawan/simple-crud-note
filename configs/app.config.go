package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv(envFile string) {
	err := godotenv.Load(envFile)

	if err != nil {
		log.Panicln("Error loading .env file")
	}
}

func getEnv(key string) string {
	return os.Getenv(key)
}

type AppConfig struct {
	APP_NAME             string
	APP_PREFIX           string
	APP_PORT             string
	DATABASE_URL         string
	REDIS_URL            string
	JWT_PRIVATE_KEY_PATH string
	JWT_PUBLIC_KEY_PATH  string
}

func LoadAppConfig() AppConfig {
	return AppConfig{
		APP_NAME:             getEnv("APP_NAME"),
		APP_PREFIX:           getEnv("APP_PREFIX"),
		APP_PORT:             getEnv("APP_PORT"),
		DATABASE_URL:         getEnv("DATABASE_URL"),
		REDIS_URL:            getEnv("REDIS_URL"),
		JWT_PRIVATE_KEY_PATH: getEnv("JWT_PRIVATE_KEY_PATH"),
		JWT_PUBLIC_KEY_PATH:  getEnv("JWT_PUBLIC_KEY_PATH"),
	}
}
