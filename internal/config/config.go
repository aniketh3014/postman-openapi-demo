package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	DSN      string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("no .env found")
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("could not find DB_PORT %w", err)
	}

	dbConfig := DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     dbPort,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	}

	dbConfig.DSN = fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password,
		dbConfig.DBName, dbConfig.SSLMode,
	)

	config := &Config{
		Server: ServerConfig{
			Port:         os.Getenv("SERVER_PORT"),
			ReadTimeout:  parseDuration(os.Getenv("READ_TIMEOUT")),
			WriteTimeout: parseDuration(os.Getenv("WRITE_TIMEOUT")),
		},
		Database: dbConfig,
	}

	return config, nil
}

func parseDuration(s string) time.Duration {
	duration, err := time.ParseDuration(s)
	if err != nil {
		return 10 * time.Second
	}
	return duration
}
