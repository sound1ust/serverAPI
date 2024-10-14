package config

import (
	"log"
	"os"

	"github.com/lpernett/godotenv"
)

type Config struct {
	Host       string
	Port       string
	DBName     string
	DBUser     string
	DBPassword string
}

var Env = initConfig()

func initConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("The env. file wasn't found: %v", err)
	}
	return Config{
		Host:       getEnv("HOST", "localhost"),
		Port:       getEnv("PORT", "5432"),
		DBName:     getEnv("POSTGRES_NAME", "postgres"),
		DBUser:     getEnv("POSTGRES_USER", "postgres"),
		DBPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
